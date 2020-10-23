package user

import (
	"fmt"
)

type Event string

const CreateEvent = Event("create")
const UpdateEvent = Event("update")

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

type EventEmitter interface {
	// Emit an event that a User has been updated.
	EmitUpdate(e Event, u *User) error
}

type User struct {
	ID        string `json:"id"`
	FirstName string              `json:"firstName" validate:"required"`
	LastName  string              `json:"lastName" validate:"required"`
	NickName  string              `json:"nickName" validate:"required"`
	Password  string              `json:"password" validate:"required"`
	Email     string              `json:"email" validate:"required,email"`
	Country   string              `json:"country" validate:"required"`
}
