package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/newestuser/faceit/user"
	"gopkg.in/go-playground/validator.v9"
)

func UserRegHandler(validate *validator.Validate, repo user.Repository) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) *StatusError {
		userReq := &user.User{}
		if err := validUnmarshal(userReq, r, validate); err != nil {
			return err
		}

		id := repo.Register(userReq)
		w.WriteHeader(http.StatusCreated)

		if _, err := w.Write([]byte(id)); err != nil {
			return &StatusError{Code: http.StatusInternalServerError, Err: err}
		}
		return nil
	}
}

func UserGetHandler(repo user.Repository) ReqHandler {
	return func(writer http.ResponseWriter, request *http.Request) *StatusError {

		return nil
	}
}

func validUnmarshal(val interface{}, req *http.Request, validate *validator.Validate) *StatusError {
	if err := json.NewDecoder(req.Body).Decode(val); err != nil {
		return &StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := validate.Struct(validate); err != nil {
		return &StatusError{Code: http.StatusBadRequest, Err: err}
	}

	return nil
}

type ReqHandler func(http.ResponseWriter, *http.Request) *StatusError

func (fn ReqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil { // e is *appError, not os.Error.
		log.Printf("HTTP %d - %s", err.Status(), err)
		http.Error(w, err.Error(), err.Status())
		//http.Error(w, http.StatusText(http.StatusInternalServerError),
		//	http.StatusInternalServerError)
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
