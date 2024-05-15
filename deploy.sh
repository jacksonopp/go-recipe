#!/usr/bin/env bash

git pull

BUILD_VERSION=$(git rev-parse HEAD)

:echo "$(date --utc +%FT%TZ) Releaseing new server version. $BUILD_VERSION"

echo "$(date --utc +%FT%TZ): Running build..."
docker compose rm -f
docker compose build

OLD_CONTAINER=$(docker ps -aqf "name=server")
echo "$(date --utc +%FT%TZ): Scaling server up..."
BUILD_VERSION=$BUILD_VERSION docker compose up -d --no-deps --scale server=2 --no-recreate server

sleep 30

echo "$(date --utc +%FT%TZ): Scaling old server down..."
docker container rm -f $OLD_CONTAINER
docker compose up -d --no-deps --scale server=1 --no-recreate server