package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/graphql-go/graphql"
)

type Input struct {
	Query         string                 `query:"query"`
	OperationName string                 `query:"operationName"`
	Variables     map[string]interface{} `query:"variables"`
}

func main() {

	var err error

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:1234",
	}))

	relationalDatabaseClient = NewRelationalDatabaseClient() // FIXME: Declaration is in demo.go...
	err = relationalDatabaseClient.ConnectToDatabase()
	if err != nil {
		log.Println("Cannot connect to the relational database: terminating...")
		return
	}
	err = relationalDatabaseClient.InitializeDatabase()
	if err != nil {
		log.Println("Cannot initialize the relational database: terminating...")
		return
	}

	schema, err := GenerateGraphQLSchema()
	if err != nil {
		return
	}

	// curl 'http://localhost:9090/?query=query%7Bhello%7D'
	app.Get("/", func(ctx *fiber.Ctx) error {
		var input Input
		if err := ctx.QueryParser(&input); err != nil {
			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString("Cannot parse query parameters: " + err.Error())
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  input.Query,
			OperationName:  input.OperationName,
			VariableValues: input.Variables,
		})

		ctx.Set("Content-Type", "application/graphql-response+json")
		return ctx.JSON(result)
	})

	// curl 'http://localhost:9090/' --header 'content-type: application/json' --data-raw '{"query":"query{hello}"}'
	app.Post("/", func(ctx *fiber.Ctx) error {
		var input Input
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString("Cannot parse body: " + err.Error())
		}

		result := graphql.Do(graphql.Params{
			Schema:         schema,
			RequestString:  input.Query,
			OperationName:  input.OperationName,
			VariableValues: input.Variables,
		})

		ctx.Set("Content-Type", "application/graphql-response+json")
		return ctx.JSON(result)
	})

	log.Fatal(app.Listen(":9090"))
}
