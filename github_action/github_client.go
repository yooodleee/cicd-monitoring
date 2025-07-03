package githubaction

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"github.com/google/go-github/v55/github"
)


// WorkflowRunSummary 구조체에 필요한 요약 정보를 담습니다.
type WorkflowRunSummary struct {
	ID 			int64
	Status 		string
	Conclusion 	string
	CommitMsg 	string
	Branch 		string
	Duration 	time.Duration
	TriggeredBy string
	CreatedAt 	time.Time
}

// GitHubClient 구조체 정의
type GitHubClient struct {
	client  *github.Client
	ctx 	context.Context
	owner 	string
	repo 	string
}

// GitHubClient 생성자
func NewGitHubClient(token, owner, repo string) *GitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &GitHubClient{
		client: github.NewClient(tc),
		ctx: 	ctx,
		owner: 	owner,
		repo: 	repo,
	}
}

// workflow 실행 이력 조회
func (gh *GitHubClient) ListWorkflowRuns() ([]WorkflowRunSummary, error) {
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	runs, _, err := gh.client.Actions.ListRepositoryWorkflowRuns(gh.ctx, gh.owner, gh.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflow runs: %w", err)
	}

	var results []WorkflowRunSummary
	for _, run := range runs.WorkflowRuns {
		duration := time.Duration(0)
		if run.RunStartedAt != nil && run.UpdatedAt != nil {
			duration = run.UpdatedAt.Sub(run.RunStartedAt.Time)
		}

		summary := WorkflowRunSummary {
			ID: run.GetID(),
			Status: run.GetStatus(),
			Conclusion: run.GetConclusion(),
			CommitMsg: run.GetHeadCommit().GetMessage(),
			Branch: run.GetHeadBranch(),
			TriggeredBy: run.GetActor().GetLogin(),
			CreatedAt: run.GetCreatedAt().Time,
			Duration: duration,
		}
		results = append(results, summary)
	}

	return results, nil
}


// 수동 trigger 기능
func (gh *GitHubClient) TriggerWorkflow(workflowFileName, branch string, inputs map[string]interface{}) error {
	dispatch := github.CreateWorkflowDispatchEventRequest{
		Ref: branch,
		Inputs: inputs,
	}

	_, err := gh.client.Actions.CreateWorkflowDispatchEventByFileName(
		gh.ctx,
		gh.owner,
		gh.repo,
		workflowFileName,
		dispatch,
	)
	if err != nil {
		return fmt.Errorf("failed to dispatch workflow: %w", err)
	}
	return nil
}