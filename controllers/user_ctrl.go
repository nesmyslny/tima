package controllers

import (
	"encoding/json"
	"errors"
	"gnomon/models"
	"gnomon/services"
	"io"
	"io/ioutil"
	"net/http"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService}
}

func (this *UserController) Signin(w http.ResponseWriter, r *http.Request) {
	var formData models.UserSignin
	err := unmarshalRequestData(r.Body, &formData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !this.userService.Authenticate(formData.Username, formData.Password) {
		http.Error(w, "invalid user/password", http.StatusBadRequest)
	}
}

func unmarshalRequestData(body io.Reader, model interface{}) error {
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
