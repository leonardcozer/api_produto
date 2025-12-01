package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequestDuration é um histograma para duração de requisições HTTP
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP em segundos",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestTotal é um contador para total de requisições HTTP
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestErrors é um contador para erros HTTP
	HTTPRequestErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total de erros em requisições HTTP",
		},
		[]string{"method", "path", "status"},
	)

	// DatabaseOperations é um contador para operações de banco de dados
	DatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total de operações no banco de dados",
		},
		[]string{"operation", "collection", "status"},
	)

	// DatabaseOperationDuration é um histograma para duração de operações de banco
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_operation_duration_seconds",
			Help:    "Duração das operações de banco de dados em segundos",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"operation", "collection"},
	)

	// ActiveConnections é um gauge para conexões ativas
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_connections",
			Help: "Número de conexões HTTP ativas",
		},
	)
)

// RecordHTTPRequest registra uma requisição HTTP
func RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	status := http.StatusText(statusCode)
	if status == "" {
		status = "unknown"
	}

	HTTPRequestDuration.WithLabelValues(method, path, status).Observe(duration.Seconds())
	HTTPRequestTotal.WithLabelValues(method, path, status).Inc()

	if statusCode >= 400 {
		HTTPRequestErrors.WithLabelValues(method, path, status).Inc()
	}
}

// RecordDatabaseOperation registra uma operação de banco de dados
func RecordDatabaseOperation(operation, collection, status string, duration time.Duration) {
	DatabaseOperations.WithLabelValues(operation, collection, status).Inc()
	DatabaseOperationDuration.WithLabelValues(operation, collection).Observe(duration.Seconds())
}

// GetHandler retorna o handler do Prometheus
func GetHandler() http.Handler {
	return promhttp.Handler()
}

