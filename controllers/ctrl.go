package controllers

import (
	"encoding/json"
	"errors"
	"github.com/nesmyslny/tima/models"
	"io"
	"io/ioutil"
	"net/http"
)

type CtrlHandlerError struct {
	Error   error
	Message string
	Code    int
}

type CtrlHandlerStruct struct {
	HandlerFunc func(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError)
	AuthFunc    func(r *http.Request) bool
}

func (this CtrlHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authorized := true
	if this.AuthFunc != nil {
		authorized = this.AuthFunc(r)
	}

	if !authorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, hErr := this.HandlerFunc(w, r)

	if hErr != nil {
		http.Error(w, hErr.Message, hErr.Code)
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
