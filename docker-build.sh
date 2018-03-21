#!/bin/bash

version=$1

docker build --no-cache --build-arg "VERSION=$version" --tag "korylprince/bisd-inventory-identifier-server:$version" .

docker push "korylprince/bisd-inventory-identifier-server:$version"
