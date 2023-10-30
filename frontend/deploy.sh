#!/bin/bash

set -o errexit

readonly NAME="expenses_web"
readonly RESOURCES="/data/local/data/resources"
readonly PORT=5000

docker build . -t $NAME
docker stop $NAME
docker rm $NAME 
docker run --name $NAME  --net nettyw -d -v $RESOURCES:/app/resources $NAME
