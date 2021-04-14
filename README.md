# Reddit Scan Bot

## About this repo
A script that sends a notification when certain keywords are matched in a new post for a given subreddit.

To run this locally, clone it and run `go run main.go`. Be sure to set the following env vars!
```
export TELEGRAM_CHAT_ID=""
export TELEGRAM_TOKEN=""
export REDDIT_SEARCH_SUBREDDIT="buildapcsales"
export REDDIT_SEARCH_TERMS="Monitor,Xbox"
```

#### \* This repo is a WIP and provided 'as-is' with no guarantee for support. You'll notice several commented out blocks and some TODOs scattered throughout the code. I didn't intend to post this code until it was complete, but here we are! ¯\\\_(ツ)\_/¯  