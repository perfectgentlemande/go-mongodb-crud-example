# go-mongodb-crud-example
Example of REST API. Includes such things as MongoDB, Zap logger, chi router etc...

## Generate

### oapi-codegen

API boilerplate code is generated using `oapi-codegen` tool from the `openapi.yaml` file. See `api.go`.

Get it here:
`https://github.com/deepmap/oapi-codegen`

And make sure that your `GOPATH/bin` path presents in `PATH` variable.

Use this command to generate the `api.go` file:  
`oapi-codegen --generate=types,chi-server openapi/openapi.yaml > internal/api/api.go`