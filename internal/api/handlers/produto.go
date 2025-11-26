package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

var repo *repository.ProdutoRepository

// SetRepository injeta o repositório usado pelos handlers
func SetRepository(r *repository.ProdutoRepository) {
	repo = r
}

// GET - Listar todos os produtos
func GetProdutos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	produtos, err := repo.FindAll(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao buscar produtos"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produtos)
}

// GET - Obter produto por ID
func GetProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}
	ctx := r.Context()
	produto, err := repo.FindByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(produto)
}

// POST - Criar novo produto
func CreateProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var novoProduto model.Produto
	err := json.NewDecoder(r.Body).Decode(&novoProduto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}
	ctx := r.Context()
	created, err := repo.Create(ctx, novoProduto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Erro ao criar produto"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// PUT - Atualizar produto completo
func UpdateProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}

	var produtoAtualizado model.Produto
	err = json.NewDecoder(r.Body).Decode(&produtoAtualizado)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}
	ctx := r.Context()
	updated, err := repo.Update(ctx, id, produtoAtualizado)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// PATCH - Atualizar produto parcialmente
func PatchProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Dados inválidos"})
		return
	}
	// clean possible id field
	delete(updates, "id")
	ctx := r.Context()
	updated, err := repo.Patch(ctx, id, updates)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// DELETE - Deletar produto
func DeleteProduto(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"erro": "ID inválido"})
		return
	}
	ctx := r.Context()
	if err := repo.Delete(ctx, id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"erro": "Produto não encontrado"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Health check
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
