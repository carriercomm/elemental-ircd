FROM yikaus/alpine-base

RUN apk update && apk add go
ENV GOPATH /go

ADD . /go/src/github.com/Elemental-IRCd/elemental-ircd/testsuite/docker/runner

WORKDIR /go/src/github.com/Elemental-IRCd/elemental-ircd/testsuite/docker/runner

CMD sh ./run.sh
