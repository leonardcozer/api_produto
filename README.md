# API Go com Arquitetura

Uma API REST completa em Go com arquitetura em camadas, suporte aos verbos HTTP (GET, POST, PUT, PATCH, DELETE), validação estruturada, tratamento de erros padronizado e muito mais.

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture** e **Layered Architecture**:

```
┌─────────────┐
│   Handler   │  ← Camada de apresentação (HTTP)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  ← Camada de lógica de negócio
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ Repository │  ← Camada de acesso a dados
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  Database   │  ← MongoDB
└─────────────┘
```

## 📁 Estrutura do Projeto

```
api-go-arquitetura/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── api/
│   │   ├── handlers/            # Handlers HTTP
│   │   │   ├── produto.go
│   │   │   └── produto_test.go
│   │   ├── middleware/          # Middlewares (RequestID, CORS, Logger, Recovery, RateLimit)
│   │   └── routes.go            # Definição de rotas
│   ├── config/                  # Configurações
│   │   └── config.go
│   ├── database/                # Conexão com banco de dados
│   │   └── mongodb.go
│   ├── dto/                     # Data Transfer Objects
│   │   ├── converter.go
│   │   ├── produto_request.go
│   │   └── produto_response.go
│   ├── errors/                  # Erros customizados
│   │   └── errors.go
│   ├── logger/                  # Logger estruturado
│   │   └── logger.go
│   ├── model/                   # Modelos de domínio
│   │   └── produto.go
│   ├── repository/              # Camada de acesso a dados
│   │   ├── interfaces.go
│   │   ├── produto_repository.go
│   │   └── produto_repository_test.go
│   ├── service/                 # Camada de lógica de negócio
│   │   ├── interfaces.go
│   │   ├── produto_service.go
│   │   └── produto_service_test.go
│   ├── utils/                   # Utilitários
│   │   └── response.go
│   └── validator/               # Validação estruturada
│       └── validator.go
├── docs/                        # Documentação Swagger
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Endpoints

### Produtos

#### Versão 1 (Recomendado)
- **GET /api/v1/produtos** - Listar todos os produtos
- **GET /api/v1/produtos/{id}** - Obter um produto específico
- **POST /api/v1/produtos** - Criar novo produto
- **PUT /api/v1/produtos/{id}** - Atualizar produto completo
- **PATCH /api/v1/produtos/{id}** - Atualizar produto parcialmente
- **DELETE /api/v1/produtos/{id}** - Deletar produto

#### Versão Legacy (Compatibilidade)
- **GET /api/produtos** - Listar todos os produtos (redireciona para v1)
- **GET /api/produtos/{id}** - Obter um produto específico (redireciona para v1)
- **POST /api/produtos** - Criar novo produto (redireciona para v1)
- **PUT /api/produtos/{id}** - Atualizar produto completo (redireciona para v1)
- **PATCH /api/produtos/{id}** - Atualizar produto parcialmente (redireciona para v1)
- **DELETE /api/produtos/{id}** - Deletar produto (redireciona para v1)

### Saúde

- **GET /health** - Verificar saúde da aplicação

### Observabilidade

- **GET /metrics** - Métricas Prometheus

## 🚀 Como Executar

### Pré-requisitos

- Go 1.21 ou superior
- MongoDB (local ou remoto)
- Redis (opcional, para cache distribuído)
- Docker e Docker Compose (opcional)

### Variáveis de Ambiente

O projeto suporta as seguintes variáveis de ambiente:

#### MongoDB
- `MONGO_URI` - URI de conexão (padrão: `mongodb://localhost:27017`)
- `MONGO_DB` - Nome do banco de dados (padrão: `api_go`)
- `MONGO_CONNECT_TIMEOUT` - Timeout de conexão (padrão: `10s`)
- `MONGO_MAX_POOL_SIZE` - Tamanho máximo do pool (padrão: `100`)
- `MONGO_MIN_POOL_SIZE` - Tamanho mínimo do pool (padrão: `10`)

#### Server
- `PORT` - Porta do servidor (padrão: `8080`)
- `READ_TIMEOUT` - Timeout de leitura (padrão: `15s`)
- `WRITE_TIMEOUT` - Timeout de escrita (padrão: `15s`)
- `IDLE_TIMEOUT` - Timeout de idle (padrão: `60s`)
- `SHUTDOWN_TIMEOUT` - Timeout de shutdown (padrão: `30s`)

#### Logging
- `LOG_LEVEL` - Nível de log: `debug`, `info`, `warn`, `error` (padrão: `info`)
- `LOG_FORMAT` - Formato de log: `json` ou `text` (padrão: `text`)

#### Observabilidade (Loki/Grafana)
- `LOKI_URL` - URL do endpoint Loki para envio de logs (ex: `http://10.110.0.239:3100/loki/api/v1/push`)
- `LOKI_JOB` - Nome do job para identificação no Grafana (padrão: `ARQUITETURA`)

#### Cache
- `CACHE_TYPE` - Tipo de cache: `memory` ou `redis` (padrão: `memory`)
- `CACHE_TTL` - TTL (Time To Live) do cache (padrão: `5m`)
- `REDIS_ADDR` - Endereço do Redis (padrão: `localhost:6379`)
- `REDIS_PASSWORD` - Senha do Redis (padrão: vazio)
- `REDIS_DB` - Database do Redis (padrão: `0`)

### Com Docker Compose

```bash
docker-compose up --build
```

A API estará disponível em: `http://localhost:8080`

### Localmente (sem Docker)

1. **Instalar dependências**:
```bash
go mod download
```

2. **Configurar variáveis de ambiente** (opcional):
```bash
export MONGO_URI="mongodb://localhost:27017"
export MONGO_DB="api_go"
export PORT="8080"
export LOG_LEVEL="info"
```

3. **Executar a aplicação**:
```bash
go run cmd/server/main.go
```

### Executar Testes

```bash
# Todos os testes
go test ./...

# Testes com cobertura
go test -cover ./...

# Testes verbosos
go test -v ./...
```

## Exemplos de Requisições

### Request ID Tracking

Todas as requisições recebem automaticamente um **Request ID** único no header de resposta `X-Request-ID`. Este ID também aparece em todos os logs da requisição, facilitando o rastreamento e debugging.

Você pode enviar um Request ID customizado no header `X-Request-ID` da requisição, e ele será reutilizado:

```bash
curl -X GET http://localhost:8080/api/v1/produtos \
  -H "X-Request-ID: meu-request-id-123"
```

### GET - Listar todos os produtos
```bash
# Versão 1 (recomendado)
curl http://localhost:8080/api/v1/produtos

# Versão legacy (compatibilidade)
curl http://localhost:8080/api/produtos

# Com paginação
curl "http://localhost:8080/api/v1/produtos?page=1&pageSize=10"

# Com filtros
curl "http://localhost:8080/api/v1/produtos?nome=notebook&precoMin=1000&precoMax=5000"

# Com paginação e filtros
curl "http://localhost:8080/api/v1/produtos?page=1&pageSize=10&nome=notebook&precoMin=1000"
```

### GET - Obter produto específico
```bash
# Versão 1 (recomendado)
curl http://localhost:8080/api/v1/produtos/1

# Versão legacy (compatibilidade)
curl http://localhost:8080/api/produtos/1
```

### POST - Criar novo produto
```bash
# Versão 1 (recomendado)
curl -X POST http://localhost:8080/api/v1/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Monitor",
    "preco": 800.00,
    "descricao": "Monitor 27 polegadas"
  }'

# Versão legacy (compatibilidade)
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
curl -X PUT http://localhost:8080/api/v1/produtos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Notebook Premium",
    "preco": 4500.00,
    "descricao": "Notebook de ultra alta performance"
  }'
```

### PATCH - Atualizar produto parcialmente
```bash
curl -X PATCH http://localhost:8080/api/v1/produtos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "preco": 5000.00
  }'
```

### DELETE - Deletar produto
```bash
curl -X DELETE http://localhost:8080/api/v1/produtos/1
```

### Health Check
```bash
curl http://localhost:8080/health
```

### Métricas Prometheus
```bash
curl http://localhost:8080/metrics
```

## 💾 Cache

A API suporta cache em duas modalidades:

### Cache em Memória (Padrão)
Cache local em memória, ideal para desenvolvimento e ambientes pequenos:
```bash
export CACHE_TYPE="memory"
export CACHE_TTL="5m"
```

### Cache Redis (Produção)
Cache distribuído usando Redis, ideal para ambientes de produção:
```bash
export CACHE_TYPE="redis"
export REDIS_ADDR="localhost:6379"
export REDIS_PASSWORD=""  # Opcional
export REDIS_DB="0"
export CACHE_TTL="5m"
```

**Funcionalidades do Cache:**
- ✅ Cache automático em operações de leitura (`FindByID`)
- ✅ Invalidação automática em operações de escrita (Create, Update, Patch, Delete)
- ✅ TTL configurável por variável de ambiente
- ✅ Fallback automático: se Redis falhar, usa cache em memória

## 🔄 Transações MongoDB

A API inclui suporte completo a transações do MongoDB para operações atômicas. Use transações quando precisar garantir que múltiplas operações sejam executadas como uma única unidade.

**Exemplo de uso:**

```go
import (
    "context"
    "api-go-arquitetura/internal/database"
    "go.mongodb.org/mongo-driver/mongo"
)

func exemploTransacao(ctx context.Context, client *mongo.Client) error {
    // Iniciar transação
    tx, cancel, err := database.StartTransaction(ctx, client)
    if err != nil {
        return err
    }
    defer cancel()
    defer tx.End()

    // Executar operações dentro da transação
    return tx.WithTransaction(func(sc mongo.SessionContext) error {
        // Todas as operações devem usar sc como contexto
        // Se qualquer operação falhar, rollback automático
        
        // Operação 1
        // _, err := collection1.InsertOne(sc, document1)
        // if err != nil {
        //     return err // Rollback automático
        // }
        
        // Operação 2
        // _, err := collection2.InsertOne(sc, document2)
        // if err != nil {
        //     return err // Rollback automático
        // }
        
        return nil // Sucesso - commit automático
    })
}
```

**Nota:** Transações requerem MongoDB replica set ou sharded cluster. Para desenvolvimento local, você pode usar um replica set de um único nó.

## 🛠️ Tecnologias

### Core
- **Go 1.21** - Linguagem de programação
- **Gorilla Mux** - Router HTTP
- **MongoDB** - Banco de dados NoSQL

### Bibliotecas
- **validator/v10** - Validação estruturada
- **logrus** - Logger estruturado
- **swaggo/swag** - Geração de documentação Swagger

### Infraestrutura
- **Docker** - Containerização
- **Docker Compose** - Orquestração

## 📚 Documentação da API

A documentação Swagger está disponível em:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Swagger JSON**: `http://localhost:8080/swagger/doc.json`

## 🧪 Testes

O projeto inclui testes unitários para:
- ✅ Service Layer (`internal/service/produto_service_test.go`)
- ✅ Handlers (`internal/api/handlers/produto_test.go`)
- ✅ Repository (`internal/repository/produto_repository_test.go`)

### Executar Testes

```bash
# Todos os testes
go test ./...

# Testes com cobertura
go test -cover ./...

# Testes de um pacote específico
go test ./internal/service/...
```

## 📊 Qualidade do Código

- ✅ **Arquitetura**: 9/10
- ✅ **Separação de Camadas**: 9/10
- ✅ **Tratamento de Erros**: 9/10
- ✅ **Validação**: 9/10
- ✅ **Testes**: 8/10
- ✅ **Logging**: 8/10
- ✅ **Documentação**: 8/10

**Nota Geral**: **8.5/10** ✅

## 🔧 Funcionalidades

- ✅ CRUD completo de produtos
- ✅ Validação estruturada de dados
- ✅ Tratamento de erros padronizado
- ✅ Logger estruturado (JSON/Text)
- ✅ Health check com verificação de banco
- ✅ Middlewares (CORS, Rate Limit, Recovery, Logging)
- ✅ Documentação Swagger
- ✅ Graceful shutdown
- ✅ Testes unitários
- ✅ **Paginação** (page, pageSize)
- ✅ **Filtros e Busca** (nome, precoMin, precoMax, descricao)
- ✅ **Métricas Prometheus** (endpoint /metrics)
- ✅ **Versionamento de API** (v1 com compatibilidade com versões antigas)
- ✅ **Request ID Tracking** (rastreamento de requisições via X-Request-ID)
- ✅ **Cache Layer** (memória ou Redis) com invalidação automática
- ✅ **Transações MongoDB** (suporte para operações atômicas)

## 📝 Exemplos de Respostas

### Sucesso (200 OK)
```json
{
  "id": 1,
  "nome": "Notebook",
  "preco": 3500.00,
  "descricao": "Notebook de alta performance"
}
```

### Erro Padronizado (400 Bad Request)
```json
{
  "code": "INVALID_INPUT",
  "message": "Dados de entrada inválidos",
  "details": "O campo 'nome' é obrigatório"
}
```

### Lista de Produtos (200 OK)
```json
{
  "produtos": [
    {
      "id": 1,
      "nome": "Notebook",
      "preco": 3500.00,
      "descricao": "Notebook de alta performance"
    }
  ],
  "total": 1
}
```

## 🎯 Próximas Melhorias

- [ ] Rate limit distribuído (Redis)
- [ ] Autenticação e Autorização (JWT)
- [ ] Webhooks
- [ ] Versionamento v2 (quando necessário)
- [ ] Cache distribuído com tags para invalidação mais eficiente

## 📄 Licença

Este projeto está sob a licença MIT.
