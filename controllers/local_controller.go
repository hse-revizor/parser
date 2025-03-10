package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"parser/services"
)

type LocalController struct{}

func NewLocalController() *LocalController {
	return &LocalController{}
}

func (lc *LocalController) ContentHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query.Get("path")
	if path == "" {
		http.Error(w, "Missing 'path' query parameter", http.StatusBadRequest)
		return
	}

	rowParam := query.Get("row")
	startRow, endRow, err := services.ParseRowParam(rowParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid row parameter: %v", err), http.StatusBadRequest)
		return
	}
	parsedPath, err := url.Parse(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid path: %v", err), http.StatusBadRequest)
		return
	}

	parsedPath.RawQuery = ""
	cleanPath := parsedPath.Path

	filePath, err := services.ConvertLocalPathToFilePath(cleanPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid local path: %v", err), http.StatusBadRequest)
		return
	}

	content, err := services.FetchFileContentFromLocalPath(filePath)
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
