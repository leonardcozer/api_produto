package repository

import (
	"context"
	"errors"

	"api-go-arquitetura/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProdutoRepository struct {
	Collection *mongo.Collection
}

func NewProdutoRepository(col *mongo.Collection) *ProdutoRepository {
	return &ProdutoRepository{Collection: col}
}

func (r *ProdutoRepository) getNextID(ctx context.Context) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{"id", -1}})
	var p model.Produto
	err := r.Collection.FindOne(ctx, bson.M{}, opts).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 1, nil
		}
		return 0, err
	}
	return p.ID + 1, nil
}

func (r *ProdutoRepository) Create(ctx context.Context, produto model.Produto) (model.Produto, error) {
	id, err := r.getNextID(ctx)
	if err != nil {
		return model.Produto{}, err
	}
	produto.ID = id
	_, err = r.Collection.InsertOne(ctx, produto)
	if err != nil {
		return model.Produto{}, err
	}
	return produto, nil
}

func (r *ProdutoRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var produtos []model.Produto
	if err = cursor.All(ctx, &produtos); err != nil {
		return nil, err
	}
	return produtos, nil
}

func (r *ProdutoRepository) FindByID(ctx context.Context, id int) (model.Produto, error) {
	var produto model.Produto
	err := r.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&produto)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Produto{}, errors.New("not found")
		}
		return model.Produto{}, err
	}
	return produto, nil
}

func (r *ProdutoRepository) Update(ctx context.Context, id int, produto model.Produto) (model.Produto, error) {
	produto.ID = id
	res, err := r.Collection.ReplaceOne(ctx, bson.M{"id": id}, produto)
	if err != nil {
		return model.Produto{}, err
	}
	if res.MatchedCount == 0 {
		return model.Produto{}, errors.New("not found")
	}
	return produto, nil
}

func (r *ProdutoRepository) Patch(ctx context.Context, id int, updates map[string]interface{}) (model.Produto, error) {
	update := bson.M{"$set": updates}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated model.Produto
	err := r.Collection.FindOneAndUpdate(ctx, bson.M{"id": id}, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Produto{}, errors.New("not found")
		}
		return model.Produto{}, err
	}
	return updated, nil
}

func (r *ProdutoRepository) Delete(ctx context.Context, id int) error {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("not found")
	}
	return nil
}
