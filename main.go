package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// глобальные метрики
var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 12), // ~1ms .. ~2s
		},
		[]string{"path", "method", "status"},
	)
)

// маленький враппер, чтобы поймать status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// middleware для измерения запросов
func instrument(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		latency := time.Since(start).Seconds()
		statusCode := fmt.Sprintf("%d", rec.status)

		requestsTotal.WithLabelValues(r.URL.Path, r.Method, statusCode).Inc()
		requestDuration.WithLabelValues(r.URL.Path, r.Method, statusCode).Observe(latency)
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Всегда быстрый ответ")
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
	delay := 50 + rand.Intn(451) // 50–500 ms
	time.Sleep(time.Duration(delay) * time.Millisecond)
	fmt.Fprintf(w, "Медленный эндпоинт, задержка %d ms\n", delay)
}

func main() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)

	http.Handle("/healthz", instrument(http.HandlerFunc(homeHandler)))
	http.Handle("/", instrument(http.HandlerFunc(slowHandler)))

	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Сервер на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка сервера:", err)
	}
}
