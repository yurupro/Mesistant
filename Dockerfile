FROM golang:1.12-alpine3.9
LABEL maintainer="Prokuma <prokuma@prokuma.kr>"

RUN apk add git

ADD . /go/src/Mesistant

WORKDIR /go/src/Mesistant

ENV GO111MODULE=on
RUN go mod download

RUN go build .

EXPOSE 8080

CMD ["./Mesistant"]
