VERSIONS="1.8 1.8.1 1.8.2 1.8.3 1.8.4 1.8.5 1.9 1.9.1 1.9.2 1.9.3 1.9.4 1.10"
LATEST="1.10"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
done

docker tag vieux/gocover:1.8.5 vieux/gocover:1.8.x
docker tag vieux/gocover:1.9.4 vieux/gocover:1.9.x
docker tag vieux/gocover:1.10 vieux/gocover:1.10.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
