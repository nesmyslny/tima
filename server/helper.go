package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func jsonResult(boolResult bool, stringResult string) (interface{}, *CtrlHandlerError) {
	return JsonResult{boolResult, stringResult}, nil
}

func jsonResultBool(boolResult bool) (interface{}, *CtrlHandlerError) {
	return JsonResult{BoolResult: boolResult}, nil
}

func jsonResultString(stringResult string) (interface{}, *CtrlHandlerError) {
	return JsonResult{StringResult: stringResult}, nil
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
