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

func (s *mgoUserRepository) Find(id string) (*user.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &user.ErrNotFound{UID: id}
	}

	result := s.users().FindOne(context.TODO(), bson.M{"_id": oid})

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, &user.ErrNotFound{UID: id}
	} else if result.Err() != nil {
		return nil, err
	}

	u := &user.User{}
	err = result.Decode(u)
	return u, err
}

func (s *mgoUserRepository) Register(u *user.User) (*user.User, error) {
	uCopy := *u
	uCopy.ID = nil

	result, err := s.users().InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	uCopy.ID = &id

	return &uCopy, nil
}

func (s *mgoUserRepository) Update(u *user.User) (*user.User, error) {
	if u.ID == nil {
		return nil, &user.ErrNotFound{UID: "'nill'"}
	}

	result, err := s.users().UpdateOne(context.TODO(), bson.M{"_id": u.ID},
		bson.M{
			"$set": bson.M{
				"firstName": u.FirstName,
				"lastName":  u.LastName,
				"nickName":  u.NickName,
				"password":  u.Password,
				"email":     u.Email,
				"country":   u.Country,
			},
		})
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, &user.ErrNotFound{UID: u.ID.Hex()}
	}

	return u, nil
}

func (s *mgoUserRepository) users() *mongo.Collection {
	return s.Collection(users)
}

func NewUserRepository(db *MongoDB) user.Repository {
	return &mgoUserRepository{db}
}
