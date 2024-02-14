#!/usr/bin/env bash

IMAGE_BASE="rinha-de-backend-2024-q1-edimarlnx:latest"
IMAGE_FINAL="edimarlnx/$IMAGE_BASE"

docker build -t $IMAGE_BASE -f docker/api.Dockerfile .
docker tag $IMAGE_BASE $IMAGE_FINAL
docker push $IMAGE_FINAL
#docker rmi $IMAGE_FINAL