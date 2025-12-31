# go-microservice
$ go run tools/create_service.go -name financial

kubectl port-forward service/postgres-service -n production 5433:5432
