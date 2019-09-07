//go:generate swagger generate server --target=./swagger --spec=./swagger/swagger.yml --exclude-main
// Here we're specifying some flags:
// --target              the base directory for generating the files;
// --spec                path to the swagger specification;
// --exclude-main        generates only the library code and not a
//                       sample CLI application;
// --name                the name of the application.
package main

import (
	"log"

	"github.com/slonegd-otus-go/swaggertest/swagger/models"
	"github.com/slonegd-otus-go/swaggertest/swagger/restapi"
	"github.com/slonegd-otus-go/swaggertest/swagger/restapi/operations"
	apipet "github.com/slonegd-otus-go/swaggertest/swagger/restapi/operations/pet"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

func main() {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPetsAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = 8080

	var pets []*models.Pet
	var id int64 = 1

	api.PetCreateHandler = apipet.CreateHandlerFunc(
		func(params apipet.CreateParams) middleware.Responder {
			pet := params.Pet
			pet.ID = id
			id++
			pets = append(pets, pet)
			return apipet.NewCreateCreated()
		})

	api.PetListHandler = apipet.ListHandlerFunc(
		func(params apipet.ListParams) middleware.Responder {
			// TODO тут в параметрах фильтр по Kind
			return apipet.NewListOK().WithPayload(pets)
		})

	api.PetGetHandler = apipet.GetHandlerFunc(
		func(params apipet.GetParams) middleware.Responder {
			for _, pet := range pets {
				if pet.ID == params.PetID {
					return apipet.NewGetOK().WithPayload(pet)
				}
			}
			return apipet.NewGetNotFound()
		})

	// Start listening using having the handlers and port already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
