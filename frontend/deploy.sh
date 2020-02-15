#!/bin/bash

readonly NAME="expenses"

docker build . -t $NAME
docker stop $NAME
docker rm $NAME expenses
docker run --name $NAME -d -v /data/local/data/resources:/app/resources -p 5000:80 $NAME
