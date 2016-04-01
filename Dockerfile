FROM ubuntu:14.04

# required for rsvg-convert dependency
RUN apt-get -y update && apt-get install -y librsvg2-bin && apt-get install -y curl

ENV GOLANG_VERSION 1.6
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 5470eac05d273c74ff8bac7bef5bad0b5abbd1c4052efbdbc8db45332e836b0b

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

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