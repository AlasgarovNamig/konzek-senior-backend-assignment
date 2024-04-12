# Prometheus
microk8s kubectl apply -f ./deployments/k8s/prometheus/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/configmap.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/prometheus/service.yaml

# Grafana
microk8s kubectl apply -f ./deployments/k8s/grafana/pvc.yaml
microk8s kubectl apply -f ./deployments/k8s/grafana/deployment.yaml
microk8s kubectl apply -f ./deployments/k8s/grafana/service.yaml

