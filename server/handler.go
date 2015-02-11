package server

import (
	"encoding/json"
	"net/http"
)

type HandlerError struct {
	Error   error
	Message string
	Code    int
}

type AnonHandler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *HandlerError)
}

type AuthHandler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError)
	AuthFunc    func(r *http.Request) (bool, *User)
}

func NewAnonHandler(handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *HandlerError)) AnonHandler {
	return AnonHandler{handlerFunc}
}

func NewAuthHandler(
	handlerFunc func(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *HandlerError),
	authFunc func(r *http.Request) (bool, *User)) AuthHandler {
	return AuthHandler{handlerFunc, authFunc}
}

func (anonHandler AnonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, hErr := anonHandler.HandlerFunc(w, r)
	serveHTTP(w, response, hErr)
}

func (authHandler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorized, user := authHandler.AuthFunc(r)

	if !authorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, hErr := authHandler.HandlerFunc(w, r, user)
	serveHTTP(w, response, hErr)
}

func serveHTTP(w http.ResponseWriter, response interface{}, handlerError *HandlerError) {
	if handlerError != nil {
		http.Error(w, handlerError.Message, handlerError.Code)
		return
	}

	if response == nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}
