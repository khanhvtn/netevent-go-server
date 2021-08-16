package api

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/khanhvtn/netevent-go/graph"
	"github.com/khanhvtn/netevent-go/graph/generated"
	"github.com/khanhvtn/netevent-go/middlewares"
)

// Defining the Graphql handler
func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.Init()}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

/* Init: set the port, cors, route, api and serve the api */
func Init() {
	// Setting up Gin
	app := gin.New()

	//load env file
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file")
	}

	//Middlewares
	app.Use(middlewares.CheckDB())
	app.Use(middlewares.ContextToContextMiddleware())

	//GraphQL
	app.POST("/query", graphqlHandler())
	app.GET("/", playgroundHandler())
	app.Run(":5000")
}
