FROM golang:latest

ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]
