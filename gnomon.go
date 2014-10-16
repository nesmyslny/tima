package main

import (
	"github.com/gorilla/mux"
	"gnomon/services"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/signin", signin).Methods("POST")
	router.HandleFunc("/upgrade", upgrade).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "invalid user/password", http.StatusBadRequest)
}

func upgrade(w http.ResponseWriter, r *http.Request) {
	err := services.MigrationService.Run()
	if err != nil {
		// todo: logging
		// in this case, the internal error is directly exposed to the user.
		// upgrading is an admin task and the internal error is needed to resolve issues.
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
