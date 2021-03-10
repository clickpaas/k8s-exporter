package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

var(
	crHandler *handler
	once sync.Once
)

type handler struct {
	// cache
	enableInternalMetric bool
	collectors []prometheus.Collector
}


func NewCollectorHandler(enableInternalMetric bool)*handler{
	if crHandler == nil{
		once.Do(func() {
			crHandler = &handler{
				enableInternalMetric: enableInternalMetric,
				collectors:           []prometheus.Collector{},
			}
		})
	}
	return crHandler
}

func(h *handler)GetCollectors()[]prometheus.Collector{
	return h.collectors
}

func(h *handler)MustRegister(collector ...prometheus.Collector){
	h.collectors = append(h.collectors, collector...)
}


func(h *handler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	register := prometheus.NewRegistry()

	if h.enableInternalMetric{
		register.MustRegister(prometheus.NewGoCollector())
		register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
		register.MustRegister(prometheus.NewBuildInfoCollector())
	}
	for _,collector := range h.collectors{
		register.MustRegister(collector)
	}


	sh := promhttp.HandlerFor(register, promhttp.HandlerOpts{
		Registry:            register,
		EnableOpenMetrics:   false,
	})
	sh.ServeHTTP(w, r)
}

