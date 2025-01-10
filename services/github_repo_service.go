package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GitHubFile defines github file structure
type GitHubFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

// ListRepFiles returns list of all files in repository
func ListRepFiles(owner, repo, token string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get contents: %s", resp.Status)
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		if file.Type == "file" {
			fileList = append(fileList, file.Path)
		} else if file.Type == "dir" {
			subFiles, err := ListFilesInDir(owner, repo, token, file.Path)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, subFiles...)
		}
	}

	return fileList, nil
}

// ListFilesInDir returns list of all files in directory
func ListFilesInDir(owner, repo, token, dir string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, dir)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get contents: %s", resp.Status)
	}

	var files []GitHubFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	var fileList []string
	for _, file := range files {
		if file.Type == "file" {
			fileList = append(fileList, file.Path)
		} else if file.Type == "dir" {
			subFiles, err := ListFilesInDir(owner, repo, token, file.Path)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, subFiles...)
		}
	}

	return fileList, nil
}
