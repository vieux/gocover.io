VERSIONS="1.12 .1.12.1 1.12.2 1.12.3 1.12.4 1.12.5 1.12.6 1.12.7 1.12.8 1.12.9 1.13"
LATEST="1.13.x"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
    docker push vieux/gocover:$version
done

docker tag vieux/gocover:1.12.9 vieux/gocover:1.12.x
docker push vieux/gocover:1.12.x
docker tag vieux/gocover:1.13 vieux/gocover:1.13.x
docker push vieux/gocover:1.13.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
docker push vieux/gocover:latest

