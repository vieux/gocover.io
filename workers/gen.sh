VERSIONS="1.18 1.17.8"
LATEST="1.18.x"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
    docker push vieux/gocover:$version
done

docker tag vieux/gocover:1.17.8 vieux/gocover:1.17.x
docker push vieux/gocover:1.17.x
docker tag vieux/gocover:1.18 vieux/gocover:1.18.x
docker push vieux/gocover:1.18.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
docker push vieux/gocover:latest
