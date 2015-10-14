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
	HandlerFunc func(context *HandlerContext) (interface{}, *HandlerError)
}

type AuthHandler struct {
	HandlerFunc      func(context *HandlerContext) (interface{}, *HandlerError)
	AuthenticateFunc func(context *HandlerContext) (bool, string, error)
	AuthorizeFunc    func(context *HandlerContext) (bool, error)
}

func NewAnonHandler(handlerFunc func(context *HandlerContext) (interface{}, *HandlerError)) AnonHandler {
	return AnonHandler{handlerFunc}
}

func NewAuthHandler(
	handlerFunc func(context *HandlerContext) (interface{}, *HandlerError),
	authenticateFunc func(context *HandlerContext) (bool, string, error),
	authorizeFunc func(context *HandlerContext) (bool, error)) AuthHandler {
	return AuthHandler{handlerFunc, authenticateFunc, authorizeFunc}
}

func (anonHandler AnonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewHandlerContext(w, r)
	response, hErr := anonHandler.HandlerFunc(context)
	serveHTTP(w, response, hErr)
}

func (authHandler AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := NewHandlerContext(w, r)

	authenticated, newToken, err := authHandler.AuthenticateFunc(context)
	if err != nil {
		http.Error(w, "Error while authenticating the user.", http.StatusInternalServerError)
		return
	}
	if !authenticated {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	authorized, err := authHandler.AuthorizeFunc(context)

	if err != nil {
		http.Error(w, "Error while authorize the user.", http.StatusInternalServerError)
		return
	}
	if !authorized {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	response, hErr := authHandler.HandlerFunc(context)
	w.Header().Set("Authorization", newToken)
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
