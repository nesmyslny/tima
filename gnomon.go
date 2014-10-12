package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/signin", signin).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("web/")))
	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func signin(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "invalid user/password", http.StatusBadRequest)
}
