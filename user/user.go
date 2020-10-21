package user

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ErrNotFound struct {
	UID string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("user with id: %s not found", e.UID)
}

type Repository interface {
	// Search for a user by ID. If the user does not exist return ErrNotFound.
	Find(id string) (*User, error)

	// Register a new user and return the registered user with his assigned ID.
	Register(u *User) (*User, error)

	// Perform an update of an existing user and return the updated one.
	// If the user does not exist return ErrNotFound.
	Update(u *User) (*User, error)
}

type User struct {
	ID        *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName string             `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string             `json:"lastName" bson:"lastName" validate:"required"`
	NickName  string             `json:"nickName" bson:"nickName" validate:"required"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Country   string             `json:"country" bson:"country" validate:"required"`
}
