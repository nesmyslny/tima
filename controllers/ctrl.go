package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/models"
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
	HandlerFunc func(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError)
	AuthFunc    func(r *http.Request) (bool, *models.User)
}

func NewAnonHandler(handlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError)) anonHandler {
	return anonHandler{handlerFunc}
}

func NewAuthHandler(
	handlerFunc func(w http.ResponseWriter, r *http.Request, user *models.User) (interface{}, *CtrlHandlerError),
	authFunc func(r *http.Request) (bool, *models.User)) authHandler {
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

func jsonResult(boolResult bool, stringResult string) (interface{}, *CtrlHandlerError) {
	return models.JsonResult{boolResult, stringResult}, nil
}

func jsonResultBool(boolResult bool) (interface{}, *CtrlHandlerError) {
	return models.JsonResult{BoolResult: boolResult}, nil
}

func jsonResultString(stringResult string) (interface{}, *CtrlHandlerError) {
	return models.JsonResult{StringResult: stringResult}, nil
}

func unmarshalJson(body io.Reader, model interface{}) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return errors.New("invalid request")
	}

	err = json.Unmarshal(data, model)
	if err != nil {
		return errors.New("invalid data")
	}

	return nil
}

func getRouteVar(r *http.Request, name string) string {
	vars := mux.Vars(r)
	return vars[name]
}
