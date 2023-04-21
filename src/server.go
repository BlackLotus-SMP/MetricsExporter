package src

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"metrics-exporter/src/logger"
)

type Server struct {
	engine *gin.Engine
}

func NewServer() *Server {
	server := new(Server)
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/healthcheck"), gin.Recovery())
	server.engine = engine
	routes := Loader{}
	for _, route := range routes.Load() {
		route.Route(server.engine)
	}
	return server
}

func (s *Server) Start(port string) {
	logg := logger.NewColorLogger("API")
	logg.Info("Listening on 0.0.0.0:%s", port)
	err := s.engine.Run(fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.engine
}
