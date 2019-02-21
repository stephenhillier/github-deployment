package main

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
)

func TestCreateDeployment(t *testing.T) {

	repoName := "testowner/testrepo"

	calledWithOwner := github.String("")
	calledWithRepo := github.String("")

	client := newTestGitHubClient()

	_, err := createDeployment(
		client,
		PullRequestEvent{
			Repository:  GitHubRepository{FullName: repoName},
			PullRequest: GitHubPullRequest{Number: 123},
		},
		"test",
	)

	if err != nil {
		t.Error("error running createDeployment:", err)
	}

	if *calledWithOwner != "testowner" {
		t.Errorf("createDeployment did not call github.CreateRepository with correct owner. Expected %s, received %s", "testowner", *calledWithOwner)
	}

	if *calledWithRepo != "testrepo" {
		t.Errorf("createDeployment did not call github.CreateRepository with correct repo. Expected %s, received %s", "testrepo", *calledWithRepo)
	}

}

type testRepoClient struct {
	CalledWithOwner string
	CalledWithRepo  string
}

func TestCreateDeploymentStatus(t *testing.T) {

}

func newTestGitHubClient() GitHubClient {
	return GitHubClient{Repositories: testRepoClient{}}
}

func (c testRepoClient) CreateDeployment(ctx context.Context, owner string, repo string, req *github.DeploymentRequest) (*github.Deployment, *github.Response, error) {
	return &github.Deployment{
		ID: github.Int64(123),
	}, &github.Response{}, nil
}

func (c testRepoClient) CreateDeploymentStatus(ctx context.Context, owner string, repo string, deployment int64, request *github.DeploymentStatusRequest) (*github.DeploymentStatus, *github.Response, error) {
	return &github.DeploymentStatus{}, &github.Response{}, nil
}
