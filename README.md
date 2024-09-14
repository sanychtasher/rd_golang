# rd_golang

#### Run Client: 
`go mod tidy && go mod vendor`

`go run -mod=vendor ./cmd/client --name yermakov --host :8456`


#### Run Server
`go mod tidy && go mod vendor`

`go run -mod=vendor ./cmd/server --host :8456 .`

#### Docker Build 
`docker build --tag server .`

#### Docker Run
`docker run -p 8456:8456 server`