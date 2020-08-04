# 🦉 Jia

Keeping an eye on the Hack Club Slack’s [#counttoamillion](https://hackclub.slack.com/archives/CDJMS683D) channel.

## Setup

### Creating the Slack app

You’ll need to create a Slack app (not a classic one) with at least the following bot token scopes. The reasons each scope is required are also included:

- `channels:history` (or `groups:history`, if it’s a private channel): Used to listen to messages sent.
- `chat:write`: Used for sending messages.
- `reactions:write`: For reacting to invalid messages.

Then you’ll need to subscribe the app to a few events. The server has an endpoint open at `/slack/events`, so when you’re asked for a request URL, just put `https://<SERVER>/slack/events`. Only the following events are needed:

- `message.channels` (or `message.groups` if it’s a private channel)

### Environment variables

Here are all the variables you need to set up, with hints.

```bash
# The port to run the app server on
PORT=3000
# Redis database to store the last number and its sender
REDIS_URL=redis://…
# App config. Obtained from the "Basic Information" page of your app.
SLACK_BOT_TOKEN=xoxb-…
SLACK_VERIFICATION_TOKEN=xxxx…
# The channel where Jia should validate counted numbers in.
SLACK_CHANNEL_ID=C…
```

### Deploying

```bash
# Run it…
make

# …or build a binary and run that instead
make build
./bin/jia
```

After you’ve followed all the above steps, you should see something like this:

```bash
Starting Jia…
Listening on port 3000
```

## License

[MIT License](LICENSE.txt)
