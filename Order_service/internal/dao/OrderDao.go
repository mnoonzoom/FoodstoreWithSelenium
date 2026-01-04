package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"order/internal/model"
	"time"
)

type OrderDao struct {
	Collection *mongo.Collection
	Cache      *redis.Client
}

func NewOrderDao(db *mongo.Database, cache *redis.Client) *OrderDao {
	return &OrderDao{
		Collection: db.Collection("orders"),
		Cache:      cache,
	}
}

func (r *OrderDao) Create(ctx context.Context, order model.Order) (string, error) {
	res, err := r.Collection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(interface {
		Hex() string
	}).Hex(), nil
}

func (r *OrderDao) GetByID(ctx context.Context, id string) (*model.Order, error) {
	cacheKey := "order:id:" + id

	if r.Cache != nil {
		cached, err := r.Cache.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var cachedOrder model.Order
			if err := json.Unmarshal([]byte(cached), &cachedOrder); err == nil {
				return &cachedOrder, nil
			}
		}
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	var order model.Order
	err = r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		return nil, err
	}

	if r.Cache != nil {
		data, err := json.Marshal(order)
		if err == nil {
			r.Cache.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return &order, nil
}

func (r *OrderDao) UpdateStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status}},
	)

	if err == nil && r.Cache != nil {
		r.Cache.Del(ctx, "order:id:"+id)
		err = r.invalidateUserOrdersCache(ctx)
	}

	return err
}

func (r *OrderDao) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ObjectID: %v", err)
	}

	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": objID})

	if err == nil && r.Cache != nil {
		r.Cache.Del(ctx, "order:id:"+id)
		err = r.invalidateUserOrdersCache(ctx)
	}

	return err
}
func (r *OrderDao) FindOrdersByUserId(ctx context.Context, userId string) ([]model.Order, error) {
	cacheKey := "orders:user:" + userId

	if r.Cache != nil {
		cached, err := r.Cache.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var cachedOrders []model.Order
			if err := json.Unmarshal([]byte(cached), &cachedOrders); err == nil {
				fmt.Println("Cache hit for user:", userId)
				return cachedOrders, nil
			}

			fmt.Println("Cache unmarshal error for user:", userId, err)
		} else {
			fmt.Println("Cache miss for user:", userId)
		}
	}

	filter := bson.M{"user_id": userId}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	if r.Cache != nil {
		data, err := json.Marshal(orders)
		if err == nil {
			err = r.Cache.Set(ctx, cacheKey, data, 10*time.Minute).Err()
			if err != nil {
				fmt.Println("Failed to set cache for user:", userId, err)
			} else {
				fmt.Println("Cache set for user:", userId)
			}
		}
	}

	return orders, nil
}

func (r *OrderDao) invalidateUserOrdersCache(ctx context.Context) error {
	iter := r.Cache.Scan(ctx, 0, "orders:user:*", 0).Iterator()
	for iter.Next(ctx) {
		err := r.Cache.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

func (r *OrderDao) List(ctx context.Context, limit int64, skip int64) ([]model.Order, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
