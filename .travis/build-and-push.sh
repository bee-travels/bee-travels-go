#!/bin/bash
prefix=go
echo "prefix set to $prefix"

docker build -t beetravels/destination-v1:${prefix}-$TRAVIS_COMMIT -f services/destination-v1/docker/Dockerfile destination-v1
docker build -t beetravels/destination-v2:${prefix}-$TRAVIS_COMMIT -f services/destination-v2/docker/Dockerfile destination-v2

echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push beetravels/destination-v1:${prefix}-$TRAVIS_COMMIT
docker push beetravels/destination-v2:${prefix}-$TRAVIS_COMMIT
