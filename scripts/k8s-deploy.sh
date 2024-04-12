#!/bin/bash
# Adding secret.txt file to microk8s as default secret
microk8s kubectl create secret generic default-secret --dry-run=client -o yaml --from-env-file ./deployments/k8s/secret.txt | microk8s kubectl apply -f -

# Postgres
microk8s kubectl apply -f ./deployments/k8s/postgres/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/postgres/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/postgres/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/postgres/service.yaml


# Keycloak
microk8s kubectl apply -f ./deployments/k8s/keycloak/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/keycloak/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/keycloak/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/keycloak/service.yaml

# Redis
#microk8s kubectl apply -f ./deployments/k8s/cache/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/redis/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/redis/service.yaml

# Elasticsearch
microk8s kubectl apply -f ./deployments/k8s/elasticsearch/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/elasticsearch/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/elasticsearch/service.yaml

# Logstash
microk8s kubectl apply -f ./deployments/k8s/logstash/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/logstash/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/logstash/service.yaml

# Kibana
microk8s kubectl apply -f ./deployments/k8s/kibana/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/kibana/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/kibana/service.yaml

# Auth Service
#microk8s kubectl apply -f ./deployments/k8s/product-catalog-service/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/auth-service/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/auth-service/service.yaml

# Product Catalog Service
microk8s kubectl apply -f ./deployments/k8s/product-catalog-service/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/product-catalog-service/service.yaml

# Product Category Service
microk8s kubectl apply -f ./deployments/k8s/product-category-service/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/product-category-service/service.yaml

# Prometheus
microk8s kubectl apply -f ./deployments/k8s/prometheus/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/service.yaml

# Grafana
microk8s kubectl apply -f ./deployments/k8s/grafana/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/grafana/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/grafana/service.yaml
