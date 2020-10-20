package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Register(u *User) (string, error)

	Find(id string) (*User, error)
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string             `json:"lastName" bson:"lastName" validate:"required"`
	NickName  string             `json:"nickName" bson:"nickName" validate:"required"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Country   string             `json:"country" bson:"country" validate:"required"`
}
