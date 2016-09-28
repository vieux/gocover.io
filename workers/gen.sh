VERSIONS="1.2.2 1.3 1.3.1 1.3.2 1.3.3 1.4 1.4.1 1.4.2 1.4.3 1.5 1.5.1 1.5.2 1.5.3 1.6 1.6.1 1.6.2 1.6.3 1.7 1.7.1"
LATEST="1.7.1"

for version in $VERSIONS; do
    docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version . 
done
    
docker tag vieux/gocover:$LATEST vieux/gocover:latest
