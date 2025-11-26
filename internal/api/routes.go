package api

import (
	"api-go-arquitetura/internal/api/handlers"

	"github.com/gorilla/mux"
)

// NewRouter monta e retorna o router com as rotas registradas pelos handlers
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// Rotas para produtos
	router.HandleFunc("/api/produtos", handlers.GetProdutos).Methods("GET")
	router.HandleFunc("/api/produtos/{id}", handlers.GetProduto).Methods("GET")
	router.HandleFunc("/api/produtos", handlers.CreateProduto).Methods("POST")
	router.HandleFunc("/api/produtos/{id}", handlers.UpdateProduto).Methods("PUT")
	router.HandleFunc("/api/produtos/{id}", handlers.PatchProduto).Methods("PATCH")
	router.HandleFunc("/api/produtos/{id}", handlers.DeleteProduto).Methods("DELETE")

	// Rota de health check
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	return router
}
