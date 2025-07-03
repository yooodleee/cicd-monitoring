package main

import (
	githubaction "cicd-monitoring/github_action"
	router "cicd-monitoring/router"
	"fmt"
	"os"
	"github.com/gofiber/fiber/v2"
)


func main() {
	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	client := githubaction.NewGitHubClient(token, owner, repo)

	app := fiber.New()
	router.SetupRoutes(app, client)

	app.Listen(":8080")

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