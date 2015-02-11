package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/server"
)

func main() {
	// todo: configuration
	db := server.NewDB("root:pwd@tcp(localhost:3307)/tima?parseTime=true")
	defer db.Close()

	auth := server.NewAuth()
	userAPI := server.NewUserAPI(db, auth)
	projectAPI := server.NewProjectAPI(db)
	activityTypeAPI := server.NewActivityTypeAPI(db)
	activityAPI := server.NewActivityAPI(db)
	migrationAPI := server.NewMigrationAPI(db, userAPI)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading
	router.Handle("/upgrade", server.NewAnonHandler(migrationAPI.UpgradeHandler)).Methods("POST")

	router.Handle("/signin", server.NewAnonHandler(userAPI.SigninHandler)).Methods("POST")
	router.Handle("/issignedin", server.NewAnonHandler(userAPI.IsSignedInHandler)).Methods("GET")

	router.Handle("/activities/{day}", server.NewAuthHandler(activityAPI.GetByDayHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activities", server.NewAuthHandler(activityAPI.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activities/{id}", server.NewAuthHandler(activityAPI.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projects", server.NewAuthHandler(projectAPI.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectAPI.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projects", server.NewAuthHandler(projectAPI.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/projects/{id}", server.NewAuthHandler(projectAPI.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeAPI.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeAPI.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeAPI.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeAPI.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projectActivityTypes", server.NewAuthHandler(activityTypeAPI.GetActivityViewListHandler, auth.AuthenticateRequest)).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
