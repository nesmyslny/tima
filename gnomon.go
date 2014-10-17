package main

import (
	"github.com/gorilla/mux"
	"gnomon/controllers"
	"gnomon/dbaccess"
	"gnomon/services"
	"net/http"
)

func main() {
	// todo: configuration
	db := DbAccess.New("root:pwd@tcp(localhost:3307)/gnomon?parseTime=true")
	defer db.Close()

	migrationService := services.NewMigrationService(db)
	migrationController := controllers.NewMigrationController(migrationService)
	userController := controllers.UserController{}

	router := mux.NewRouter()
	router.HandleFunc("/signin", userController.Signin).Methods("POST")
	router.HandleFunc("/upgrade", migrationController.Upgrade).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
