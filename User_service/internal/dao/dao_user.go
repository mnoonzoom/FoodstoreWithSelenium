package dao

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"user/internal/model"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user model.User) (string, error) {

	res, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(interface {
		Hex() string
	}).Hex(), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return &user, err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ObjectID: %v", err)
	}

	var user model.User
	err = r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
