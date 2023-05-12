package minecraft

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"metrics-exporter/src/logger"
	"metrics-exporter/src/minecraft/packet"
	"sync"
	"time"
)

type MCMetrics interface {
	GetMetrics() *Response
}

type Listener struct {
	log      logger.Logger
	lock     *sync.Mutex
	interval int
	mcAddr   string
	mcPort   uint
	metrics  *Response
	updater  *MetricUpdater
}

func NewMCMetricsListener(interval int, mcAddr string, mcPort uint, promReg *prometheus.Registry) *Listener {
	listener := new(Listener)
	listener.log = logger.NewColorLogger("MinecraftMetrics")
	listener.lock = &sync.Mutex{}
	listener.interval = interval
	listener.mcAddr = mcAddr
	listener.mcPort = mcPort
	listener.updater = NewMetricUpdater(promReg)
	return listener
}

func (l *Listener) GetMetrics() *Response {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.metrics
}

func (l *Listener) Do() {
	ticker := time.NewTicker(time.Duration(l.interval) * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	l.collect()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			l.collect()
		}
	}
}

func (l *Listener) collect() {
	dialer := Dialer{}
	conn, connErr := dialer.DialMCContext(context.Background(), l.mcAddr, l.mcPort)
	if connErr != nil {
		l.onError()
		l.log.Critical("Unable to establish a valid connection to the server %s:%d", l.mcAddr, l.mcPort)
		l.log.Critical(connErr.Error())
		return
	}
	const Handshake = 0x00
	pErr := conn.WritePacket(packet.Marshal(
		Handshake,
		packet.VarInt(762),
		packet.String(l.mcAddr),
		packet.UnsignedShort(l.mcPort),
		packet.Byte(9),
	))
	if pErr != nil {
		l.onError()
		l.log.Critical("Unable to send the metrics packet to the server %s:%d", l.mcAddr, l.mcPort)
		l.log.Critical(pErr.Error())
		return
	}
	p2Err := conn.WritePacket(packet.Marshal(0))
	if p2Err != nil {
		l.onError()
		l.log.Critical("Unable to send the metrics packet to the server %s:%d", l.mcAddr, l.mcPort)
		l.log.Critical(pErr.Error())
		return
	}
	var p packet.Packet
	readErr := conn.ReadPacket(&p)
	if readErr != nil {
		l.onError()
		l.log.Critical("Unable to read response")
		l.log.Critical(readErr.Error())
		return
	}
	var metrics Response
	unmarshalErr := json.Unmarshal(p.Data[2:], &metrics)
	if unmarshalErr != nil {
		l.onError()
		l.log.Critical("Unable to unmarshal!")
		l.log.Critical(unmarshalErr.Error())
		return
	}
	l.lock.Lock()
	log.Printf("%#v", metrics)
	l.metrics = &metrics
	l.updater.update(metrics)
	l.lock.Unlock()
}

func (l *Listener) onError() {
	l.lock.Lock()
	l.metrics = nil
	l.lock.Unlock()
}
