#!/bin/bash

sudo docker build . -t expenses
sudo docker stop expenses
sudo docker rm expenses
sudo docker run --name expenses -d -v /data/local/data/resources:/app/resources -p 5000:80 expenses
