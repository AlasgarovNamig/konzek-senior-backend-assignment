#!/bin/bash
microk8s kubectl delete --all services
microk8s kubectl delete --all deployments
microk8s kubectl delete --all configmaps
microk8s kubectl delete --all pvc
microk8s kubectl delete --all secrets