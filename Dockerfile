FROM golang:1.7

# install depedencies
RUN go get bitbucket.org/liamstask/goose/cmd/goose

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/ml-tv/tv-api

# Install api binary globally within container
RUN cd /go/src/github.com/ml-tv/tv-api && make install

# Set binary as entrypoint
CMD /go/bin/tv-api

EXPOSE 5000