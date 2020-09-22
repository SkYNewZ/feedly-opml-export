#!/usr/bin/env bash

set -x

echo "Building skynewz/feedly-opml-export:latest"
docker buildx build \
--tag skynewz/feedly-opml-export:latest \
--push \
--platform darwin/amd64,linux/amd64,linux/arm/v7,linux/arm/v6,linux/arm64 \
--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
--build-arg SHA=${SHA} \
.

echo "Building skynewz/feedly-opml-export:${SHA}"
docker buildx build \
--tag skynewz/feedly-opml-export:${SHA} \
--push \
--platform darwin/amd64,linux/amd64,linux/arm/v7,linux/arm/v6,linux/arm64 \
--build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
--build-arg SHA=${SHA} \
.