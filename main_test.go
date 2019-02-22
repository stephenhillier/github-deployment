package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-github/github"
)

func TestCreateDeployment(t *testing.T) {

	repoName := "testowner/testrepo"

	client := newTestGitHubClient()

	deployment, err := createDeployment(
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

	if *deployment.Ref != "pull/123/head" {
		t.Errorf("createDeployment did not format ref properly. Expected %s, received %s", "pull/123/head", *deployment.Ref)
	}

	if *deployment.Environment != "test" {
		t.Errorf("environment did not match. Expected %s, received %s", "test", *deployment.Environment)
	}

	if *deployment.RepositoryURL != fmt.Sprintf("https://www.github.com/%s", repoName) {
		t.Errorf("createDeployment did not format owner/repo properly. Expected %s, received %s", fmt.Sprintf("https://www.github.com/%s", repoName), *deployment.RepositoryURL)
	}

}

type testRepoClient struct{}

func TestCreateDeploymentStatus(t *testing.T) {
	repoName := "testowner/testrepo"

	client := newTestGitHubClient()

	deploymentStatus, err := createDeploymentStatus(
		client,
		123,
		PullRequestEvent{
			Repository:  GitHubRepository{FullName: repoName},
			PullRequest: GitHubPullRequest{Number: 123},
		},
		"test",
		"success",
		"https://www.example.com",
	)

	if err != nil {
		t.Error("error running createDeploymentStatus:", err)
	}

	// mainly checking that github.CreateDeploymentStatus was properly called with the values
	// passed into this package's createDeploymentStatus function
	if deploymentStatus.GetDeploymentURL() != "https://www.example.com" {
		t.Errorf("environment URL did not match. Expected %s, received %s", "https://www.example.com", deploymentStatus.GetDeploymentURL())
	}

	if deploymentStatus.GetState() != "success" {
		t.Errorf("state did not match. Expected %s, received %s", "success", deploymentStatus.GetState())
	}

}

func newTestGitHubClient() GitHubClient {
	return GitHubClient{Repositories: testRepoClient{}}
}

// CreateDeployment does some basic checking to make sure the method is called with the expected params, according to
// https://developer.github.com/v3/repos/deployments/
func (c testRepoClient) CreateDeployment(ctx context.Context, owner string, repo string, req *github.DeploymentRequest) (*github.Deployment, *github.Response, error) {

	// check for valid owner and repo names (basically, test that they were split properly)
	if a := strings.Split(owner, "/"); len(a) != 1 {
		return &github.Deployment{}, &github.Response{}, errors.New("owner name invalid (contains slash)")
	}
	if b := strings.Split(repo, "/"); len(b) != 1 {
		return &github.Deployment{}, &github.Response{}, errors.New("repo name invalid (contains slash)")
	}

	return &github.Deployment{
		Ref:           github.String(req.GetRef()),
		ID:            github.Int64(123),
		RepositoryURL: github.String(fmt.Sprintf("https://www.github.com/%s/%s", owner, repo)),
		Environment:   req.Environment,
	}, &github.Response{}, nil
}

func (c testRepoClient) CreateDeploymentStatus(ctx context.Context, owner string, repo string, deployment int64, req *github.DeploymentStatusRequest) (*github.DeploymentStatus, *github.Response, error) {
	// check for valid owner and repo names
	if a := strings.Split(owner, "/"); len(a) != 1 {
		return &github.DeploymentStatus{}, &github.Response{}, errors.New("owner name invalid (contains slash)")
	}
	if b := strings.Split(repo, "/"); len(b) != 1 {
		return &github.DeploymentStatus{}, &github.Response{}, errors.New("repo name invalid (contains slash)")
	}

	return &github.DeploymentStatus{
		DeploymentURL: req.EnvironmentURL,
		State:         req.State,
	}, &github.Response{}, nil
}
