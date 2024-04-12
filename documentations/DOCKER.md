# Docker

## Service build
By running the script below in the root directory, a single service or all services will be pushed to the registry of microk8s:
* #### Example Command (all services):
```shell
./scripts/docker-build.sh all
```

* #### Örnek komut (single service):
```shell
./scripts/docker-build.sh auth-service
```

