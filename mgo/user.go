package mgo

import (
	"context"
	"github.com/newestuser/faceit/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const users = "users"

type mgoUserRepository struct {
	*MongoDB
}

func (s *mgoUserRepository) users() *mongo.Collection {
	return s.Collection(users)
}

func (s *mgoUserRepository) Register(u *user.User) (string, error) {
	result, err := s.users().InsertOne(context.TODO(), u)
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

func (s *mgoUserRepository) Find(id string) (*user.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := s.users().FindOne(context.TODO(), bson.M{"_id": oid})

	if result.Err() != nil {
		return nil, result.Err()
	}

	u := &user.User{}
	err = result.Decode(u)
	return u, err
}

func NewUserRepository(db *MongoDB) user.Repository {
	return &mgoUserRepository{db}
}
