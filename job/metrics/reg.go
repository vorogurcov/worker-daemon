package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	CpuUsagePercent   prometheus.Gauge
	MemUsagePercent   prometheus.Gauge
	DiskCUsagePercent prometheus.Gauge
	DiskDUsagePercent prometheus.Gauge
	NetCounter        prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		CpuUsagePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "win_cpu_usage_percent",
			Help: "Current usage of the CPU.",
		}),
		MemUsagePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "win_mem_usage_percent",
			Help: "Current usage of the Memory.",
		}),
		DiskCUsagePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "win_disk_c_usage_percent",
			Help: "Current usage of the Disk C:.",
		}),
		DiskDUsagePercent: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "win_disk_d_usage_percent",
			Help: "Current usage of the Disk D:.",
		}),
		NetCounter: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "win_net_mib_sent_total",
			Help: "Total MiB sent through Net",
		}),
	}
	reg.MustRegister(m.CpuUsagePercent)
	reg.MustRegister(m.MemUsagePercent)
	reg.MustRegister(m.DiskCUsagePercent)
	reg.MustRegister(m.DiskDUsagePercent)
	reg.MustRegister(m.NetCounter)

	return m
}
