# Ubuntu 14.04
FROM ubuntu:14.04
ENV BUMP_THIS_TO_INVALIDATE_DOCKER_CACHE 1
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update && apt-get upgrade -y && apt-get install -y tree cowsay

# Go 1.3.1
RUN apt-get install -y curl git mercurial build-essential
RUN curl --silent --location https://golang.org/dl/go1.3.1.linux-amd64.tar.gz > /tmp/go.tar.gz
RUN tar --directory=/usr/local/ -xzf /tmp/go.tar.gz
ENV GOPATH /gopath
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin
RUN mkdir -p "$GOPATH"
RUN go get -v github.com/tools/godep
RUN go get -v github.com/golang/lint/golint
RUN go get -v github.com/kisielk/errcheck

# Get Gondrian dependencies
RUN apt-get install -y postgresql
