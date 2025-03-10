package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"parser/controllers"
)

func main() {
	router := mux.NewRouter()

	gitController := controllers.NewGitController()
	localController := controllers.NewLocalController()

	router.HandleFunc("/github-content", gitController.ContentHandler).Methods("GET")
	router.HandleFunc("/local-content", localController.ContentHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, router)
}
