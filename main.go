package main

import (
	"flag"
	"metrics-exporter/src"
)

func main() {
	var port = flag.String("p", "8462", "Service port")
	server := src.NewServer()
	server.Init()
	server.Start(*port)
}
