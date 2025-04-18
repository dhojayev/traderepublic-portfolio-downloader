#!/bin/bash
set -e

# Generate the client code
$HOME/go/bin/oapi-codegen -config internal/traderepublc/api/restclient/openapi-config.yaml openapi-rest.yaml > internal/traderepublc/api/restclient/openapi.gen.go

# Format the generated code
go fmt internal/traderepublc/api/restclient/openapi.gen.go

echo "REST API client generated successfully!"
