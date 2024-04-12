# Konzek Senior Assignment Task
    The project is designed to run on local and microk8s. It is recommended to run the project
    on Linux OS. You can get information about dependencies under the following headings. 
    
## Local Environment
    You can test postman konzek-local environment
### UP to Dependencies with Docker Compose

```
docker compose up -d
```

### Endpoints
    - localhost:8080 -> keycloak
    - localhost:5432 -> postgres
    - localhost:6379 -> redis
    - localhost:9200 -> elasticsearch
    - localhost:5601 -> kibana
    - localhost:9090 -> prometheus
    - localhost:3000 -> grafana

### Passwords

    keycloak.username -> admin
    keycloak.password -> admin

    db.username -> postgres
    db.password -> postgres

    grafana.password -> admin

### Local Log Directory

    auth-service -> <project_root>/logs/auth-service/*.log
    product-catalog-service -> <project_root>/logs/product-catalog-service/*.log
    product-category-service -> <project_root>/logs/product-category-service/*.log

### Kibana Log Directory

    kibana index pattern -> log-*.log
    unique key -> <service_name>

#### Note: 
    Before running the projects, it is recommended that you log in to keycloak and renew the 
    RS256 Public Keys of the master realm and konzek realm in the .env files. 
    master realm RS Public key auth service Konzek realm RS Public Key is required for product-catalog-service 
    and product-categories-service.

    For API tests you can import collections and environments under ./postman

### Auth Service
    API -> Auth Service Register User By Admin
    Detail -> You need to be a master admin user to Create User.
    Demand -> Before Call Admin Login

    When you run the Auth Service Login query, 
    it will get the token and automatically set it to the Header of the Register User By Admin API. 
    You can call the Register User By Admin api without doing anything additional.
    
    The userRoles filter in the sample query assigns which accesses accordingly
    
    If you log in with User login with the user you created for other apis, 
    the token will be automatically set in the headers.

### Search API
    SearchOperator {
        EQUAL = 0;
        NOT_EQUAL = 1;
        GREATER_THAN = 2;
        LESS_THAN = 3;
        GREATER_THAN_EQUAL = 4;
        LESS_THAN_EQUAL = 5;
    }

    MatchType {
        AND = 0;
        OR = 1;
    }

    ### BODY
    {
        "searchFields": [
            {
                "fieldName": "id", // entity's field name in the database 
                "searchIntData": 1, // The data you want to search according to the value type 
                "searchOperator": 0 // comparison operator
            },
            {
                "fieldName": "name",  // entity's field name in the database 
                "searchStringData": "Electronics", //The data you want to search according to the value type 
                "searchOperator": 0, //// comparison operator
                "matchType": 0 // Query connection type
            }
        ],
        "page": 1,
        "limit": 10
    }

### Project Run Command
Under Project Directory Run
```
go build product-catalog-service 
```
```
./product-catalog-service
```

## Microk8s Environment
    You can access it from Dokcer.md and Kubernetes.md files under documentation folder.

    You can test postman konzek-k8s environment

    Other parameters are the same as local

### Endpoints
    - 127.0.0.1:30002 -> keycloak
    - 127.0.0.1:30001 -> postgres
    - 127.0.0.1:30603 -> elasticsearch
    - 127.0.0.1:30601 -> kibana
    - 127.0.0.1:30605 -> prometheus
    - 127.0.0.1:30604 -> grafana
