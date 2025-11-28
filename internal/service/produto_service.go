package service

import (
	"context"
	"errors"

	"api-go-arquitetura/internal/model"
	"api-go-arquitetura/internal/repository"
)

// ProdutoService implementa a lógica de negócio para produtos
type ProdutoService struct {
	repo *repository.ProdutoRepository
}

// NewProdutoService cria uma nova instância do ProdutoService
func NewProdutoService(repo *repository.ProdutoRepository) ProdutoService {
	return &ProdutoService{
		repo: repo,
	}
}

// Create cria um novo produto
func (s *ProdutoService) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	// Validações de negócio podem ser adicionadas aqui
	if produto.Nome == "" {
		return model.Produto{}, errors.New("nome do produto é obrigatório")
	}
	if produto.Preco < 0 {
		return model.Produto{}, errors.New("preço não pode ser negativo")
	}

	return s.repo.Create(ctx, produto)
}

// FindAll retorna todos os produtos
func (s *ProdutoService) FindAll(ctx context.Context) ([]model.Produto, error) {
	return s.repo.FindAll(ctx)
}

// FindByID retorna um produto pelo ID
func (s *ProdutoService) FindByID(ctx context.Context, id int) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.New("ID inválido")
	}
	return s.repo.FindByID(ctx, id)
}

// Update atualiza um produto completamente
func (s *ProdutoService) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.New("ID inválido")
	}
	if produto.Nome == "" {
		return model.Produto{}, errors.New("nome do produto é obrigatório")
	}
	if produto.Preco < 0 {
		return model.Produto{}, errors.New("preço não pode ser negativo")
	}

	return s.repo.Update(ctx, id, produto)
}

// Patch atualiza um produto parcialmente
func (s *ProdutoService) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	if id <= 0 {
		return model.Produto{}, errors.New("ID inválido")
	}

	// Validações específicas para campos que podem ser atualizados
	if nome, ok := updates["nome"].(string); ok && nome == "" {
		return model.Produto{}, errors.New("nome do produto não pode ser vazio")
	}
	if preco, ok := updates["preco"].(float64); ok && preco < 0 {
		return model.Produto{}, errors.New("preço não pode ser negativo")
	}

	return s.repo.Patch(ctx, id, updates)
}

// Delete remove um produto
func (s *ProdutoService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("ID inválido")
	}
	return s.repo.Delete(ctx, id)
}
