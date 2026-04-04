package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Metrics tracks request counts and latencies per endpoint.
type Metrics struct {
	mu       sync.RWMutex
	counters map[string]*endpointMetrics
	total    atomic.Int64
	errors   atomic.Int64
	start    time.Time
}

type endpointMetrics struct {
	requests atomic.Int64
	errors   atomic.Int64
	totalMs  atomic.Int64
}

// NewMetrics creates a new metrics tracker.
func NewMetrics() *Metrics {
	return &Metrics{
		counters: make(map[string]*endpointMetrics),
		start:    time.Now(),
	}
}

func (m *Metrics) getOrCreate(key string) *endpointMetrics {
	m.mu.RLock()
	em, ok := m.counters[key]
	m.mu.RUnlock()
	if ok {
		return em
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if em, ok = m.counters[key]; ok {
		return em
	}
	em = &endpointMetrics{}
	m.counters[key] = em
	return em
}

// Middleware returns a gin middleware that records request metrics.
func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start).Milliseconds()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		key := c.Request.Method + " " + path

		em := m.getOrCreate(key)
		em.requests.Add(1)
		em.totalMs.Add(duration)
		m.total.Add(1)

		if c.Writer.Status() >= 400 {
			em.errors.Add(1)
			m.errors.Add(1)
		}
	}
}

// Handler returns a gin handler that serves Prometheus-compatible metrics.
func (m *Metrics) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sb strings.Builder
		uptime := time.Since(m.start).Seconds()

		sb.WriteString("# HELP logicaproxy_uptime_seconds Proxy uptime in seconds\n")
		sb.WriteString("# TYPE logicaproxy_uptime_seconds gauge\n")
		sb.WriteString(fmt.Sprintf("logicaproxy_uptime_seconds %.2f\n\n", uptime))

		sb.WriteString("# HELP logicaproxy_requests_total Total requests\n")
		sb.WriteString("# TYPE logicaproxy_requests_total counter\n")
		sb.WriteString(fmt.Sprintf("logicaproxy_requests_total %d\n\n", m.total.Load()))

		sb.WriteString("# HELP logicaproxy_errors_total Total error responses (4xx/5xx)\n")
		sb.WriteString("# TYPE logicaproxy_errors_total counter\n")
		sb.WriteString(fmt.Sprintf("logicaproxy_errors_total %d\n\n", m.errors.Load()))

		sb.WriteString("# HELP logicaproxy_endpoint_requests_total Requests per endpoint\n")
		sb.WriteString("# TYPE logicaproxy_endpoint_requests_total counter\n")
		m.mu.RLock()
		for key, em := range m.counters {
			parts := strings.SplitN(key, " ", 2)
			method, path := parts[0], parts[1]
			reqs := em.requests.Load()
			errs := em.errors.Load()
			avgMs := int64(0)
			if reqs > 0 {
				avgMs = em.totalMs.Load() / reqs
			}
			sb.WriteString(fmt.Sprintf("logicaproxy_endpoint_requests_total{method=%q,path=%q} %d\n", method, path, reqs))
			sb.WriteString(fmt.Sprintf("logicaproxy_endpoint_errors_total{method=%q,path=%q} %d\n", method, path, errs))
			sb.WriteString(fmt.Sprintf("logicaproxy_endpoint_avg_latency_ms{method=%q,path=%q} %d\n", method, path, avgMs))
		}
		m.mu.RUnlock()

		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(sb.String()))
	}
}
