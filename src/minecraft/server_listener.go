package minecraft

import (
	"metrics-exporter/src/logger"
	"sync"
	"time"
)

type Listener struct {
	log      logger.Logger
	lock     *sync.Mutex
	interval int
	mcAddr   string
	mcPort   int
}

func NewMCMetricsListener(interval int, mcAddr string, mcPort int) *Listener {
	listener := new(Listener)
	listener.log = logger.NewColorLogger("MinecraftMetrics")
	listener.lock = &sync.Mutex{}
	listener.interval = interval
	listener.mcAddr = mcAddr
	listener.mcPort = mcPort
	return listener
}

func (l *Listener) Do() {
	ticker := time.NewTicker(time.Duration(l.interval) * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			l.log.Info("aa")
		}
	}
}
