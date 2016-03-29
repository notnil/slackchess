# Start from a alpine image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.6.0-wheezy

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/loganjspears/slackchess

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/loganjspears/slackchess

# required for rsvg-convert dependency
RUN apt-get -y update && apt-get install -y librsvg2-bin

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/slackchess -token=$TOKEN -url=$URL

# Document that the service listens on port 8080.
EXPOSE 5000