package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"parser/services"
)

// ContentHandler handles the GET request for file content
func ContentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uri := query.Get("uri")
	if uri == "" {
		http.Error(w, "Missing 'uri' query parameter", http.StatusBadRequest)
		return
	}

	token := os.Getenv("GIT_ACCESS_TOKEN")
	apiURL, err := services.ConvertGitHubURLToAPIURL(uri)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid GitHub URL: %v", err), http.StatusBadRequest)
		return
	}

	content, err := services.FetchFileContentFromURL(apiURL, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching content: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": content})
}

// Main function sets up the router and starts the server
func main() {
	router := mux.NewRouter()

	// Define the /content endpoint
	router.HandleFunc("/content", ContentHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s\n", port)
	http.ListenAndServe(":"+port, router)
}
