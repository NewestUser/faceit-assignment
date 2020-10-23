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

func NewUserRepository(db *MongoDB, em user.EventEmitter) user.Repository {
	return &mgoUserRepository{MongoDB: db, eventEmitter: em}
}

type mgoUserRepository struct {
	*MongoDB
	eventEmitter user.EventEmitter
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
	result, err := s.users().InsertOne(context.TODO(), adaptToBson(u))
	if err != nil {
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)
	uCopy.ID = id.Hex()

	return &uCopy, nil
}

func (s *mgoUserRepository) Update(u *user.User) (*user.User, error) {
	if u.ID == "" {
		return nil, &user.ErrNotFound{UID: "'nil'"}
	}

	oid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return nil, &user.ErrNotFound{UID: u.ID}
	}

	result, err := s.users().UpdateOne(context.TODO(), bson.M{"_id": oid}, bson.M{"$set": adaptToBson(u)})
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, &user.ErrNotFound{UID: u.ID}
	}
	// since I don't have another a service layer and I don't have a use cae for
	// how to handle errors when emitting events I am deliberately ignoring the error
	err = s.eventEmitter.EmitUpdate(user.UpdateEvent, u)
	println("failed emitting User Update Event", err)

	return u, nil
}

func (s *mgoUserRepository) users() *mongo.Collection {
	return s.Collection(users)
}

func adaptToBson(u *user.User) bson.M {
	return bson.M{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"nickName":  u.NickName,
		"password":  u.Password,
		"email":     u.Email,
		"country":   u.Country,
	}
}
