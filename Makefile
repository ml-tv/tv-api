# Build info
VERSION=1.0.0
BUILD_INFO=`git rev-parse HEAD`

# Flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD_INFO)"

install:
	go install $(LDFLAGS) .

generate:
	go install $(LDFLAGS) github.com/ml-tv/tv-api/cmd/tv-api-cli

migration:
	goose up

.PHONY:
	install