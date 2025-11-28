package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/service"
)

// ProdutoHandler gerencia os handlers de produto
type ProdutoHandler struct {
	service service.ProdutoService
}

// NewProdutoHandler cria uma nova instância do ProdutoHandler
func NewProdutoHandler(svc service.ProdutoService) *ProdutoHandler {
	return &ProdutoHandler{
		service: svc,
	}
}

// GetProdutos lista todos os produtos
// GET /api/produtos
func (h *ProdutoHandler) GetProdutos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	
	produtos, err := h.service.FindAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao buscar produtos"})
		return
	}
	
	// Converter models para DTOs
	response := dto.ToProdutoListResponse(produtos)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetProduto obtém um produto por ID
// GET /api/produtos/{id}
func (h *ProdutoHandler) GetProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}

	ctx := r.Context()
	produto, err := h.service.FindByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		return
	}

	// Converter model para DTO
	response := dto.FromModel(produto)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CreateProduto cria um novo produto
// POST /api/produtos
func (h *ProdutoHandler) CreateProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request dto.CreateProdutoRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	created, err := h.service.Create(ctx, produto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(created)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateProduto atualiza um produto completamente
// PUT /api/produtos/{id}
func (h *ProdutoHandler) UpdateProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}

	var request dto.UpdateProdutoRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	updated, err := h.service.Update(ctx, id, produto)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		}
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// PatchProduto atualiza um produto parcialmente
// PATCH /api/produtos/{id}
func (h *ProdutoHandler) PatchProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}

	var request dto.PatchProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}

	// Converter DTO para map
	updates := request.ToMap()

	ctx := r.Context()
	updated, err := h.service.Patch(ctx, id, updates)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		}
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteProduto deleta um produto
// DELETE /api/produtos/{id}
func (h *ProdutoHandler) DeleteProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"erro": err.Error()})
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck verifica o status da API
// GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
