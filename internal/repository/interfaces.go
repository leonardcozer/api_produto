package repository

import (
	"context"

	"api-go-arquitetura/internal/model"
)

// ProdutoRepository define a interface para operações de produto no repositório
type ProdutoRepository interface {
	Create(ctx context.Context, produto model.Produto) (model.Produto, error)
	FindAll(ctx context.Context) ([]model.Produto, error)
	FindByID(ctx context.Context, id int) (model.Produto, error)
	Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error)
	Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error)
	Delete(ctx context.Context, id int) error
}

