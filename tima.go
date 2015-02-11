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
	projectApi := server.NewProjectApi(db)
	activityTypeApi := server.NewActivityTypeApi(db)
	activityApi := server.NewActivityApi(db)
	migrationApi := server.NewMigrationApi(db, userApi)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading
	router.Handle("/upgrade", server.NewAnonHandler(migrationApi.UpgradeHandler)).Methods("POST")

	router.Handle("/signin", server.NewAnonHandler(userApi.SigninHandler)).Methods("POST")
	router.Handle("/issignedin", server.NewAnonHandler(userApi.IsSignedInHandler)).Methods("GET")

	router.Handle("/activities/{day}", server.NewAuthHandler(activityApi.GetByDayHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activities", server.NewAuthHandler(activityApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activities/{id}", server.NewAuthHandler(activityApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projects", server.NewAuthHandler(projectApi.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectApi.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects", server.NewAuthHandler(projectApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeApi.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeApi.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeApi.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeApi.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projectActivityTypes", server.NewAuthHandler(activityTypeApi.GetActivityViewListHandler, auth.AuthenticateRequest)).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
