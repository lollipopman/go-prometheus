// Example of counting http requests with prometheus

package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func pollHttp() {
	httpReqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_reqs_total",
		Help: "The total number of http requests",
	})
	p := http.Client{
		Transport: NewTransportWithMetrics(http.DefaultTransport, httpReqCounter),
	}
	for {
		resp, err := p.Get("http://google.com/")
		if err != nil {
			log.Printf("ERROR %v", err)
		}
		log.Printf("CODE: %v", resp.StatusCode)

	}
}

type TransportWithMetrics struct {
	tr            http.RoundTripper
	httpReqsTotal prometheus.Counter
}

func NewTransportWithMetrics(tr http.RoundTripper, httpReqCounter prometheus.Counter) *TransportWithMetrics {
	return &TransportWithMetrics{tr, httpReqCounter}
}

func (twm *TransportWithMetrics) RoundTrip(req *http.Request) (*http.Response, error) {
	twm.httpReqsTotal.Inc()
	return twm.tr.RoundTrip(req)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	go pollHttp()
	http.ListenAndServe(":2112", nil)
}
