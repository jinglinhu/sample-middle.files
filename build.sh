#!/bin/bash

TAG=1.0.1

REPOSITORY=jinglinhu/eks-workshop-x-ray-sample-middle

docker build --tag $REPOSITORY:$TAG .

docker push $REPOSITORY:$TAG
