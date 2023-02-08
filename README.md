# keycloak-slackbot

Slackbot that reports new keycloak user registrations

## Setup

### Keycloak setup

We are going to be using service accounts to assume roles rather than using a specific use or admin credentials

- Add a new client to keycloak. Make it an oidc protocol with confidential access-type.
- Allow service accounts
- In the service account tab of the client, add the "view users" role from the client role "realm-management"

### Slack Setup

- Create a new slack app to post to a channel and save the url as `$SLACK_URL` env var

### Dev setup

You need 5 more (in addition to `$SLACK_URL`) env vars of information to debug.

1. `$KEYCLOAK_URL`: The base host that keycloak is located at www.<>.com, no paths needed
2. `$KEYCLOAK_USER`: This is the client_id from keycloak that we setup above
3. `$KEYCLOAK_PASSWORD`: This is the client_secret
4. `$KEYCLOAK_REAM`: Is the name of the keycloak realm you are using.
5. `$INTERVAL`: In seconds, how often to run this check.

## Helm chart

Helm chart is stored in the skyfjell helm chart repo https://github.com/skyfjell/charts
