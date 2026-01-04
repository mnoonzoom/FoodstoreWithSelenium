package dao_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"foodstore/menu/internal/dao"
	"foodstore/menu/internal/model"
	"testing"
)

func setupMongo(t *testing.T) (*mongo.Client, func()) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	assert.NoError(t, err)

	return client, func() {
		_ = client.Disconnect(context.Background())
	}
}

func setupRedis(t *testing.T) (*redis.Client, func()) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	_, err := rdb.Ping(context.Background()).Result()
	assert.NoError(t, err)

	return rdb, func() {
		_ = rdb.FlushDB(context.Background())
		_ = rdb.Close()
	}
}

func TestGetAllMenuItems_WithCache(t *testing.T) {
	t.Log("Integration test for GetAllMenuItems with Redis cache")

	ctx := context.Background()

	mongoClient, mongoTeardown := setupMongo(t)
	defer mongoTeardown()

	redisClient, redisTeardown := setupRedis(t)
	defer redisTeardown()

	db := mongoClient.Database("testdb")
	repo := dao.NewMenuRepository(db, redisClient)

	_ = db.Collection("menu").Drop(ctx)
	_ = redisClient.FlushDB(ctx)
	item := model.MenuItem{
		Name:     "Test Pizza",
		Price:    12.99,
		Category: "Main",
	}
	id, err := repo.CreateMenuItem(ctx, item)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	items, err := repo.GetAllMenuItems(ctx, bson.M{}, 10, 0, "", true)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Test Pizza", items[0].Name)

	cacheKey := fmt.Sprintf("menu:all:filter=%v:limit=%d:skip=%d:sortBy=%s:asc=%t",
		bson.M{}, int64(10), int64(0), "", true)

	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	assert.NoError(t, err)

	var cachedItems []model.MenuItem
	err = json.Unmarshal([]byte(cachedData), &cachedItems)
	assert.NoError(t, err)
	assert.Len(t, cachedItems, 1)
	assert.Equal(t, "Test Pizza", cachedItems[0].Name)
	items2, err := repo.GetAllMenuItems(ctx, bson.M{}, 10, 0, "", true)
	assert.NoError(t, err)
	assert.Len(t, items2, 1)
	assert.Equal(t, "Test Pizza", items2[0].Name)
}
