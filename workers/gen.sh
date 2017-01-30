VERSIONS="1.6.4 1.7.5"
LATEST="1.7.5"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version . 
done
    
docker tag vieux/gocover:$LATEST vieux/gocover:latest
