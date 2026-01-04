package dao_test

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"order/internal/dao"
	"order/internal/model"
	"testing"
	"time"
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

func TestOrderDao_Caching(t *testing.T) {
	ctx := context.Background()

	mongoClient, teardownMongo := setupMongo(t)
	defer teardownMongo()

	redisClient, teardownRedis := setupRedis(t)
	defer teardownRedis()

	db := mongoClient.Database("testdb")
	dao := dao.NewOrderDao(db, redisClient)

	_ = db.Collection("orders").Drop(ctx)
	_ = redisClient.FlushDB(ctx)

	order := model.Order{
		UserID:     "user123",
		Status:     "Pending",
		TotalPrice: 20.5,
		ItemIDs:    []string{"item1", "item2"},
		CreatedAt:  time.Now(),
	}
	id, err := dao.Create(ctx, order)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	fetchedOrder, err := dao.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, fetchedOrder.ID)
	assert.Equal(t, "Pending", fetchedOrder.Status)
	cacheKey := "order:id:" + id
	cachedData, err := redisClient.Get(ctx, cacheKey).Result()
	assert.NoError(t, err)
	assert.NotEmpty(t, cachedData)

	var cachedOrder model.Order
	err = json.Unmarshal([]byte(cachedData), &cachedOrder)
	assert.NoError(t, err)
	assert.Equal(t, id, cachedOrder.ID)

	fetchedOrder2, err := dao.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, fetchedOrder2.ID)

	err = dao.UpdateStatus(ctx, id, "Completed")
	assert.NoError(t, err)

	_, err = redisClient.Get(ctx, cacheKey).Result()
	assert.Error(t, err)
	fetchedOrder3, err := dao.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "Completed", fetchedOrder3.Status)
}
