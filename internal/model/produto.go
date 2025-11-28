package model

// Produto representa a entidade de domínio de um produto
// Esta é a entidade interna usada para persistência e lógica de negócio
type Produto struct {
	ID        int     `json:"id" bson:"id"`
	Nome      string  `json:"nome" bson:"nome"`
	Preco     float64 `json:"preco" bson:"preco"`
	Descricao string  `json:"descricao" bson:"descricao"`
}
