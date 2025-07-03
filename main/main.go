package main

import (
	githubaction "cicd-monitoring/github_action"
	"fmt"
	"os"
)


func main() {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	client := githubaction.NewGitHubClient(token, owner, repo)

	runs, err := client.ListWorkflowRuns()
	if err != nil {
		panic(err)
	}

	for _, run := range runs {
		fmt.Printf("✅ [%s] %s (%s) by %s | %s\n",
			run.Conclusion,
			run.CommtiMsg,
			run.Branch,
			run.TriggeredBy,
			run.Duration,
		)
	}
}