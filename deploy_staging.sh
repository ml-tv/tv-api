#!/usr/bin/env bash

set -ex

heroku plugins:install heroku-container-registry
heroku container:login
heroku container:push web --app ml-tv
heroku run "cd /go/src/github.com/ml-tv/tv-api && goose up" --app ml-tv