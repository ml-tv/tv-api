# tv-api

## Master badges
[![Build Status](https://travis-ci.org/ml-tv/tv-api.svg?branch=master)](https://travis-ci.org/ml-tv/tv-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/ml-tv/tv-api)](https://goreportcard.com/report/github.com/ml-tv/tv-api)
[![codebeat badge](https://codebeat.co/badges/927f93b6-f37d-4bbe-aa7a-666af5bfb695)](https://codebeat.co/projects/github-com-ml-tv-tv-api)

## Staging badges
[![Build Status](https://travis-ci.org/ml-tv/tv-api.svg?branch=staging)](https://travis-ci.org/ml-tv/tv-api)

## Run the API using docker

```
docker-compose build
docker-compose up -d
```

Bash helpers can be found in `tools/docker-helpers.sh`

## travis

```
travis encrypt HEROKU_API_KEY=$(heroku auth:token) --add
travis encrypt TMDB_API_KEY=xxxxxxxxxxxxxxxxxx --add
```