docker-compose down
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker/http/server src/main.go
docker-compose build
docker-compose up