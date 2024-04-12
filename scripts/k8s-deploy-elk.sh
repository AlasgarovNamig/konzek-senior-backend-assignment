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