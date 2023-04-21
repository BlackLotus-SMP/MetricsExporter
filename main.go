package main

import (
	"flag"
	"metrics-exporter/src"
	"metrics-exporter/src/minecraft"
	"metrics-exporter/src/thread"
)

func main() {
	var port = flag.String("p", "8462", "Service port")
	var interval = flag.Int("interval", 15, "Every how many seconds the service will collect metrics from the minecraft server")
	var serverAddress = flag.String("mcAddress", "127.0.0.1", "The minecraft server address (just ip/dns)")
	var serverPort = flag.Int("mcPort", 25565, "The minecraft server port")
	if *serverPort < 0 || *serverPort > 65535 {
		panic("Invalid minecraft server port!")
	}

	minecraftListenerPool := thread.NewPool(1)
	minecraftListener := minecraft.NewMCMetricsListener(*interval, *serverAddress, *serverPort)
	go minecraftListenerPool.Start(minecraftListener)

	server := src.NewServer()
	server.Start(*port)
}
