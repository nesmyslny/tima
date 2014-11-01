package server

import (
	"encoding/json"
	"net/http"
)

type CtrlHandlerError struct {
	Error   error
	Message string
	Code    int
}

type anonHandler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError)
}

type authHandler struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *CtrlHandlerError)
	AuthFunc    func(r *http.Request) (bool, *User)
}

func NewAnonHandler(handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError)) anonHandler {
	return anonHandler{handlerFunc}
}

func NewAuthHandler(
	handlerFunc func(w http.ResponseWriter, r *http.Request, user *User) (interface{}, *CtrlHandlerError),
	authFunc func(r *http.Request) (bool, *User)) authHandler {
	return authHandler{handlerFunc, authFunc}
}

func (this anonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, hErr := this.HandlerFunc(w, r)
	serveHTTP(w, response, hErr)
}

func (this authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorized, user := this.AuthFunc(r)

	if !authorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, hErr := this.HandlerFunc(w, r, user)
	serveHTTP(w, response, hErr)
}

func serveHTTP(w http.ResponseWriter, response interface{}, handlerError *CtrlHandlerError) {
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
