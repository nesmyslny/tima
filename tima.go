package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/server"
)

func main() {
	// todo: configuration
	db := server.NewDb("root:pwd@tcp(localhost:3307)/tima?parseTime=true")
	defer db.Close()

	auth := server.NewAuth()
	userApi := server.NewUserApi(db, auth)
	projectsApi := server.NewProjectsApi(db)
	activityTypesApi := server.NewActivityTypesApi(db)
	activitiesApi := server.NewActivitiesApi(db)
	migrationApi := server.NewMigrationApi(db, userApi)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading)
	router.Handle("/upgrade", server.NewAnonHandler(migrationApi.UpgradeHandler)).Methods("POST")

	router.Handle("/signin", server.NewAnonHandler(userApi.SigninHandler)).Methods("POST")
	router.Handle("/issignedin", server.NewAnonHandler(userApi.IsSignedInHandler)).Methods("GET")

	router.Handle("/activities/{day}", server.NewAuthHandler(activitiesApi.GetByDayHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activities", server.NewAuthHandler(activitiesApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activities/{id}", server.NewAuthHandler(activitiesApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projects", server.NewAuthHandler(projectsApi.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectsApi.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects", server.NewAuthHandler(projectsApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectsApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/activityTypes", server.NewAuthHandler(activityTypesApi.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypesApi.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes", server.NewAuthHandler(activityTypesApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypesApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
