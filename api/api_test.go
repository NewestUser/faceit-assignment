package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/newestuser/faceit/user"
	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/validator.v9"
)

func Test(t *testing.T) {
	u := userReader(&user.User{
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johny",
		Password:  "qwerty",
		Email:     "john.doe@mail.com",
		Country:   "DE",
	})

	req, _ := http.NewRequest("POST", "not-important", u)
	rec := httptest.NewRecorder()

	UserRegHandler(validator.New(), newRepo()).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
}

func userReader(u *user.User) io.Reader {
	b, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

type inMemoryRepo struct {
	data map[string]*user.User
}

func newRepo() user.Repository {
	return &inMemoryRepo{
		data: make(map[string]*user.User),
	}
}

func (r *inMemoryRepo) Find(id string) *user.User {

	return r.data[id]
}

func (r *inMemoryRepo) Register(u *user.User) string {
	id, _ := uuid.NewRandom()
	r.data[id.String()] = u
	return id.String()
}
