FROM golang:1.6-wheezy

RUN go get github.com/tools/godep
ADD . $GOPATH/src/github.com/guilhermebr/pirat.as/backend

WORKDIR $GOPATH/src/github.com/guilhermebr/pirat.as/backend

RUN godep restore

CMD ["go", "run", "main.go"]
