package mgo

import (
	"context"
	"errors"
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

func (s *mgoUserRepository) Register(u *user.User) (*user.User, error) {
	result, err := s.users().InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)

	uCopy := *u
	uCopy.ID = id

	return &uCopy, nil
}

func (s *mgoUserRepository) Find(id string) (*user.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &user.NotFoundError{UID: id}
	}

	result := s.users().FindOne(context.TODO(), bson.M{"_id": oid})

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, &user.NotFoundError{UID: id}
	} else if result.Err() != nil {
		return nil, err
	}

	u := &user.User{}
	err = result.Decode(u)
	return u, err
}

func NewUserRepository(db *MongoDB) user.Repository {
	return &mgoUserRepository{db}
}
