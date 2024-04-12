#!/bin/bash
readonly serviceName="$1"

# for apple silicon
unset DOCKER_DEFAULT_PLATFORM

# change the local ip for mac and windows
export KUBERNETES_CLUSTER_ADDR=localhost:32000

if [ "$serviceName" == "all" ]
then
  declare -a services=("auth-service"
                       "product-catalog-service"
                       "product-categories-service"
                      )
else
  declare -a services=("$serviceName")
fi

for service in "${services[@]}"
do
  set -e
  printf "=======> Start: %s\n" "$service"

  docker build -t "$service:latest" "$service"
  docker tag "$service:latest" "$KUBERNETES_CLUSTER_ADDR/$service:latest"
  docker push "$KUBERNETES_CLUSTER_ADDR/$service:latest"

  printf "=======> Finish: %s\n\n" "$service"
done