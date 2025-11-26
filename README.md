# API Go com Arquitetura

Uma API REST completa em Go com suporte aos verbos HTTP: GET, POST, PUT, PATCH e DELETE.

## Estrutura do Projeto

```
.
├── main.go                 # Arquivo principal
├── handlers/
│   ├── routes.go          # Definição das rotas
│   └── produto.go         # Handlers dos endpoints
├── Dockerfile             # Configuração Docker
├── docker-compose.yml     # Orquestração de contêineres
├── go.mod                 # Dependências do módulo
└── go.sum                 # Checksums das dependências
```

## Endpoints

### Produtos

- **GET /api/produtos** - Listar todos os produtos
- **GET /api/produtos/{id}** - Obter um produto específico
- **POST /api/produtos** - Criar novo produto
- **PUT /api/produtos/{id}** - Atualizar produto completo
- **PATCH /api/produtos/{id}** - Atualizar produto parcialmente
- **DELETE /api/produtos/{id}** - Deletar produto

### Saúde

- **GET /health** - Verificar saúde da aplicação

## Como Executar

### Com Docker Compose

```bash
docker-compose up --build
```

A API estará disponível em: `http://localhost:8080`

### Localmente (sem Docker)

```bash
go mod download
go run main.go
```

## Exemplos de Requisições

### GET - Listar todos os produtos
```bash
curl http://localhost:8080/api/produtos
```

### GET - Obter produto específico
```bash
curl http://localhost:8080/api/produtos/1
```

### POST - Criar novo produto
```bash
curl -X POST http://localhost:8080/api/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Monitor",
    "preco": 800.00,
    "descricao": "Monitor 27 polegadas"
  }'
```

### PUT - Atualizar produto completo
```bash
curl -X PUT http://localhost:8080/api/produtos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Notebook Premium",
    "preco": 4500.00,
    "descricao": "Notebook de ultra alta performance"
  }'
```

### PATCH - Atualizar produto parcialmente
```bash
curl -X PATCH http://localhost:8080/api/produtos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "preco": 5000.00
  }'
```

### DELETE - Deletar produto
```bash
curl -X DELETE http://localhost:8080/api/produtos/1
```

### Health Check
```bash
curl http://localhost:8080/health
```

## Tecnologias

- **Go 1.21** - Linguagem de programação
- **Gorilla Mux** - Router HTTP
- **Docker** - Containerização
- **Docker Compose** - Orquestração

## Porta

A aplicação roda na porta **8080**.
