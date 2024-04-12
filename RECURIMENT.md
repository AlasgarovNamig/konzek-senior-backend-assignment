# Assignment: Implementing Backend Services in Distributed Microservices Platform

## Scenario:
    You are tasked with implementing backend services in distributed microservices-based platform using Go. 
    This service should be highly scalable, fault-tolerant, and capable of handling product catalog and integration 
    with k8s container orchestrator, OIDC, distributed cache etc.

## Requirements:
### Step 1: Microservices Architecture
    Design a microservices architecture for the product catalog, including defining the services, 
    their responsibilities, and how they communicate. Implement each service as a separate Go application.

### Step 2: Service Discovery and Load Balancing
    Implement service discovery and load balancing mechanisms to ensure that services can discover and 
    communicate with each other efficiently. Consider using a tool like Consul or etcd.

### Step 3: Authentication and Authorization
    Implement a central authentication and authorization mechanisms with using well known OIDC services like GCP,
    AWS or Azure's OIDC services, or local KeyCloak.

    Define role-based access control (RBAC) to restrict access to certain services or actions.

### Step 4: Product Category Service
    Develop a product categories service that allows users to create product categories. 
    Implement search and filtering functionality for product categories.

### Step 5: Product Catalog Service
    Develop a product catalog service that allows users to create and browse categorized products and view product details.
    Implement search and filtering functionality for products.

### Step 6: Caching and Data Persistence
    Implement caching mechanisms to improve performance for frequently accessed data. 
    Choose an appropriate database (PostgreSQL or MongoDB) for data persistence. Ensure data consistency and reliability.

### Step 7: Monitoring and Logging
    Implement k8s probes for monitoring and multi-level logging to track service health, performance, and errors.

### Step 8: Deployment and Orchestration
    Prepare the microservices for deployment in a production environment. 
    Provide deployment scripts and instructions for containerization (e.g., OCI compatible container image) and 
    k8s orchestration.

### Step 9: Security
    Implement security measures to protect against common web security vulnerabilities (e.g., SQL injection, CSRF).

### Step 10: Testing
    Write comprehensive unit tests and integration tests for each microservice.

### Step 11: Documentation
    Provide clear and detailed documentation for deploying, scaling, and maintaining the platform.