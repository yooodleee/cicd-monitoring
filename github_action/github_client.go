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


type StepDetail struct {
	Name        string
	Status      string
	Conclusion  string
	Number      int
	StartedAt   time.Time
	CompletedAt time.Time
}

// Job 목록 조회
type JobDetail struct {
	Name        string
	Status      string
	Conclusion  string
	StartedAt   time.Time
	CompletedAt time.Time
	Steps 		[]StepDetail
}

func (gh *GitHubClient) GetJobsForRun(runID int64) ([]JobDetail, error) {
	opts := &github.ListWorkflowJobsOptions{
		ListOptions: github.ListOptions{PerPage: 20},
	}
	jobs, _, err := gh.client.Actions.ListWorkflowJobs(gh.ctx, gh.owner, gh.repo, runID, opts)
	if err != nil {
		return nil, err
	}

	var details []JobDetail
	for _, job := range jobs.Jobs {
		var steps []StepDetail
		for _, s := range job.Steps {
			steps = append(steps, StepDetail{
				Name: s.GetName(),
				Status: s.GetStatus(),
				Conclusion: s.GetConclusion(),
				Number: int(s.GetNumber()),
				StartedAt: s.GetStartedAt().Time,
				CompletedAt: s.GetCompletedAt().Time,
			})
		}
		
		d := JobDetail{
			Name:       job.GetName(),
			Status:     job.GetStatus(),
			Conclusion: job.GetConclusion(),
			StartedAt:  job.GetStartedAt().Time,
			CompletedAt: job.GetCompletedAt().Time,
			Steps: steps,
		}
		details = append(details, d)
	}
	return details, nil
}