#!/usr/bin/env bash
set -e

echo "Checking for docker-compose..."
hash docker-compose > /dev/null 2>&1 \
    && docker-compose version | grep "docker-compose version" \
    || echo "Please install docker-compose: https://docs.docker.com/compose/install/"

echo "Checking for ab (ApacheBench)..."
hash ab > /dev/null 2>&1 \
    || echo "Please install ApacheBench;\nUbuntu: apt-get install apache2-utils"

echo "Pulling some common bases..."
docker pull smebberson/alpine-consul
docker pull smebberson/alpine-consul-redis
docker pull smebberson/alpine-consul-ui

docker pull gliderlabs/registrator:latest

echo "Pre-build some slow images..."
docker build demo-01-basics/webapp
docker build demo-02-registry/registrator
