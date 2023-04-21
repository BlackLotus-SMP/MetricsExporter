package thread

import (
	"metrics-exporter/src/logger"
	"runtime"
	"time"
)

type Pool struct {
	nWorkers   int
	log        logger.Logger
	workerChan chan Job
}

func NewPool(nWorkers int) *Pool {
	p := new(Pool)
	p.log = logger.NewColorLogger("Pool")
	p.nWorkers = nWorkers
	p.workerChan = make(chan Job)
	return p
}

func (p *Pool) Start(j Job) {
	for i := 0; i < p.nWorkers; i++ {
		go p.startJob(j)
	}
	for job := range p.workerChan {
		p.log.Error("worked died, restarting...")
		time.Sleep(time.Second)
		go p.startJob(job)
	}
}

func (p *Pool) startJob(j Job) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			p.log.Critical("panic! %v\n%s\n", r, buf)
		}
		p.workerChan <- j
	}()
	j.Do()
}
