FROM golang:1.12-alpine3.9
LABEL maintainer="Prokuma <prokuma@prokuma.kr>"

RUN apk add git

WORKDIR /go/src
RUN git clone https://github.com/yurupro/Mesistant.git

WORKDIR /go/src/Mesistant

ENV GO111MODULE=on
RUN go mod download

RUN go build .

EXPOSE 8080

CMD ["./Mesistant"]
