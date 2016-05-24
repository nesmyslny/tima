package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nesmyslny/tima/server"
)

type cliFlags struct {
	configPath     string
	dbUp           int
	dbDown         int
	dbGenerateData bool
}

func main() {
	flags := parseFlags()
	config := server.NewConfig(flags.configPath)

	db := server.NewDB(config.GetDataSourceName())
	defer db.Close()

	executed, err := db.ExecCliActions(flags.dbUp, flags.dbDown, flags.dbGenerateData)
	if executed {
		printCliActionResult(err)
		return
	}

	auth := server.NewAuth()
	departmentAPI := server.NewDepartmentAPI(db)
	userAPI := server.NewUserAPI(db, auth)
	projectAPI := server.NewProjectAPI(db)
	projectCategoryAPI := server.NewProjectCategoryAPI(db)
	activityTypeAPI := server.NewActivityTypeAPI(db)
	activityAPI := server.NewActivityAPI(db)
	migrationAPI := server.NewMigrationAPI(db, userAPI)
	reporting := server.NewReporting(db)

	router := mux.NewRouter()

	// todo: secure upgrade route (-> implement installation/upgrading
	createAnonRoute(router, "/upgrade", "POST", migrationAPI.UpgradeHandler)

	createAnonRoute(router, "/signIn", "POST", userAPI.SigninHandler)
	createAnonRoute(router, "/isSignedIn", "GET", userAPI.IsSignedInHandler)

	createAuthRoute(router, auth, server.AuthorizeUser, "/activities/{day}", "GET", activityAPI.GetByDayHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/activities", "POST", activityAPI.SaveHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/activities/{id}", "DELETE", activityAPI.DeleteHandler)

	createAuthRoute(router, auth, server.AuthorizeDeptManager, "/projects/admin", "GET", projectAPI.GetListHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/projects/user", "GET", projectAPI.GetListUserHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/projects/selectList", "GET", projectAPI.GetListSelectHandler)
	createAuthRoute(router, auth, projectAPI.AuthorizeGet, "/projects/{id}", "GET", projectAPI.GetHandler)
	createAuthRoute(router, auth, projectAPI.AuthorizeSave, "/projects", "POST", projectAPI.SaveHandler)
	createAuthRoute(router, auth, projectAPI.AuthorizeDelete, "/projects/{id}", "DELETE", projectAPI.DeleteHandler)

	createAuthRoute(router, auth, server.AuthorizeManager, "/projectCategories/tree", "GET", projectCategoryAPI.GetTreeHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/projectCategories/list", "GET", projectCategoryAPI.GetListHandler)
	createAuthRoute(router, auth, server.AuthorizeManager, "/projectCategories", "POST", projectCategoryAPI.SaveHandler)
	createAuthRoute(router, auth, server.AuthorizeManager, "/projectCategories/{id}", "DELETE", projectCategoryAPI.DeleteHandler)

	createAuthRoute(router, auth, server.AuthorizeUser, "/activityTypes", "GET", activityTypeAPI.GetListHandler)
	createAuthRoute(router, auth, server.AuthorizeManager, "/activityTypes/{id}", "GET", activityTypeAPI.GetHandler)
	createAuthRoute(router, auth, server.AuthorizeManager, "/activityTypes", "POST", activityTypeAPI.SaveHandler)
	createAuthRoute(router, auth, server.AuthorizeManager, "/activityTypes/{id}", "DELETE", activityTypeAPI.DeleteHandler)

	createAuthRoute(router, auth, server.AuthorizeUser, "/users/all", "GET", userAPI.GetListHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/users/department", "GET", userAPI.GetListDeptHandler)
	createAuthRoute(router, auth, userAPI.AuthorizeGet, "/users/{id}", "GET", userAPI.GetHandler)
	createAuthRoute(router, auth, userAPI.AuthorizeSave, "/users", "POST", userAPI.SaveHandler)

	createAuthRoute(router, auth, server.AuthorizeAdmin, "/departments/tree", "GET", departmentAPI.GetTreeHandler)
	createAuthRoute(router, auth, server.AuthorizeUser, "/departments/list", "GET", departmentAPI.GetListHandler)
	createAuthRoute(router, auth, server.AuthorizeAdmin, "/departments", "POST", departmentAPI.SaveHandler)
	createAuthRoute(router, auth, server.AuthorizeAdmin, "/departments/{id}", "DELETE", departmentAPI.DeleteHandler)

	createAuthRoute(router, auth, server.AuthorizeManager, "/report/overview", "POST", reporting.CreateReportOverview)
	createAuthRoute(router, auth, server.AuthorizeManager, "/report/projects", "POST", reporting.CreateReportProjects)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public/")))
	http.Handle("/", router)
	http.ListenAndServe(":"+config.HttpPort, nil)
}

func createAnonRoute(router *mux.Router, path string, method string,
	handlerFunc func(context *server.HandlerContext) (interface{}, *server.HandlerError)) {
	router.Handle(path, server.NewAnonHandler(handlerFunc)).Methods(method)
}

func createAuthRoute(router *mux.Router, auth *server.Auth,
	authorizeFunc func(context *server.HandlerContext) (bool, error),
	path string, method string,
	handlerFunc func(context *server.HandlerContext) (interface{}, *server.HandlerError)) {
	router.Handle(path, server.NewAuthHandler(handlerFunc, auth.AuthenticateUser, authorizeFunc)).Methods(method)
}

func parseFlags() cliFlags {
	configPath := flag.String("configPath", ".", "Path of the config file (default: current working directory)")
	dbUp := flag.Int("dbUp", -1, "Apply the given number of database migrations (if 0 is specified, all pending migrations will be applied)")
	dbDown := flag.Int("dbDown", -1, "Undo the given number of database migration (if 0 is specified, all migration will be undone)")
	dbGenerateData := flag.Bool("dbGenerateData", false, "Generate test data")
	flag.Parse()

	return cliFlags{*configPath, *dbUp, *dbDown, *dbGenerateData}
}

func printCliActionResult(err error) {
	if err == nil {
		fmt.Println("done!")
	} else {
		fmt.Println(err.Error())
	}
}
