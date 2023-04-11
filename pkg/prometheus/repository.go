package prometheus

type PrometheusServer interface {
	Start()
	Stop()
}
