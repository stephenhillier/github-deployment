# github-deployment
create GitHub deployments in a container-based script or CI/CD pipeline

## Brigade usage

To use the container in a Brigade pipeline, set the following environment variables:

```javascript
const env = {
    EVENT_PAYLOAD: e.payload,
    ENVIRONMENT_NAME: "staging",
    ENVIRONMENT_URL: "https://staging.example.com",
    DEPLOYMENT_STATE: "pending",
    GITHUB_TOKEN: "a github token that has rights to the repository",
    DEPLOYMENT_ID: 1234 // omit if not known, one will be generated and returned
}
```
