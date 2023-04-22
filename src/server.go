package src

import (
	"fmt"
	"metrics-exporter/src/logger"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	server := new(Server)
	mux := http.NewServeMux()
	server.mux = mux
	routes := Loader{}
	for _, route := range routes.Load() {
		route.Route(server.mux)
	}
	return server
}

func (s *Server) Start(port string) {
	logg := logger.NewColorLogger("API")
	logg.Info("Listening on 0.0.0.0:%s", port)
	panic(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), s.mux))
}
