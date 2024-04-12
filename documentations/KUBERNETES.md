# Kubernetes
I added the necessary configuration deployment files under deployments folder to run the projects locally.

* ## NOTE:
  After the keycloak is up, it may be necessary to get the RS256 public keys of the master and konzek realms 
    in the secret, update the secret accordingly, renew the secret and restart the deployments.
## Microk8s
I chose microk8s for easy kubernetes installation in local environment.
```
sudo snap install microk8s --classic
```

After installation, you need to activate the following plugins in microk8s to stand up the services.
```
microk8s enable dns registry ingress
```

## Upload to microk8s ( in the root directory)
* ### !!! First docker images must be loaded into the local microk8s registry (necessary information is in the Docker.md file)
``` 
/bin/bash ./scripts/k8s-deploy.sh
```

## Microk8s to remove all component connected to the project in the cluster
``` 
/bin/bash ./scripts/delete-microk8s-component.sh
```


# Container Logs
    1. Elactic search
    2. After connecting to the container, 
    you can access a folder in the app folder from the log folder outside



## Deployment restart Service Base

```
microk8s kubectl rollout restart deployment <service-name>
```

## All deployment  restart 

```
microk8s kubectl rollout restart deployment
```

