package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/TinyRogue/lembook-serv/cmd/gql/graph"
	"github.com/TinyRogue/lembook-serv/cmd/gql/graph/generated"
	"github.com/TinyRogue/lembook-serv/pkg/middleware"
	service "github.com/TinyRogue/lembook-serv/pkg/mongo"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load("../../../.env")
	port := os.Getenv("PORT")
	var mode string
	if len(os.Args) == 1 || os.Args[1] == "--dev" {
		mode = "dev"
	} else {
		mode = "prod"
	}

	log.Printf("Server running in %s mode\n", mode)
	service.InitDb()
	defer service.Disconnect()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", middleware.Cors(middleware.Auth(srv), mode))
	log.Printf("GraphiQL http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
