cd src
for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        echo "Building for \033[1;33m$GOOS-$GOARCH\033[0m"
        start=`date +%s`
        rm -rf build/needle-url-$GOOS-$GOARCH
        CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o ../build/needle-url-$GOOS-$GOARCH
        end=`date +%s`
        runtime=$((end-start))
        echo "Build for \033[1;33m$GOOS-$GOARCH\033[0m in $runtime seconds"
    done
done