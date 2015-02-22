FROM golang:1.4.1

ENV GOPATH /go

ADD . /go/src/github.com/Elemental-IRCd/elemental-ircd/test

WORKDIR /go/src/github.com/Elemental-IRCd/elemental-ircd/test

RUN go get github.com/thoj/go-ircevent

CMD sh ./run.sh
