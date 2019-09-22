#!/bin/bash

docker build --build-arg name=1 -t r.kfd.me/piu:hello-world-1 .
docker build --build-arg name=2 -t r.kfd.me/piu:hello-world-2 .
