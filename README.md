# pr-bot
PR-Bot is an event driven service that listens to GitHub webhooks and performs actions based on the events received.

## Development

### Local Setup
In order to run the PR-Bot locally, you will need to setup a webhook proxy to receive the GitHub webhooks and forward them to your local service. You can use [smee.io](https://smee.io/) for this purpose.
1. https://smee.io/ to setup the webhook proxy
2. Setup the webhook in the github repository with the smee.io url with the pull request scopes. NOTE: Make sure that the `content-type` is set to `application/json` in the webhook settings.
3. Get your AWS credentials
4. Run the pr-bot service locally with `make run`
5. In a separate terminal run `smee -u SMEE_CHANNEL_URL --target http://localhost:9090/v1/webhook`
