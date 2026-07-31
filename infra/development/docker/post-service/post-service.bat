set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

go build -o build/post-service ./services/post-service