package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"parser/services"
)

func main() {
	// Parsing command line arguments
	repoOwner := flag.String("owner", "", "Owner of the GitHub repository")
	repoName := flag.String("repo", "", "Name of the GitHub repository")
	flag.Parse()

	if *repoOwner == "" || *repoName == "" {
		fmt.Println("You must specify the repository owner, name, and GitHub token using the -owner, -repo, and -token flags.")
		return
	}

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found.")
		return
	}

	gitAccessToken, exists := os.LookupEnv("GIT_ACCESS_TOKEN")
	if !exists {
		fmt.Println("No GITHUB access token provided.")
		return
	}

	// Getting all repo files' names
	files, err := services.ListRepFiles(*repoOwner, *repoName, gitAccessToken)
	if err != nil {
		fmt.Println("Error listing files:", err)
		return
	}

	fmt.Println("Files in repository:")
	for _, file := range files {
		fmt.Println(file)
	}
}
