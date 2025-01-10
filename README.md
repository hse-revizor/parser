The service for parsing data from GitHub repositories written in Go.
# Usage

I. Change "YOUR_TOKEN_HERE" to your GitHub access token in .env file.
```
GIT_ACCESS_TOKEN=YOUR_TOKEN_HERE
```

II. Run the command below passing repository name and owner with -repo and -owner flags
```
go run main.go -owner "owner's name" -repo "repo's -name"
```