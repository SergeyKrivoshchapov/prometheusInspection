package main

import (
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	name := os.Getenv("SERVER_NAME")
	if name == "" {
		name = "srv1"
	}

	site := os.Getenv("SITE")
	if site == "" {
		site = "dc1"
	}

	cpu := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_cpu_usage_percent",
		Help: "CPU usage percentage",
	})

	memory := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_memory_usage_percent",
		Help: "Memory usage percentage",
	})

	disk := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "server_disk_usage_percent",
		Help: "Disk usage percentage",
	})

	registry := prometheus.NewRegistry()
	registry.MustRegister(cpu, memory, disk)

	updateMetrics := func() {
		cpu.Set(10 + rand.Float64()*85)
		memory.Set(20 + rand.Float64()*70)
		disk.Set(30 + rand.Float64()*50)
	}

	updateMetrics()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			updateMetrics()
		}
	}()

	http.Handle(
		"/metrics",
		promhttp.HandleFor(registry, promhttp.HandlerOpts{}),
	)

	log.Printf("Server %s (%s), listening on :8080", name, site)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
