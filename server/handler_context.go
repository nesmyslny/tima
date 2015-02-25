package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HandlerContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	User           *User
	body           []byte
}

func NewHandlerContext(w http.ResponseWriter, r *http.Request) *HandlerContext {
	return &HandlerContext{w, r, nil, nil}
}

func (handlerContext *HandlerContext) getReqBody() ([]byte, error) {
	var err error
	if handlerContext.body == nil {
		handlerContext.body, err = ioutil.ReadAll(handlerContext.Request.Body)
	}
	return handlerContext.body, err
}

func (handlerContext *HandlerContext) GetReqBodyJSON(model interface{}) error {
	body, err := handlerContext.getReqBody()
	if err != nil {
		return errors.New("invalid request")
	}

	err = json.Unmarshal(body, model)
	if err != nil {
		return err
	}
	return nil
}

func (handlerContext *HandlerContext) GetRouteVarString(name string) (string, error) {
	vars := mux.Vars(handlerContext.Request)
	if val, ok := vars[name]; ok {
		return val, nil
	}
	return "", errors.New("invalid parameter: " + name)
}

func (handlerContext *HandlerContext) GetRouteVarInt(name string) (int, error) {
	str, err := handlerContext.GetRouteVarString(name)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(str, 0, 32)
	if err != nil {
		return 0, errors.New("invalid parameter: " + name)
	}
	return int(i), nil
}

func (handlerContext *HandlerContext) GetRouteVarTime(name string, layout string) (time.Time, error) {
	str, err := handlerContext.GetRouteVarString(name)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, errors.New("invalid parameter: " + name)
	}
	return t, nil
}
