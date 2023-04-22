package main

import (
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"metrics-exporter/src"
	"metrics-exporter/src/minecraft"
	"metrics-exporter/src/thread"
)

func main() {
	var port = flag.String("p", "8462", "Service port")
	var interval = flag.Int("interval", 5, "Every how many seconds the service will collect metrics from the minecraft server")
	var serverAddress = flag.String("mcAddress", "127.0.0.1", "The minecraft server address (just ip/dns)")
	var serverPort = flag.Uint("mcPort", 25565, "The minecraft server port")
	flag.Parse()
	if *serverPort < 0 || *serverPort > 65536 {
		panic("Invalid minecraft server port!")
	}

	promReg := prometheus.NewRegistry()

	minecraftListenerPool := thread.NewPool(1)
	minecraftListener := minecraft.NewMCMetricsListener(*interval, *serverAddress, *serverPort, promReg)
	go minecraftListenerPool.Start(minecraftListener)

	server := src.NewServer(promReg)
	server.Start(*port)
}
