package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"parser/services"
)

type GitController struct{}

func NewGitController() *GitController {
	return &GitController{}
}

func (gc *GitController) ContentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	uri := query.Get("uri")
	if uri == "" {
		http.Error(w, "Missing 'uri' query parameter", http.StatusBadRequest)
		return
	}

	rowParam := query.Get("row")
	startRow, endRow, err := services.ParseRowParam(rowParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid row parameter: %v", err), http.StatusBadRequest)
		return
	}

	parsedURI, err := url.Parse(uri)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid URI: %v", err), http.StatusBadRequest)
		return
	}

	parsedURI.RawQuery = ""
	cleanURI := parsedURI.String()

	token := os.Getenv("GIT_ACCESS_TOKEN")
	apiURL, err := services.ConvertGitHubURLToAPIURL(cleanURI)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid GitHub URL: %v", err), http.StatusBadRequest)
		return
	}

	content, err := services.FetchFileContentFromURL(apiURL, token)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching content: %v", err), http.StatusInternalServerError)
		return
	}

	if startRow > 0 {
		if startRow == endRow {
			content, err = services.ExtractRow(content, startRow)
		} else {
			content, err = services.ExtractRowRange(content, startRow, endRow)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Error extracting row(s): %v", err), http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"content": content})
}
