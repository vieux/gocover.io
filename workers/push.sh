VERSIONS="1.8 latest"
for version in $VERSIONS; do
    docker push vieux/gocover:$version
done
