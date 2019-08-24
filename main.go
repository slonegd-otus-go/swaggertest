//go:generate swagger generate server --target=./swagger --spec=./swagger/swagger.yml --exclude-main --name=hello
// Here we're specifying some flags:
// --target              the base directory for generating the files;
// --spec                path to the swagger specification;
// --exclude-main        generates only the library code and not a 
//                       sample CLI application;
// --name                the name of the application.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slonegd-otus-go/swaggertest/swagger/restapi"
	"github.com/slonegd-otus-go/swaggertest/swagger/restapi/operations"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = 8080

	// Implement the handler functionality.
	// As all we need to do is give an implementation to the interface
	// we can just override the `api` method giving it a method with a valid
	// signature (we didn't need to have this implementation here, it could
	// even come from a different package).
	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
		func(params operations.GetHostnameParams) middleware.Responder {
			response, _ := os.Hostname()
			response = fmt.Sprintf("host: %s\n", response)
			return operations.NewGetHostnameOK().WithPayload(response)
		})

	api.GetTimeHandler = operations.GetTimeHandlerFunc(
		func(params operations.GetTimeParams) middleware.Responder {
			response := fmt.Sprintf("текущее время: %s\n", time.Now())
			return operations.NewGetTimeOK().WithPayload(response)
		})

	// Start listening using having the handlers and port already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
