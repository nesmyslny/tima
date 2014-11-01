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
	migrationApi := server.NewMigrationApi(db, userApi)
	activitiesApi := server.NewActivitiesApi(db)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading)
	router.Handle("/upgrade", server.NewAnonHandler(migrationApi.UpgradeHandler)).Methods("POST")

	router.Handle("/signin", server.NewAnonHandler(userApi.SigninHandler)).Methods("POST")
	router.Handle("/issignedin", server.NewAnonHandler(userApi.IsSignedInHandler)).Methods("GET")

	router.Handle("/activities/{day}", server.NewAuthHandler(activitiesApi.GetByDayHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activities", server.NewAuthHandler(activitiesApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activities/{id}", server.NewAuthHandler(activitiesApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
