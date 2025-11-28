package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/service"
	"api-go-arquitetura/internal/utils"
	"api-go-arquitetura/internal/validator"
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
	ctx := r.Context()
	
	produtos, err := h.service.FindAll(ctx)
	if err != nil {
		utils.ErrorResponse(w, errors.WrapError(err, errors.ErrDatabase))
		return
	}
	
	// Converter models para DTOs
	response := dto.ToProdutoListResponse(produtos)
	
	utils.SuccessResponse(w, http.StatusOK, response)
}

// GetProduto obtém um produto por ID
// GET /api/produtos/{id}
func (h *ProdutoHandler) GetProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	ctx := r.Context()
	produto, err := h.service.FindByID(ctx, id)
	if err != nil {
		if errors.IsAPIError(err) {
			utils.ErrorResponse(w, err)
		} else {
			utils.ErrorResponse(w, errors.ErrProdutoNotFound)
		}
		return
	}

	// Converter model para DTO
	response := dto.FromModel(produto)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// CreateProduto cria um novo produto
// POST /api/produtos
func (h *ProdutoHandler) CreateProduto(w http.ResponseWriter, r *http.Request) {
	var request dto.CreateProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	created, err := h.service.Create(ctx, produto)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(created)

	utils.SuccessResponse(w, http.StatusCreated, response)
}

// UpdateProduto atualiza um produto completamente
// PUT /api/produtos/{id}
func (h *ProdutoHandler) UpdateProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	var request dto.UpdateProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para model
	produto := request.ToModel()

	ctx := r.Context()
	updated, err := h.service.Update(ctx, id, produto)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// PatchProduto atualiza um produto parcialmente
// PATCH /api/produtos/{id}
func (h *ProdutoHandler) PatchProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	var request dto.PatchProdutoRequest
	
	// Decodificar JSON
	if err := utils.DecodeJSON(r.Body, &request); err != nil {
		utils.BadRequestResponse(w, "Erro ao decodificar JSON: "+err.Error())
		return
	}

	// Validar DTO (validação opcional para PATCH)
	if validationErrors := validator.Validate(&request); len(validationErrors) > 0 {
		utils.ValidationErrorResponse(w, validationErrors)
		return
	}

	// Converter DTO para map
	updates := request.ToMap()

	ctx := r.Context()
	updated, err := h.service.Patch(ctx, id, updates)
	if err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	// Converter model para DTO de resposta
	response := dto.FromModel(updated)

	utils.SuccessResponse(w, http.StatusOK, response)
}

// DeleteProduto deleta um produto
// DELETE /api/produtos/{id}
func (h *ProdutoHandler) DeleteProduto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, errors.ErrInvalidID)
		return
	}

	ctx := r.Context()
	if err := h.service.Delete(ctx, id); err != nil {
		utils.ErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheckHandler gerencia o health check da API
type HealthCheckHandler struct {
	healthCheckFunc func(ctx context.Context) error
}

// NewHealthCheckHandler cria uma nova instância do HealthCheckHandler
func NewHealthCheckHandler(healthCheckFunc func(ctx context.Context) error) *HealthCheckHandler {
	return &HealthCheckHandler{
		healthCheckFunc: healthCheckFunc,
	}
}

// HealthCheck verifica o status da API e da conexão com o banco de dados
// GET /health
func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Verificar conexão com banco de dados se função disponível
	if h.healthCheckFunc != nil {
		if err := h.healthCheckFunc(ctx); err != nil {
			utils.JSONResponse(w, http.StatusServiceUnavailable, map[string]interface{}{
				"status":  "unhealthy",
				"message": "Conexão com banco de dados falhou",
				"error":   err.Error(),
			})
			return
		}
	}

	utils.JSONResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"message": "API e banco de dados estão funcionando",
	})
}
