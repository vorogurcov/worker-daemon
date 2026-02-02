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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var jobDto CreateJobDto
		if err := json.NewDecoder(r.Body).Decode(&jobDto); err != nil {
			http.Error(w, "Incorrect JSON", http.StatusBadRequest)
			return
		}

		switch jobDto.Type {
		case "WaitingJob", "MonitoringCPUJob":
			// OK
		default:
			http.Error(w, "Unsupported job type", http.StatusBadRequest)
			return
		}

		if err := CreateJob(metrics, worker, jobDto); err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK"))
	}

}

func SetAndGetMux(reg *prometheus.Registry, metrics *metrics.Metrics, worker *worker.BasicWorker) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	mux.HandleFunc("/create", createJobHandlerFunc(metrics, worker))

	return mux
}
