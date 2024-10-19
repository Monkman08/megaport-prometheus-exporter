package metrics

import (
	"net/http"
	"time"

	"github.com/Monkman08/megaport-prometheus-exporter/pkg/megaport"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define Prometheus metrics
var (
	megaportMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "megaport_metric",
			Help: "Example metric from Megaport API",
		},
		[]string{"label"},
	)
	scrapeDuration = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "megaport_scrape_duration_seconds",
			Help: "Duration of the Megaport API scrape",
		},
	)
)

// RegisterMetrics registers the metrics with the default registry
func RegisterMetrics() {
	// Register custom metrics
	prometheus.MustRegister(megaportMetric)
	prometheus.MustRegister(scrapeDuration)
}

// MetricsHandler returns an HTTP handler for all metrics
func MetricsHandler(client *megaport.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		client.GenerateToken()
		// Placeholder for actual API calls
		megaportMetric.WithLabelValues("example").Set(1)
		scrapeDuration.Set(time.Since(start).Seconds())
		promhttp.Handler().ServeHTTP(w, r)
	})
}
