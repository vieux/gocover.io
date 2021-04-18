VERSIONS="1.16.3 1.15.11"
LATEST="1.16.x"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
    docker push vieux/gocover:$version
done

docker tag vieux/gocover:1.15.11 vieux/gocover:1.15.x
docker push vieux/gocover:1.15.x
docker tag vieux/gocover:1.16.3 vieux/gocover:1.16.x
docker push vieux/gocover:1.16.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
docker push vieux/gocover:latest

