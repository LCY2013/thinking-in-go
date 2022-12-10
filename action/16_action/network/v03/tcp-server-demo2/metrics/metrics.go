package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// metricsHTTPPort for prometheus to connect
var metricsHTTPPort = 8889

var (
	ClientConnected prometheus.Gauge
	ReqRecvTotal    prometheus.Counter
	RespSendTotal   prometheus.Counter
)

func init() {
	ReqRecvTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_demo2_req_recv_total",
	})

	RespSendTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tcp_server_demo2_resp_send_total",
	})

	ClientConnected = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_server_demo2_client_connected",
	})

	prometheus.MustRegister(ReqRecvTotal, RespSendTotal, ClientConnected)

	// start the metrics server
	metricsServer := &http.Server{
		Addr: fmt.Sprintf(":%d", metricsHTTPPort),
	}

	mu := http.NewServeMux()
	mu.Handle("/metrics", promhttp.Handler())
	metricsServer.Handler = mu
	go func() {
		err := metricsServer.ListenAndServe()
		if err != nil {
			fmt.Printf("prometheus-exporter http server start failed: [%+v]\n", err)
		}
	}()
	fmt.Println("prometheus metrics server start ok(*:88899)")
}
