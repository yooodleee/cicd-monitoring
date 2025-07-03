package router

import (
	"fmt"
	"html/template"
	"log"

	"github.com/gofiber/fiber/v2"

	githubaction "cicd-monitoring/github_action"
)


func SetupRoutes(app *fiber.App, client *githubaction.GitHubClient) {
	app.Get("/dashboard", func(c *fiber.Ctx) error {
		runs, err := client.ListWorkflowRuns()
		if err != nil {
			return c.Status(500).SendString("GitHub API 에러")
		}

		// 간단한 템플릿 엔진 사용 (html/template)
		tmpl, err := template.ParseFiles("views/dashboard.html")
		if err != nil {
			log.Println("Template error:", err)
			return c.Status(500).SendString("Failed to parse template")
		}

		// Fiber의 c.SendWriter() 사용
		c.Type("html")
		return tmpl.Execute(c.Response().BodyWriter(), runs)
	})

	// /trigger 라우트 추가
	app.Post("/trigger", func(c *fiber.Ctx) error {
		type TriggerRequest struct {
			WorkflowFile string `json:"workflow_file"`
			Branch		 string `json:"branch"`
		}

		var req TriggerRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).SendString("Invalid request")
		}

		err := client.TriggerWorkflow(req.WorkflowFile, req.Branch, nil)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed: %v", err))
		}

		return c.SendString("✅ Workflow triggered!")
	})
}	