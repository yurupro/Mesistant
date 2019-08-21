FROM golang:1.12-alpine

RUN apk add git
RUN apk add mongodb

VOLUME /data/db

WORKDIR /go/src
RUN git clone https://github.com/yurupro/Mesistant.git
EXPOSE 8080

CMD ["/usr/local/go/bin/go", "run", "/go/src/Messistant/server.go"]
