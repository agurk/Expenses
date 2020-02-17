#!/bin/bash

set -o errexit

readonly NAME="expenses"
readonly RESOURCES="/data/local/data/resources"
readonly PORT=5000

docker build . -t $NAME
docker stop $NAME
docker rm $NAME 
docker run --name $NAME -d -v $RESOURCES:/app/resources -p $PORT:80 $NAME
