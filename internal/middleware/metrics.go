package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Всего HTTP запросов",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Длительность HTTP запросов в секундах",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	ActiveLinks = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_links_total",
			Help: "Всего активных ссылок",
		},
	)

	TotalClicks = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "total_clicks_all_time",
			Help: "Всего переходов по всем ссылкам",
		},
	)

	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users_total",
			Help: "Всего зарегистрированных пользователей",
		},
	)

	LoginAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "login_attempts_total",
			Help: "Всего попыток входа",
		},
		[]string{"result"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path

		HttpRequestsTotal.WithLabelValues(r.Method, path, strconv.Itoa(rw.statusCode)).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
	})
}

func RecordClick() {
	TotalClicks.Inc()
}

func SetActiveLinks(count int64) {
	ActiveLinks.Set(float64(count))
}

func SetActiveUsers(count int64) {
	ActiveUsers.Set(float64(count))
}

func RecordLoginAttempt(success bool) {
	result := "failed"
	if success {
		result = "success"
	}
	LoginAttemptsTotal.WithLabelValues(result).Inc()
}
