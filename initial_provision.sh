echo "Setting GOPATH to ${PWD}"
export GOPATH=$PWD
echo "Getting dependencies"
go get -u github.com/go-redis/redis
go get -u github.com/gorilla/mux
echo "Done"
