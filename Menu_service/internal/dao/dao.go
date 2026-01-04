package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"foodstore/menu/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MenuRepository interface {
	CreateMenuItem(ctx context.Context, item model.MenuItem) (string, error)
	GetAllMenuItems(ctx context.Context, filter interface{}, limit, skip int64, sortBy string, asc bool) ([]model.MenuItem, error)
	GetMenuItemByID(ctx context.Context, id string) (*model.MenuItem, error)
	Update(ctx context.Context, id string, update bson.M) error
	Delete(ctx context.Context, id string) error
	CountMenuItems(ctx context.Context, filter interface{}) (int64, error)
}

type MongoMenuRepository struct {
	coll  *mongo.Collection
	Cache *redis.Client
}

func NewMenuRepository(db *mongo.Database, cache *redis.Client) MenuRepository {
	return &MongoMenuRepository{
		coll:  db.Collection("menu"),
		Cache: cache,
	}
}

func (r *MongoMenuRepository) CreateMenuItem(ctx context.Context, item model.MenuItem) (string, error) {
	res, err := r.coll.InsertOne(ctx, item)
	if err != nil {
		return "", err
	}
	oid := res.InsertedID.(primitive.ObjectID).Hex()
	return oid, nil
}

func (r *MongoMenuRepository) GetMenuItemByID(ctx context.Context, id string) (*model.MenuItem, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var item model.MenuItem
	if err := r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *MongoMenuRepository) Update(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": update})
	return err
}

func (r *MongoMenuRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (r *MongoMenuRepository) GetAllMenuItems(ctx context.Context, filter interface{}, limit, skip int64, sortBy string, asc bool) ([]model.MenuItem, error) {
	cacheKey := fmt.Sprintf("menu:all:filter=%v:limit=%d:skip=%d:sortBy=%s:asc=%t", filter, limit, skip, sortBy, asc)

	if r.Cache != nil {
		cached, err := r.Cache.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var cachedItems []model.MenuItem
			if err := json.Unmarshal([]byte(cached), &cachedItems); err == nil {
				fmt.Println("Cache hit for GetAllMenuItems")
				return cachedItems, nil
			}
			fmt.Println("Cache unmarshal error:", err)
		} else {
			fmt.Println("Cache miss for GetAllMenuItems")
		}
	}

	opts := options.Find().SetLimit(limit).SetSkip(skip)

	if sortBy != "" {
		order := 1
		if !asc {
			order = -1
		}
		opts.SetSort(bson.D{{Key: sortBy, Value: order}})
	}

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.MenuItem
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	if r.Cache != nil {
		data, err := json.Marshal(items)
		if err == nil {
			err := r.Cache.Set(ctx, cacheKey, data, 10*time.Minute).Err()
			if err != nil {
				fmt.Println("Failed to set cache:", err)
			}
		}
	}

	return items, nil
}

func (r *MongoMenuRepository) CountMenuItems(ctx context.Context, filter interface{}) (int64, error) {
	return r.coll.CountDocuments(ctx, filter)
}
