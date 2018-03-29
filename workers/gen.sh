VERSIONS="1.8 1.8.1 1.8.2 1.8.3 1.8.4 1.8.5 1.9 1.9.1 1.9.2 1.9.3 1.9.4 1.9.5 1.10 1.10.1"
LATEST="1.10.x"

for version in $VERSIONS; do
    echo "docker build --build-arg GO_VERSION=$version -t vieux/gocover:$version ."
done

echo "docker tag vieux/gocover:1.8.5 vieux/gocover:1.8.x"
echo "docker tag vieux/gocover:1.9.5 vieux/gocover:1.9.x"
echo "docker tag vieux/gocover:1.10.1 vieux/gocover:1.10.x"
echo "docker tag vieux/gocover:$LATEST vieux/gocover:latest"
