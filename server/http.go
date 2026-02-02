package server

import (
	"encoding/json"
	"main/job/metrics"
	"main/worker"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func createJobHandlerFunc(metrics *metrics.Metrics, worker *worker.BasicWorker) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			w.Write([]byte("method not allowed"))
			return
		}
		var jobDto CreateJobDto
		if err := json.NewDecoder(r.Body).Decode(&jobDto); err != nil {
			w.WriteHeader(500)
			w.Write([]byte("incorrect json"))
			return
		}

		if jobDto.Type != "WaitingJob" && jobDto.Type != "MonitoringCPUJob" {
			w.WriteHeader(500)
			w.Write([]byte("incorrect job type, types supported are 'MonitoringCPUJob' and 'WaitingJob'"))
			return
		}

		if err := CreateJob(metrics, worker, jobDto); err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("OK"))
	}

}

func SetAndGetMux(reg *prometheus.Registry, metrics *metrics.Metrics, worker *worker.BasicWorker) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	mux.HandleFunc("/create", createJobHandlerFunc(metrics, worker))

	return mux
}
