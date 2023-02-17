#!/bin/bash

# Exit on Error.
set -e

# Image ID.
IMAGE_ID="$1"
if [ -z "$IMAGE_ID" ]
then
  echo "Image ID is not set."
  exit 1
fi

# Mark the Image and upload it.
REGISTRY="localhost:5000"
IMAGE_NAME="libre_office"
docker tag $IMAGE_ID $REGISTRY/$IMAGE_NAME
docker push $REGISTRY/$IMAGE_NAME
