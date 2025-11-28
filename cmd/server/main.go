package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api-go-arquitetura/docs"
	"api-go-arquitetura/internal/api"
	"api-go-arquitetura/internal/api/handlers"
	"api-go-arquitetura/internal/api/middleware"
	"api-go-arquitetura/internal/database"
	"api-go-arquitetura/internal/repository"
	"api-go-arquitetura/internal/service"

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
	// Obter URI do MongoDB (usar MONGO_URI env var ou default)
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		log.Println("MONGO_URI não definida, usando padrão:", mongoURI)
	}

	// Conectar ao MongoDB com tratamento de erro robusto
	opts := database.DefaultConnectOptions(mongoURI)
	client, err := database.Connect(opts)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := database.Disconnect(ctx, client); err != nil {
			log.Printf("Erro ao fechar conexão com MongoDB: %v", err)
		}
	}()

	// Obter coleção de produtos com tratamento de erro
	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "api_go"
	}

	col, err := database.GetCollection(client, dbName, "produtos")
	if err != nil {
		log.Fatalf("Erro ao obter coleção: %v", err)
	}

	// Criar repositório
	prodRepo := repository.NewProdutoRepository(col)

	// Criar service e injetar o repositório
	prodService := service.NewProdutoService(prodRepo)

	// Criar handler e injetar o service
	produtoHandler := handlers.NewProdutoHandler(prodService)

	// Criar health check handler com verificação de banco de dados
	healthCheckFunc := func(ctx context.Context) error {
		return database.HealthCheck(ctx, client)
	}
	healthCheckHandler := handlers.NewHealthCheckHandler(healthCheckFunc)

	// Criar router e injetar os handlers
	router := api.NewRouter(produtoHandler, healthCheckHandler)

	// Rota do Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Aplicar middlewares
	handler := middleware.ApplyMiddlewares(router)

	// Configurar servidor HTTP
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para receber sinais do sistema
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		log.Println("Servidor iniciando na porta 8080...")
		log.Println("Swagger disponível em http://localhost:8080/swagger/index.html")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Servidor sendo encerrado...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao encerrar servidor: %v", err)
	}

	log.Println("Servidor encerrado com sucesso")
}
