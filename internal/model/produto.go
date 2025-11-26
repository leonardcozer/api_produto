package model

// Produto representa um produto da API
// @Description Representa um produto da API
type Produto struct {
	ID        int     `json:"id" bson:"id" example:"1"`
	Nome      string  `json:"nome" bson:"nome" example:"Notebook"`
	Preco     float64 `json:"preco" bson:"preco" example:"3500.00"`
	Descricao string  `json:"descricao" bson:"descricao" example:"Notebook de alta performance"`
}
