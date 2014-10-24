package main

import (
	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/controllers"
	"github.com/nesmyslny/tima/dbaccess"
	"github.com/nesmyslny/tima/services"
	"net/http"
)

func main() {
	// todo: configuration
	db := DbAccess.New("root:pwd@tcp(localhost:3307)/tima?parseTime=true")
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
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
