package main

import (
	"github.com/gorilla/mux"
	"gnomon/controllers"
	"gnomon/services"
	"net/http"
)

func main() {
	migrationController := controllers.MigrationController{&services.MigrationService{}}
	userController := controllers.UserController{}

	router := mux.NewRouter()
	router.HandleFunc("/signin", userController.Signin).Methods("POST")
	router.HandleFunc("/upgrade", migrationController.Upgrade).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
