package router

import (
	"fmt"
	"html/template"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	githubaction "cicd-monitoring/github_action"
)


func SetupRoutes(app *fiber.App, client *githubaction.GitHubClient) {
	app.Get("/", func (c *fiber.Ctx) error {
		return c.Redirect("/dashboard")
	})

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
		WorkflowFile := c.FormValue("workflow_file")
		Branch := c.FormValue("branch")

		if WorkflowFile == "" || Branch == "" {
			return c.Status(400).SendString("⚠️ workflow_file 또는 branch가 비어 있습니다.")
		}

		err := client.TriggerWorkflow(WorkflowFile, Branch, nil)
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed: %v", err))
		}

		return c.SendString("✅ Workflow triggered!")
	})

	// Job 목록 조회 router 추가
	app.Get("/run/:id", func(c *fiber.Ctx) error {
		runIDStr := c.Params("id")
		runID, err := strconv.ParseInt(runIDStr, 10, 64)
		if err != nil {
			return c.Status(400).SendString("Invalid run ID")
		}

		jobs, err := client.GetJobsForRun(runID)
		if err != nil {
			return c.Status(500).SendString("Failed to fetch job details")
		}

		tmpl, err := template.ParseFiles("views/job_detail.html")
		if err != nil {
			log.Println("Template error:", err)
			return c.Status(500).SendString("Template error")
		}

		c.Type("html")
		return tmpl.Execute(c.Response().BodyWriter(), jobs)
	})
}	