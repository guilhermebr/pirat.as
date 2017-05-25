FROM golang:1.8-wheezy

ADD . $GOPATH/src/github.com/guilhermebr/pirat.as/backend
WORKDIR $GOPATH/src/github.com/guilhermebr/pirat.as/backend
VOLUME /data

CMD ["go", "run", "main.go"]
