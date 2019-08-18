// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	"log"
	"os"
	"fmt"
	"time"

	"github.com/alexflint/go-arg"

	"github.com/slonegd-otus-go/swaggertest/swagger/restapi"
	"github.com/slonegd-otus-go/swaggertest/swagger/restapi/operations"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 8080,
	}
)

func main() {
	arg.MustParse(args)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = args.Port

	// Implement the handler functionality.
	// As all we need to do is give an implementation to the interface
	// we can just override the `api` method giving it a method with a valid
	// signature (we didn't need to have this implementation here, it could
	// even come from a different package).
	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
		func(params operations.GetHostnameParams) middleware.Responder {
			response, _ := os.Hostname()
			return operations.NewGetHostnameOK().WithPayload(response)
		})

	api.GetTimeHandler = operations.GetTimeHandlerFunc(
		func(params operations.GetTimeParams) middleware.Responder {
			response := fmt.Sprintf("текущее время: %s\n", time.Now())
			return operations.NewGetTimeOK().WithPayload(response)
		})

	// Start listening using having the handlers and port
	// already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
