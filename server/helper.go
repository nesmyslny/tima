package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var errItemInUse = errors.New("Item is already in use")
var errUsernameUnavailable = errors.New("Username unavailable")

func jsonResult(boolResult bool, stringResult string, intResult int) (interface{}, *HandlerError) {
	return JsonResult{boolResult, stringResult, intResult}, nil
}

func jsonResultBool(boolResult bool) (interface{}, *HandlerError) {
	return JsonResult{BoolResult: boolResult}, nil
}

func jsonResultString(stringResult string) (interface{}, *HandlerError) {
	return JsonResult{StringResult: stringResult}, nil
}

func jsonResultInt(intResult int) (interface{}, *HandlerError) {
	return JsonResult{IntResult: intResult}, nil
}

func unmarshalJSON(body io.Reader, model interface{}) error {
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

func getRouteVarString(r *http.Request, name string) (string, error) {
	vars := mux.Vars(r)
	if val, ok := vars[name]; ok {
		return val, nil
	}
	return "", errors.New("invalid parameter: " + name)
}

func getRouteVarInt(r *http.Request, name string) (int, error) {
	str, err := getRouteVarString(r, name)
	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(str, 0, 32)
	if err != nil {
		return 0, errors.New("invalid parameter: " + name)
	}

	return int(i), nil
}

func getRouteVarTime(r *http.Request, name string, layout string) (time.Time, error) {
	str, err := getRouteVarString(r, name)
	if err != nil {
		return time.Time{}, err
	}

	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Time{}, errors.New("invalid parameter: " + name)
	}

	return t, nil
}
