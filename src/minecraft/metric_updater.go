package minecraft

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricUpdater struct {
	version           *prometheus.GaugeVec
	onlinePlayers     prometheus.Gauge
	mspt              prometheus.Gauge
	ramUsage          *prometheus.GaugeVec
	tpsAverage        *prometheus.GaugeVec
	players           *prometheus.GaugeVec
	onlinePlayersUUID []string
}

func NewMetricUpdater(registry *prometheus.Registry) *MetricUpdater {
	namespace := "minecraft"
	versionMetric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "version",
			Help:      "minecraft server version",
		},
		[]string{"version"},
	)
	onlinePlayers := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "online_players",
			Help:      "total players online",
		},
	)
	mspt := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "mspt",
			Help:      "milliseconds per tick",
		},
	)
	ramUsage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "total_ram",
			Help:      "total ram the sever has allocated",
		},
		[]string{"data"},
	)
	tpsAverage := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "tps",
			Help:      "average tps of the server",
		},
		[]string{"time"},
	)
	players := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "players",
			Help:      "General information about players",
		},
		[]string{"name", "uuid", "dimension", "x", "y", "z"},
	)
	mu := new(MetricUpdater)
	mu.version = versionMetric
	mu.onlinePlayers = onlinePlayers
	mu.mspt = mspt
	mu.ramUsage = ramUsage
	mu.tpsAverage = tpsAverage
	mu.players = players

	registry.MustRegister(versionMetric)
	registry.MustRegister(onlinePlayers)
	registry.MustRegister(mspt)
	registry.MustRegister(ramUsage)
	registry.MustRegister(tpsAverage)
	registry.MustRegister(players)

	return mu
}

func (mu *MetricUpdater) update(metrics Response) {
	mu.version.With(prometheus.Labels{"version": metrics.Version})
	mu.onlinePlayers.Set(float64(len(metrics.Players)))
	mu.mspt.Set(metrics.Mspt)
	mu.ramUsage.With(prometheus.Labels{"data": "max"}).Set(metrics.Ram.Max)
	mu.ramUsage.With(prometheus.Labels{"data": "used"}).Set(metrics.Ram.Used)
	mu.tpsAverage.With(prometheus.Labels{"time": "5s"}).Set(metrics.Tps.FiveSec)
	mu.tpsAverage.With(prometheus.Labels{"time": "30s"}).Set(metrics.Tps.ThirtySec)
	mu.tpsAverage.With(prometheus.Labels{"time": "1m"}).Set(metrics.Tps.OneMin)
	mu.playerMetrics(metrics)
}

func (mu *MetricUpdater) playerMetrics(metrics Response) {
	nowOnline := make(map[string]bool)
	for _, p := range metrics.Players {
		nowOnline[p.UUID] = true
	}
	for _, uuid := range mu.onlinePlayersUUID {
		if connected := nowOnline[uuid]; !connected {
			mu.players.DeletePartialMatch(prometheus.Labels{"uuid": uuid})
		}
	}
	mu.onlinePlayersUUID = make([]string, 0)
	for _, player := range metrics.Players {
		mu.players.DeletePartialMatch(prometheus.Labels{"uuid": player.UUID})
		mu.players.With(prometheus.Labels{
			"name":      player.Name,
			"uuid":      player.UUID,
			"dimension": player.Dim,
			"x":         fmt.Sprintf("%f", player.X),
			"y":         fmt.Sprintf("%f", player.Y),
			"z":         fmt.Sprintf("%f", player.Z),
		}).Set(1)
		mu.onlinePlayersUUID = append(mu.onlinePlayersUUID, player.UUID)
	}
}
