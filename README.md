# go-microservice
$ go run tools/create_service.go -name financial

kubectl port-forward service/postgres-service -n production 5433:5432


# Microservice Production Deployment

## 🏗 Architecture
- **API Gateway**: Go, Port 8080
- **Auth Service**: Go + gRPC, Port 9092
- **PostgreSQL**: Version 15, Port 5432
- **RabbitMQ**: 3-management, Ports 5672/15672

## 🚀 Deployment

### Prerequisites
- Kubernetes cluster (v1.24+)
- ArgoCD installed
- GitHub Repository access
- Container Registry access

### Quick Start
```bash
# 1. Clone repository
git clone https://github.com/YOUR_USERNAME/microservice.git

# 2. Setup ArgoCD Application
kubectl apply -f argocd/applications/microservice-production.yaml

# 3. Monitor deployment
argocd app get microservice-production