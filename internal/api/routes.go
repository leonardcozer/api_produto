package api

import (
	"api-go-arquitetura/internal/api/handlers"

	"github.com/gorilla/mux"
)

// NewRouter monta e retorna o router com as rotas registradas pelos handlers
func NewRouter(produtoHandler *handlers.ProdutoHandler) *mux.Router {
	router := mux.NewRouter()

	// Rotas para produtos
	router.HandleFunc("/api/produtos", produtoHandler.GetProdutos).Methods("GET")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.GetProduto).Methods("GET")
	router.HandleFunc("/api/produtos", produtoHandler.CreateProduto).Methods("POST")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.UpdateProduto).Methods("PUT")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.PatchProduto).Methods("PATCH")
	router.HandleFunc("/api/produtos/{id}", produtoHandler.DeleteProduto).Methods("DELETE")

	// Rota de health check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	return router
}
