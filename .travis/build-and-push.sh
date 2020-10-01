#!/bin/bash
prefix=go
echo "prefix set to $prefix"

docker build -t beetravels/destination-v1:${prefix}-$TRAVIS_COMMIT -f destination-v1/docker/Dockerfile destination-v1

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push beetravels/destination-v1:${prefix}-$TRAVIS_COMMIT