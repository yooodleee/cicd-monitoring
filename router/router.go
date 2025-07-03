package router

import (
	"html/template"
	"github.com/gofiber/fiber/v2"

	githubaction "cicd-monitoring/github_action"
)


func SetupRoutes(app *fiber.App, client *githubaction.GitHubClient) {
	app.Get("/router", func(c *fiber.Ctx) error {
		runs, err := client.ListWorkflowRuns()
		if err != nil {
			return c.Status(500).SendString("Failed to load workflows")
		}

		// 간단한 템플릿 엔진 사용 (html/template)
		tmpl, err := template.ParseFiles("cicd-monitoring/views/dashboard.html")
		if err != nil {
			return c.Status(500).SendString("Failed to parse template")
		}

		// Fiber의 c.SendWriter() 사용
		c.Type("html")
		return tmpl.Execute(c.Response().BodyWriter(), runs)
	})
}