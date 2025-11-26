package main

import (
	"log"
	"net/http"
	"os"

	_ "api-go-arquitetura/docs"
	"api-go-arquitetura/internal/api"
	"api-go-arquitetura/internal/api/handlers"
	"api-go-arquitetura/internal/api/middleware"
	"api-go-arquitetura/internal/database"
	"api-go-arquitetura/internal/repository"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title API Go com Arquitetura
// @version 1.0
// @description Uma API REST completa em Go com suporte aos verbos HTTP
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @basePath /api
// @schemes http
func main() {
	// Conectar ao MongoDB (usar MONGO_URI env var ou default)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	client := database.Connect(mongoURI)

	// obter coleção de produtos
	col := client.Database("api_go").Collection("produtos")

	// criar repositório e injetar nos handlers
	prodRepo := repository.NewProdutoRepository(col)
	handlers.SetRepository(prodRepo)

	router := api.NewRouter()

	// Rota do Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Aplicar middlewares
	handler := middleware.ApplyMiddlewares(router)

	log.Println("Servidor iniciando na porta 8080...")
	log.Println("Swagger disponível em http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
