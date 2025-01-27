package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GitHubFileContent represents the structure of the file content response from the GitHub API
type GitHubFileContent struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// ConvertGitHubURLToAPIURL converts a GitHub URL to a GitHub API URL
func ConvertGitHubURLToAPIURL(gitHubURL string) (string, error) {
	parsedURL, err := url.Parse(gitHubURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Check if the URL belongs to GitHub
	if parsedURL.Host != "github.com" {
		return "", fmt.Errorf("URL must belong to github.com")
	}

	pathParts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(pathParts) < 5 || pathParts[2] != "blob" {
		return "", fmt.Errorf("URL must point to a file in a repository in the format: https://github.com/{owner}/{repo}/blob/{branch}/{filePath}")
	}

	owner := pathParts[0]
	repo := pathParts[1]
	branch := pathParts[3]
	filePath := strings.Join(pathParts[4:], "/")

	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, filePath, branch)
	return apiURL, nil
}

// FetchFileContentFromURL fetches the content of a file using a GitHub API URL
func FetchFileContentFromURL(apiURL, token string) (string, error) {
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file content: %s", resp.Status)
	}

	var fileContent GitHubFileContent
	if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if fileContent.Encoding != "base64" {
		return "", fmt.Errorf("unsupported file encoding: %s", fileContent.Encoding)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode file content: %w", err)
	}

	return string(decodedContent), nil
}
