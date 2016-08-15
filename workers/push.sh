VERSIONS="1.2.2 1.3 1.3.1 1.3.2 1.3.3 1.4 1.4.1 1.4.2 1.4.3 1.5 1.5.1 1.5.2 1.5.3 1.6 1.6.1 1.6.2 1.6.3 1.7 latest"
for version in $VERSIONS; do
    docker push vieux/gocover:$version
done
