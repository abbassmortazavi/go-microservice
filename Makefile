PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := pkg/proto
include .env
MIGRATION_PATH = ./migrations

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)




## Database Migrations

.PHONY: migrate-create
migrate-create:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDRESS) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDRESS) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-status
migrate-status:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDRESS) version

.PHONY: migrate-force
migrate-force:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDRESS) force $(filter-out $@,$(MAKECMDGOALS))