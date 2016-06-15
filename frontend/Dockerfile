FROM golang:1.6-wheezy

RUN go get github.com/mholt/caddy/caddy

ADD . /site
WORKDIR /site

EXPOSE 2015
CMD ["caddy"]
