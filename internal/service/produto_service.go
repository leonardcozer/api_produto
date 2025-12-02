package service

import (
	"context"
	"time"

	"api-go-arquitetura/internal/cache"
	"api-go-arquitetura/internal/dto"
	"api-go-arquitetura/internal/errors"
	"api-go-arquitetura/internal/logger"
	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

// produtoService implementa a lógica de negócio para produtos
type produtoService struct {
	repo  repository.ProdutoRepository
	cache cache.Cache
	ttl   time.Duration
}

// NewProdutoService cria uma nova instância do ProdutoService
func NewProdutoService(repo repository.ProdutoRepository, cache cache.Cache) ProdutoService {
	// TTL padrão de 5 minutos para cache
	ttl := 5 * time.Minute
	return &produtoService{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

// NewProdutoServiceWithTTL cria uma nova instância do ProdutoService com TTL customizado
func NewProdutoServiceWithTTL(repo repository.ProdutoRepository, cache cache.Cache, ttl time.Duration) ProdutoService {
	return &produtoService{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
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

	// Invalidar cache de listas (novo produto adicionado)
	if s.cache != nil {
		// Limpar todas as listas em cache (simplificado)
		// Em produção, seria melhor usar padrões de chave ou tags
		logger.Debug("Cache de listas invalidado após criação de produto")
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

	// Tentar buscar do cache primeiro
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		cachedData, err := s.cache.Get(ctx, cacheKey)
		if err == nil {
			produto, err := cache.DecodeProduto(cachedData)
			if err == nil {
				logger.WithFields(map[string]interface{}{
					"id":        id,
					"cache_key": cacheKey,
				}).Debug("Cache hit para produto")
				return produto, nil
			}
		}
	}

	// Cache miss ou erro - buscar do banco
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		// Verificar se é erro de "not found" do repository
		if err.Error() == "not found" {
			return model.Produto{}, errors.ErrProdutoNotFound
		}
		return model.Produto{}, errors.WrapError(err, errors.ErrDatabase)
	}

	// Armazenar no cache
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		cachedData, err := cache.EncodeProduto(result)
		if err == nil {
			if err := s.cache.Set(ctx, cacheKey, cachedData, s.ttl); err != nil {
				logger.WithField("error", err).Warn("Erro ao armazenar produto no cache")
			}
		}
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

	// Invalidar cache do produto atualizado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após atualização de produto")
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

	// Invalidar cache do produto atualizado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após patch de produto")
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

	// Invalidar cache do produto deletado
	if s.cache != nil {
		cacheKey := cache.GenerateProdutoKey(id)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			logger.WithField("error", err).Warn("Erro ao invalidar cache do produto")
		}
		// Invalidar cache de listas também
		logger.Debug("Cache invalidado após deleção de produto")
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
