build:
	@go build -o bin/openapi-examples cmd/openapiexamplescli/openapiexamples.go

run:
	@go run cmd/openapiexamplescli/openapiexamples.go