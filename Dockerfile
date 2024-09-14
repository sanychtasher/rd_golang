FROM golang:1.23.1-alpine3.20

WORKDIR server
COPY /cmd/server .

COPY go.mod .
COPY go.sum .
COPY vendor .

RUN go mod tidy && \
    go mod vendor && \
    go build -mod=vendor -o server .

EXPOSE 8456

ENTRYPOINT ["./server", "--host", ":8456"]
