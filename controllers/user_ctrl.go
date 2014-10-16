package controllers

import "net/http"

type UserController struct{}

func (this *UserController) Signin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "invalid user/password", http.StatusBadRequest)
}
