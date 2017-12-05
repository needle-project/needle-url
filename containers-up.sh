echo "Compiling..."
cd src
CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o ../docker/http/server

echo "Building compose"
docker-compose build

echo "Launching up containers"
docker-compose up -d