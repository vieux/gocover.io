VERSIONS="1.8 1.8.1 1.8.2 1.8.3 1.8.4 1.8.5 1.9 1.9.1 1.9.2 1.9.3 1.9.4 1.10 latest"

for version in $VERSIONS; do
    docker push vieux/gocover:$version
done
