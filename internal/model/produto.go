package model

import "time"

// Produto representa a entidade de domínio de um produto
// Esta é a entidade interna usada para persistência e lógica de negócio
type Produto struct {
	ID        int       `json:"id" bson:"id"`
	Nome      string    `json:"nome" bson:"nome"`
	Preco     float64   `json:"preco" bson:"preco"`
	Descricao string    `json:"descricao" bson:"descricao"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// BeforeCreate inicializa os timestamps antes de criar
func (p *Produto) BeforeCreate() {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
}

// BeforeUpdate atualiza o timestamp de atualização
func (p *Produto) BeforeUpdate() {
	p.UpdatedAt = time.Now()
}
