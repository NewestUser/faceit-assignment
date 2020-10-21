package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/newestuser/faceit/user"
	"gopkg.in/go-playground/validator.v9"
)

func UserRegHandler(validate *validator.Validate, repo user.Repository) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) *StatusError {
		userReq := &user.User{}
		if err := validUnmarshal(userReq, r, validate); err != nil {
			return err
		}

		regUser, err := repo.Register(userReq)
		if err != nil {
			return &StatusError{Code: http.StatusInternalServerError, Err: err}
		}

		return respond(w, regUser, http.StatusCreated)
	}
}

func UserUpdateHandler(validate *validator.Validate, repo user.Repository) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) *StatusError {
		userReq := &user.User{}
		if err := validUnmarshal(userReq, r, validate); err != nil {
			return err
		}

		regUser, err := repo.Update(userReq)
		if err != nil {
			if _, ok := err.(*user.ErrNotFound); ok {
				return &StatusError{Code: http.StatusNotFound, Err: err}
			}

			return &StatusError{Code: http.StatusInternalServerError, Err: err}
		}

		return respond(w, regUser, http.StatusOK)
	}
}

func UserGetHandler(repo user.Repository) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) *StatusError {
		userId := mux.Vars(r)["id"]
		foundUser, err := repo.Find(userId)

		if err != nil {
			if _, ok := err.(*user.ErrNotFound); ok {
				return &StatusError{Code: http.StatusNotFound, Err: err}
			}

			return &StatusError{Code: http.StatusInternalServerError, Err: err}
		}

		return respond(w, foundUser, http.StatusOK)
	}
}

func respond(w http.ResponseWriter, val interface{}, status int) *StatusError {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		return &StatusError{Code: http.StatusInternalServerError, Err: err}
	}
	return nil
}

func validUnmarshal(val interface{}, req *http.Request, validate *validator.Validate) *StatusError {
	if err := json.NewDecoder(req.Body).Decode(val); err != nil {
		return &StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := validate.Struct(val); err != nil {
		return &StatusError{Code: http.StatusBadRequest, Err: err}
	}

	return nil
}

type ReqHandler func(http.ResponseWriter, *http.Request) *StatusError

// https://blog.questionable.services/article/http-handler-error-handling-revisited/
func (fn ReqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("HTTP %d - %s", err.Status(), err)
		http.Error(w, err.Error(), err.Status())
	}
}

type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

func (se StatusError) Status() int {
	return se.Code
}
