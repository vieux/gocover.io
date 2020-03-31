VERSIONS="1.13.9 1.14 1.14.1"
LATEST="1.14.x"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
    docker push vieux/gocover:$version
done

docker tag vieux/gocover:1.13.9 vieux/gocover:1.13.x
docker push vieux/gocover:1.13.x
docker tag vieux/gocover:1.14.1 vieux/gocover:1.14.x
docker push vieux/gocover:1.14.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
docker push vieux/gocover:latest

