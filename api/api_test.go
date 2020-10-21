package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/newestuser/faceit/user"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Running integration tests with httptest.Server -> http://developers--production.almamedia.fi.s3-website-eu-west-1.amazonaws.com/2014/painless-mongodb-testing-with-docker-and-golang/
// 1. Integration testing with make and docker -> https://blog.gojekengineering.com/golang-integration-testing-made-easy-a834e754fa4c
// 2. Integration testing with make and docker -> https://medium.com/@rabin_gaire/integration-test-on-golang-using-docker-852f4c933cbe

func TestFindCreatedUser(t *testing.T) {
	u := &user.User{
		FirstName: "John",
		LastName:  "Doe",
		NickName:  "johny",
		Password:  "qwerty",
		Email:     "john.doe@mail.com",
		Country:   "DE",
	}

	req := newRequest("POST", "/users", u)
	resp := doRequest(req)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	regUser := &user.User{}
	_ = json.NewDecoder(resp.Body).Decode(regUser)

	req = newRequest("GET", fmt.Sprintf("/users/%s", regUser.ID.Hex()), nil)
	resp = doRequest(req)

	actual := &user.User{}
	_ = json.NewDecoder(resp.Body).Decode(actual)

	assert.Equal(t, regUser, actual)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReturn400WhenUserInvalid(t *testing.T) {
	emptyName := ""
	u := &user.User{
		FirstName: emptyName,
		LastName:  "Doe",
		NickName:  "johny",
		Password:  "qwerty",
		Email:     "john.doe@mail.com",
		Country:   "DE",
	}

	req := newRequest("POST", "/users", u)
	resp := doRequest(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestReturn404WhenUserNotFound(t *testing.T) {
	uid := "non-existing-user"

	req := newRequest("GET", fmt.Sprintf("/users/%s", uid), nil)
	resp := doRequest(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func readRespBytes(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	return b
}

type inMemoryRepo struct {
	data map[primitive.ObjectID]*user.User
}

func newRepo() user.Repository {
	return &inMemoryRepo{
		data: make(map[primitive.ObjectID]*user.User),
	}
}

func (r *inMemoryRepo) Find(id string) (*user.User, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	return r.data[oid], nil
}

func (r *inMemoryRepo) Register(u *user.User) (*user.User, error) {
	id := primitive.NewObjectID()
	r.data[id] = u
	uCopy := *u
	uCopy.ID = id
	return &uCopy, nil
}

func doRequest(req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func newRequest(method, path string, body interface{}) *http.Request {
	addr := os.Getenv("FORM3_ADDR")
	if len(addr) == 0 {
		addr = "http://localhost:8080"
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", addr, path), newReader(body))
	if err != nil {
		panic(err)
	}

	return req
}

func newReader(val interface{}) io.Reader {
	if val == nil {
		return nil
	}
	b, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}
