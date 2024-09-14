FROM golang:1.23.1-alpine3.20

WORKDIR /go/src/github.com/sanychtasher/rd_golang
COPY /cmd/server .

RUN go mod init && \
    go mod tidy && \
    go mod vendor && \
    go build -mod=vendor -o server .

EXPOSE 8456

ENTRYPOINT ["./server", "--host", ":8456"]
