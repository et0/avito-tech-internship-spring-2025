.PHONY: generate migrate-up migrate-down
generate:
	@echo "Generating OpenAPI"

	oapi-codegen \
		-generate types \
		-package openapi \
		-o ./api/gen/openapi/types.go \
		./api/openapi/swagger.yaml
	
	oapi-codegen \
		-generate server \
		-package openapi \
		-o ./api/gen/openapi/server.go \
		./api/openapi/swagger.yaml
	
	oapi-codegen \
		-generate client \
		-package openapi \
		-o ./api/gen/openapi/client.go \
		./api/openapi/swagger.yaml

	@echo "Done!"

MIGRATION_URI = "postgres://postgres:postgres@localhost:5432/pvz?sslmode=disable"
migration-up:
	migrate \
		-path ./migrations \
		-database $(MIGRATION_URI) \
		up

migration-down:
	migrate \
		-path ./migrations \
		-database $(MIGRATION_URI) \
		down