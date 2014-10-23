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

	authService := services.NewAuthService(db)
	authService = services.NewAuthService(db)
	userService := services.NewUserService(db)
	migrationService := services.NewMigrationService(db, userService)
	migrationController := controllers.NewMigrationController(migrationService)
	userController := controllers.NewUserController(authService, userService)

	router := mux.NewRouter()
	router.Handle("/signin", controllers.CtrlHandlerStruct{userController.Signin, nil}).Methods("POST")
	router.Handle("/issignedin", controllers.CtrlHandlerStruct{userController.IsSignedIn, nil}).Methods("GET")
	router.Handle("/secret", controllers.CtrlHandlerStruct{userController.Secret, authService.ValidateToken}).Methods("GET")
	router.Handle("/upgrade", controllers.CtrlHandlerStruct{migrationController.Upgrade, nil}).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
