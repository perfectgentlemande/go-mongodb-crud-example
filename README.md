# go-mongodb-crud-example
Example of REST API.  
Includes:  
- oapi-codegen
- go-chi
- MongoDB, Zap logger, etc...

## Generate

API boilerplate code is generated using `oapi-codegen` tool from the `openapi.yaml` file.  
It's great tool that makes your actual API reflect the documentation.  

Get it there:
`https://github.com/deepmap/oapi-codegen`

And make sure that your `GOPATH/bin` path presents in `PATH` variable.

Use this command to generate the `api.go` file:  
`oapi-codegen --package=api --generate=types,chi-server openapi/openapi.yaml > internal/api/api.go`