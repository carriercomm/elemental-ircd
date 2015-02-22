FROM golang:1.4.1

ENV GOPATH /go

RUN go get github.com/thoj/go-ircevent

ADD . /go/src/github.com/Elemental-IRCd/elemental-ircd/test

WORKDIR /go/src/github.com/Elemental-IRCd/elemental-ircd/test

CMD sh ./run.sh
