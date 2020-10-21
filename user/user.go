package user

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotFoundError struct {
	UID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("user with id: %s not found", e.UID)
}

type Repository interface {
	Register(u *User) (*User, error)

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
