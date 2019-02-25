# github-deployment
create GitHub deployments in a container-based script or CI/CD pipeline

## Basic usage

github-deployment helps create GitHub deployment notifications for pull requests that are deployed to testing environments. After running the application (either using CLI or by using the
`stephenhillier/github-deployment` container), a notification will appear in your pull request's page with the status of the deployment and a link to the environment.

For more information on the GitHub Deployment API, see https://developer.github.com/v3/repos/deployments/

Options can be passed as arguments to the cli command (single dash) or as UPPERCASE environment variables. Arguments will override environment variables.

* `event_payload`: if using github-deployment with a a CI/CD system triggered by webhooks, this is the webhook payload that arrives from GitHub after a pull request is opened. Currently only supports `pull_request` events.
* `environment_name`: the name of the environment (e.g. "staging", "pr-123", etc.)
* `environment_url`: URL where the preview environment can be accessed
* `deployment_state`: one of `pending`, `success`, `failure`, `in_progress`, `queued`, `inactive`
* `github_token`: a GitHub OAuth token with permissions to create deployments and statuses in the target repo. The account associated with this token will also appear alongside notifications; e.g. "my_github_user deployed to staging 12 minutes ago"
* `deployment_id`:  when updating a deployment status, this is the ID of the deployment to update. If you don't have a deployment yet, omit this value and one will be automatically created.  The application's output will be a JSON string containing the deployment ID.  When you want to update the deployment (to "success" or "inactive"), store the response and call `github-deployment` again with the new deployment_id.

Example: `$ DEPLOYMENT_ID=123 github-deployment -event_payload $PAYLOAD -environment_name "staging" -deployment_state "success" -environment_url "https://www.example.com"`

## Brigade usage

To use the `stephenhillier/github-deployment:v0.1.0` container in a Brigade pipeline, set the following environment variables and run the container as a job:

```javascript
const env = {
    EVENT_PAYLOAD: e.payload,
    ENVIRONMENT_NAME: "staging",
    ENVIRONMENT_URL: "https://staging.example.com",
    DEPLOYMENT_STATE: "pending",
    GITHUB_TOKEN: "a github token that has rights to the repository",
    DEPLOYMENT_ID: 1234 // omit if not known, one will be generated and returned
}

const notify = new Job("deploy-notify", "stephenhillier/github-deployment:v0.1.0")
notify.imageForcePull = true
notify.env = env
notify.run()
```

You can omit `DEPLOYMENT_ID` on the first run, and a new deployment will be generated and returned. To update that deployment (e.g. `pending` -> `success`), run a second job but pass in the deployment id.

## GitHub Actions usage

Add a step to your workflow like this:

```hcl
action "GitHub Deployment" {
    uses = "stephenhillier/github-deployment@master"
    args = [
        "-environment_name", "staging",
        "-deployment_state", "success" 
    ]
}
```

Other parameters can be passed by setting environment variables for the stage.
Note: GitHub Actions support is experimental. The `event_payload` parameter may be a blocker. Feel free to create an issue/PR to report your experience.
