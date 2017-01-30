VERSIONS="1.6.4 1.7.5 latest"
for version in $VERSIONS; do
    docker push vieux/gocover:$version
done
