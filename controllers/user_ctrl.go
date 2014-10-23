package controllers

import (
	"gnomon/models"
	"gnomon/services"
	"net/http"
	"time"
)

type UserController struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewUserController(authService *services.AuthService, userService *services.UserService) *UserController {
	return &UserController{authService, userService}
}

func (this *UserController) Signin(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	var formData models.UserSignin
	err := unmarshalJson(r.Body, &formData)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}

	token, err := this.authService.Authenticate(formData.Username, formData.Password)
	if err != nil {
		return nil, &CtrlHandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	return jsonResultString(token)
}

func (this *UserController) IsSignedIn(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	signedIn := this.authService.ValidateToken(r)
	return jsonResultBool(signedIn)
}

func (this *UserController) Secret(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	return jsonResult(true, time.Now().String())
}
