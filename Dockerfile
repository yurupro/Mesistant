FROM golang:1.12-alpine3.9

RUN apk add git

WORKDIR /go/src
RUN git clone https://github.com/yurupro/Mesistant.git

WORKDIR /go/src/Mesistant
RUN git checkout dev

ENV GO111MODULE=on
RUN go mod download

RUN go build .

EXPOSE 8080

CMD ["./Mesistant"]
