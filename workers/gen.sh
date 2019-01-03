VERSIONS="1.8 1.8.1 1.8.2 1.8.3 1.8.4 1.8.5 1.8.6 1.8.7 1.9 1.9.1 1.9.2 1.9.3 1.9.4 1.9.5 1.9.6 1.9.7 1.10 1.10.1 1.10.2 1.10.3 1.10.4 1.10.5 1.10.6 1.10.7 1.11 1.11.1 1.11.2 1.11.3 1.11.4 1.12beta1"
LATEST="1.11.x"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version .
    docker push vieux/gocover:$version
done

docker tag vieux/gocover:1.8.7 vieux/gocover:1.8.x
docker push vieux/gocover:1.8.x
docker tag vieux/gocover:1.9.7 vieux/gocover:1.9.x
docker push vieux/gocover:1.9.x
docker tag vieux/gocover:1.10.7 vieux/gocover:1.10.x
docker push vieux/gocover:1.10.x
docker tag vieux/gocover:1.11.4 vieux/gocover:1.11.x
docker push vieux/gocover:1.11.x
docker tag vieux/gocover:$LATEST vieux/gocover:latest
docker push vieux/gocover:latest

