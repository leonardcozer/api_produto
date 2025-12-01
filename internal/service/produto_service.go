package service

import (
	"context"

	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

// produtoService implementa a lógica de negócio para produtos
type produtoService struct {
	repo repository.ProdutoRepository
}

// NewProdutoService cria uma nova instância do ProdutoService
func NewProdutoService(repo repository.ProdutoRepository) ProdutoService {
	return &produtoService{
		repo: repo,
	}
}

// Create cria um novo produto
func (s *produtoService) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	// Validações de negócio
	if produto.Nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if produto.Preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Create(ctx, produto)
	if err != nil {
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// FindAll retorna todos os produtos
func (s *produtoService) FindAll(ctx context.Context) ([]model.Produto, error) {
	result, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// FindByID retorna um produto pelo ID
func (s *produtoService) FindByID(ctx context.Context, id int) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}
	
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// Verificar se é erro de "not found" do repository
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// Update atualiza um produto completamente
func (s *produtoService) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}
	if produto.Nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if produto.Preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Update(ctx, id, produto)
	if err != nil {
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// Patch atualiza um produto parcialmente
func (s *produtoService) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.ErrInvalidID
	}

	// Validações específicas para campos que podem ser atualizados
	if nome, ok := updates["nome"].(string); ok && nome == "" {
		return model.Produto{}, errors.ErrNomeObrigatorio
	}
	if preco, ok := updates["preco"].(float64); ok && preco <= 0 {
		return model.Produto{}, errors.ErrPrecoInvalido
	}

	result, err := s.repo.Patch(ctx, id, updates)
	if err != nil {
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}
	return result, nil
}

// Delete remove um produto
func (s *produtoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.ErrInvalidID
	}
	
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			return errors.ErrProdutoNotFound
		}
		return errors.WrapError(err, errors.ErrDatabase)
	}
	return nil
}

// FindAllPaginated retorna produtos paginados com filtros
func (s *produtoService) FindAllPaginated(ctx context.Context, pagination dto.PaginationRequest, filter dto.FilterRequest) ([]model.Produto, dto.PaginationResponse, error) {
	// Validar paginação
	pagination.Validate()

	// Converter filtro para MongoDB
	mongoFilter := filter.ToMongoFilter()

	// Contar total de documentos
	totalItems, err := s.repo.Count(ctx, mongoFilter)
	if err != nil {
		return nil, dto.PaginationResponse{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Buscar produtos paginados
	produtos, err := s.repo.FindAllPaginated(ctx, pagination.GetSkip(), pagination.GetLimit(), mongoFilter)
	if err != nil {
		return nil, dto.PaginationResponse{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Criar resposta de paginação
	paginationResp := dto.NewPaginationResponse(pagination.Page, pagination.PageSize, int(totalItems))

	return produtos, paginationResp, nil
}
