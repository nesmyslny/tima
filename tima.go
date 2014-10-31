package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/controllers"
	"github.com/nesmyslny/tima/dbaccess"
	"github.com/nesmyslny/tima/services"
)

func main() {
	// todo: configuration
	db := DbAccess.New("root:pwd@tcp(localhost:3307)/tima?parseTime=true")
	defer db.Close()

	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	activitiesService := services.NewActivitiesService(db)
	migrationService := services.NewMigrationService(db, userService)
	migrationController := controllers.NewMigrationController(migrationService)
	userController := controllers.NewUserController(authService, userService)
	activitiesController := controllers.NewActivitiesController(activitiesService)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading)
	router.Handle("/upgrade", controllers.NewAnonHandler(migrationController.Upgrade)).Methods("POST")

	router.Handle("/signin", controllers.NewAnonHandler(userController.Signin)).Methods("POST")
	router.Handle("/issignedin", controllers.NewAnonHandler(userController.IsSignedIn)).Methods("GET")

	router.Handle("/activities/{day}", controllers.NewAuthHandler(activitiesController.GetActivities, authService.AuthenticateRequest)).Methods("GET")
	router.Handle("/activities", controllers.NewAuthHandler(activitiesController.SaveActivity, authService.AuthenticateRequest)).Methods("POST")
	router.Handle("/activities/{id}", controllers.NewAuthHandler(activitiesController.DeleteActivity, authService.AuthenticateRequest)).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func blub(user string) error {
	return nil
}
