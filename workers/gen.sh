VERSIONS="1.8"
LATEST="1.8"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version . 
done
    
docker tag vieux/gocover:$LATEST vieux/gocover:latest
