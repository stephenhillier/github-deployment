package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/namsral/flag"
	"golang.org/x/oauth2"
)

// PullRequestEvent is the payload sent by GitHub
// when a pull request event occurs (opened, closed, labeled, etc.)
type PullRequestEvent struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	PullRequest struct {
		Number int `json:"number"`
	} `json:"pull_request"`
}

// DeploymentResult contains the results of deployment requests
// (deployment and deployment status), and will be outputted
// when the program exits.
type DeploymentResult struct {
	Deployment int64  `json:"deployment"`
	State      string `json:"state"`
}

func main() {
	var payload string
	var envURL string
	var envName string
	var token string
	var status string
	var deploymentID int64

	// these flags can be provided as environment variables (uppercase) e.g. EVENT_PAYLOAD
	flag.StringVar(&payload, "event_payload", "", "the github pull_request event payload")
	flag.StringVar(&envURL, "environment_url", "", "the URL that the deployment can be accessed at")
	flag.StringVar(&envName, "environment_name", "", "the name of the deployment environment (e.g. staging)")
	flag.StringVar(&token, "github_token", "", "GitHub OAuth token")
	flag.StringVar(&status, "deployment_state", "pending", "the status of the deployment (success, pending, inactive, etc)")
	flag.Int64Var(&deploymentID, "deployment_id", 0, "the ID of the deployment (if omitted, one will be generated)")
	flag.Parse()

	// make sure a token is provided to access GitHub API
	// TODO: this isn't strictly necessary - check could potentially be removed
	if token == "" {
		log.Fatal("GITHUB_TOKEN not set")
	}
	oauthToken := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// set up the client to access GitHub, using the token
	httpClient := oauth2.NewClient(context.Background(), oauthToken)
	client := github.NewClient(httpClient)

	// set up the event using the payload that comes from Brigade pull_request event
	var event PullRequestEvent
	err := json.Unmarshal([]byte(payload), &event)
	if err != nil {
		log.Fatal("error parsing pull request event payload:", err)
	}

	// this program is configured to deploy a pull request ref e.g pull/1/head
	// in the future, this could be made more general (for example, to deploy a branch ref)
	if event.PullRequest.Number == 0 {
		log.Fatal("invalid event: event must be a pull request, but payload did not contain a pull request number")
	}

	result := DeploymentResult{Deployment: deploymentID}

	// if a deployment ID wasn't provided, create a new deployment
	if result.Deployment == 0 {
		deployment, err := createDeployment(client, event, envName)
		if err != nil {
			log.Fatal("error creating deployment:", err)
		}
		result.Deployment = deployment.GetID()
	}

	// set the deployment status
	deploymentStatus, err := createDeploymentStatus(client, result.Deployment, event, envName, status, envURL)
	if err != nil {
		log.Fatal("error creating deployment status:", err)
	}
	result.State = deploymentStatus.GetState()

	// successful deployment: print the deployment and state as JSON so it can be recorded
	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("error converting result to JSON")
	}
	fmt.Println(string(jsonResult))
}

// createDeploymentStatus takes a client and a GitHub Deployment and updates its status and environment URL.
// https://developer.github.com/v3/repos/deployments/
func createDeploymentStatus(client *github.Client, deployment int64, event PullRequestEvent, envName string, status string, URL string) (*github.DeploymentStatus, error) {
	repoName := strings.Split(event.Repository.FullName, "/")
	owner, repo := repoName[0], repoName[1]

	req := &github.DeploymentStatusRequest{
		State:          github.String(status),
		Environment:    github.String(envName),
		EnvironmentURL: github.String(URL),
	}
	ctx := context.Background()
	deploymentStatus, _, err := client.Repositories.CreateDeploymentStatus(
		ctx,
		owner,
		repo,
		deployment,
		req,
	)
	if err != nil {
		return deploymentStatus, err
	}
	return deploymentStatus, nil
}

// createDeployment sends a request to the GitHub Deployment API to create
// a new deployment for a pull request
// https://developer.github.com/v3/repos/deployments/
func createDeployment(client *github.Client, event PullRequestEvent, envName string) (*github.Deployment, error) {
	repoName := strings.Split(event.Repository.FullName, "/")
	owner, repo := repoName[0], repoName[1]
	ref := fmt.Sprintf("pull/%v/head", event.PullRequest.Number)

	req := &github.DeploymentRequest{
		Ref:                  github.String(ref),
		TransientEnvironment: github.Bool(true),
		Environment:          github.String(envName),
		RequiredContexts:     &[]string{},
	}
	ctx := context.Background()
	deployment, _, err := client.Repositories.CreateDeployment(ctx, owner, repo, req)
	if err != nil {
		return deployment, err
	}
	return deployment, nil
}
