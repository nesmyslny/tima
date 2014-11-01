package controllers

import (
	"net/http"

	"github.com/nesmyslny/tima/models"
	"github.com/nesmyslny/tima/services"
)

type UserController struct {
	authService *services.AuthService
	userService *services.UserService
}

func NewUserController(authService *services.AuthService, userService *services.UserService) *UserController {
	return &UserController{authService, userService}
}

func (this *UserController) Signin(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	var credentials models.UserCredentials
	err := unmarshalJson(r.Body, &credentials)
	if err != nil {
		return nil, &CtrlHandlerError{err, err.Error(), http.StatusBadRequest}
	}

	token, err := this.authService.Authenticate(credentials.Username, credentials.Password)
	if err != nil {
		return nil, &CtrlHandlerError{err, "Invalid username/password", http.StatusBadRequest}
	}

	return jsonResultString(token)
}

func (this *UserController) IsSignedIn(w http.ResponseWriter, r *http.Request) (interface{}, *CtrlHandlerError) {
	signedIn := this.authService.ValidateToken(r)
	return jsonResultBool(signedIn)
}
