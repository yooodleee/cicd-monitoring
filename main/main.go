package main

import (
	githubaction "cicd-monitoring/github_action"
	router "cicd-monitoring/router"
	"fmt"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)


func main() {
	// .env 파일 로드
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ .env 파일을 불러올 수 없습니다.")
	}

	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	fmt.Println("🔑 GITHUB_TOKEN set:", len(token) > 0)
	fmt.Println("🧑 GITHUB_OWNER:", owner)
	fmt.Println("🫙 GITHUB_REPO:", repo)

	client := githubaction.NewGitHubClient(token, owner, repo)

	// 진단용 workflow 목록 확인
	runs, err := client.ListWorkflowRuns()
	if err != nil {
		fmt.Println("❌ GitHub API 호출 실패:", err)
	} else {
		for _, run := range runs {
			fmt.Printf("✅ [%s] %s (%s) by %s | %s\n",
				run.Conclusion,
				run.CommitMsg,
				run.Branch,
				run.TriggeredBy,
				run.Duration,
			)
		}
	}

	app := fiber.New()
	router.SetupRoutes(app, client)

	fmt.Println("🚀 서버 실행 중 -> http://localhost:7070/dashboard")
	// 서버 실행은 마지막 라인에 -> blocking
	app.Listen(":7070")
}