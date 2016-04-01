# Go installed and a workspace (GOPATH) configured at /go.
FROM golang:1.6.0-wheezy

# required for rsvg-convert dependency
RUN apt-get -y update && apt-get install -y librsvg2-bin

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/loganjspears/slackchess

# Build the slackchess command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/loganjspears/slackchess/cmd/slackchess

# Run the command by default when the container starts.
ENTRYPOINT /go/bin/slackchess -token=$TOKEN -url=$URL

# Document that the service listens on port 5000.
EXPOSE 5000