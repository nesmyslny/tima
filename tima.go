package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/server"
)

func main() {
	if execFlags() {
		return
	}

	db := initDB()
	defer db.Close()

	auth := server.NewAuth()
	userAPI := server.NewUserAPI(db, auth)
	projectAPI := server.NewProjectAPI(db)
	projectCategoryAPI := server.NewProjectCategoryAPI(db)
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

	router.Handle("/projectCategories/tree", server.NewAuthHandler(projectCategoryAPI.GetTreeHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projectCategories/list", server.NewAuthHandler(projectCategoryAPI.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/projectCategories", server.NewAuthHandler(projectCategoryAPI.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/projectCategories/{id}", server.NewAuthHandler(projectCategoryAPI.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeAPI.GetListHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeAPI.GetHandler, auth.AuthenticateRequest)).Methods("GET")
	router.Handle("/activityTypes", server.NewAuthHandler(activityTypeAPI.SaveHandler, auth.AuthenticateRequest)).Methods("POST")
	router.Handle("/activityTypes/{id}", server.NewAuthHandler(activityTypeAPI.DeleteHandler, auth.AuthenticateRequest)).Methods("DELETE")

	router.Handle("/projectActivityTypes", server.NewAuthHandler(activityTypeAPI.GetActivityViewListHandler, auth.AuthenticateRequest)).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func initDB() *server.DB {
	// todo: configuration
	db := server.NewDB("root:pwd@tcp(localhost:3307)/tima?parseTime=true")
	return db
}

func execFlags() bool {
	dbUp := flag.Int("dbUp", -1, "Applies the given number of database migrations (if 0 is specified, all pending migrations will be applied)")
	dbDown := flag.Int("dbDown", -1, "Undos the given number of database migration (if 0 is specified, all migration will be undone)")
	dbGenerateData := flag.Bool("dbGenerateData", false, "Generate test data")
	flag.Parse()

	if *dbUp > -1 && *dbDown > -1 {
		fmt.Println("dbUp and dbDown can't be used at once")
		return true
	}

	db := initDB()
	defer db.Close()

	if *dbUp > -1 {
		err := db.Upgrade(*dbUp)
		printCliActionResult(err)
		return true
	}

	if *dbDown > -1 {
		err := db.Downgrade(*dbDown)
		printCliActionResult(err)
		return true
	}

	if *dbGenerateData {
		return true
	}

	return false
}

func printCliActionResult(err error) {
	if err == nil {
		fmt.Println("done!")
	} else {
		fmt.Println(err.Error())
	}
}
