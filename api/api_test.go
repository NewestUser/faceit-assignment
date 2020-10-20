package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/newestuser/faceit/user"
	"github.com/stretchr/testify/assert"
)

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
	u.ID, _ = primitive.ObjectIDFromHex(string(readRespBytes(resp)))

	req = newRequest("GET", fmt.Sprintf("/users/%s", u.ID.Hex()), nil)
	resp = doRequest(req)

	actual := &user.User{}
	_ = json.NewDecoder(resp.Body).Decode(actual)

	assert.Equal(t, u, actual)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func readRespBytes(resp *http.Response) []byte {
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err = resp.Body.Close(); err != nil {
		panic(err)
	}
	return bytes
}

type inMemoryRepo struct {
	data map[string]*user.User
}

func newRepo() user.Repository {
	return &inMemoryRepo{
		data: make(map[string]*user.User),
	}
}

func (r *inMemoryRepo) Find(id string) (*user.User, error) {

	return r.data[id], nil
}

func (r *inMemoryRepo) Register(u *user.User) (string, error) {
	id, _ := uuid.NewRandom()
	r.data[id.String()] = u
	return id.String(), nil
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
