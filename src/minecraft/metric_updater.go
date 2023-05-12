package minecraft

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricUpdater struct {
	version             *prometheus.GaugeVec
	onlinePlayers       prometheus.Gauge
	mspt                prometheus.Gauge
	day                 prometheus.Gauge
	ramUsage            *prometheus.GaugeVec
	tpsAverage          *prometheus.GaugeVec
	players             *prometheus.GaugeVec
	entities            *prometheus.GaugeVec
	blockEntities       *prometheus.GaugeVec
	chunks              *prometheus.GaugeVec
	dimensions          *prometheus.GaugeVec
	onlinePlayersUUID   []string
	loadedEntities      map[string][]string
	loadedBlockEntities map[string][]string
	loadedChunks        []string
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
			Help:      "information about players",
		},
		[]string{"name", "uuid", "dimension", "x", "y", "z"},
	)
	entities := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "entities",
			Help:      "count of entities per dimension",
		},
		[]string{"name", "dimension"},
	)
	blockEntities := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "block_entities",
			Help:      "count of block entities per dimension",
		},
		[]string{"name", "dimension"},
	)
	chunks := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "chunks",
			Help:      "count of chunks per dimension",
		},
		[]string{"dimension"},
	)
	dimensions := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "dimensions",
			Help:      "list of dimensions",
		},
		[]string{"name"},
	)
	day := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "day",
			Help:      "minecraft day",
		},
	)
	mu := new(MetricUpdater)
	mu.version = versionMetric
	mu.onlinePlayers = onlinePlayers
	mu.mspt = mspt
	mu.day = day
	mu.ramUsage = ramUsage
	mu.tpsAverage = tpsAverage
	mu.players = players
	mu.entities = entities
	mu.blockEntities = blockEntities
	mu.chunks = chunks
	mu.dimensions = dimensions

	registry.MustRegister(versionMetric)
	registry.MustRegister(onlinePlayers)
	registry.MustRegister(mspt)
	registry.MustRegister(day)
	registry.MustRegister(ramUsage)
	registry.MustRegister(tpsAverage)
	registry.MustRegister(players)
	registry.MustRegister(entities)
	registry.MustRegister(blockEntities)
	registry.MustRegister(chunks)
	registry.MustRegister(dimensions)

	return mu
}

func (mu *MetricUpdater) update(metrics Response) {
	mu.version.With(prometheus.Labels{"version": metrics.Version}).Set(1)
	mu.onlinePlayers.Set(float64(len(metrics.Players)))
	mu.mspt.Set(metrics.Mspt)
	mu.day.Set(metrics.Day)
	mu.ramUsage.With(prometheus.Labels{"data": "max"}).Set(metrics.Ram.Max)
	mu.ramUsage.With(prometheus.Labels{"data": "used"}).Set(metrics.Ram.Used)
	mu.tpsAverage.With(prometheus.Labels{"time": "5s"}).Set(metrics.Tps.FiveSec)
	mu.tpsAverage.With(prometheus.Labels{"time": "30s"}).Set(metrics.Tps.ThirtySec)
	mu.tpsAverage.With(prometheus.Labels{"time": "1m"}).Set(metrics.Tps.OneMin)
	mu.playerMetrics(metrics)
	mu.entitiesMetrics(metrics)
	mu.blockEntitiesMetrics(metrics)
	mu.chunkMetrics(metrics)
	for _, dimension := range metrics.Dimensions {
		mu.dimensions.With(prometheus.Labels{"name": dimension}).Set(1)
	}
}

func (mu *MetricUpdater) playerMetrics(metrics Response) {
	nowOnline := make(map[string]bool)
	for _, p := range metrics.Players {
		nowOnline[p.UUID] = true
	}
	for _, uuid := range mu.onlinePlayersUUID {
		if !nowOnline[uuid] {
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

func (mu *MetricUpdater) chunkMetrics(metrics Response) {
	currentChunks := make(map[string]bool)
	for _, c := range metrics.DimChunks {
		currentChunks[c.Dim] = true
	}
	for _, dim := range mu.loadedChunks {
		if !currentChunks[dim] {
			mu.chunks.Delete(prometheus.Labels{"dimension": dim})
		}
	}
	mu.loadedChunks = make([]string, 0)
	for _, c := range metrics.DimChunks {
		mu.chunks.DeletePartialMatch(prometheus.Labels{"dimension": c.Dim})
		mu.chunks.With(prometheus.Labels{"dimension": c.Dim}).Set(c.Count)
		mu.loadedChunks = append(mu.loadedChunks, c.Dim)
	}
}

func (mu *MetricUpdater) blockEntitiesMetrics(metrics Response) {
	currentBlockEntities := make(map[string]map[string]bool)
	for _, dimBlockEntities := range metrics.DimBlockEntities {
		if currentBlockEntities[dimBlockEntities.Dim] == nil {
			currentBlockEntities[dimBlockEntities.Dim] = make(map[string]bool)
		}
		for _, blockEntity := range dimBlockEntities.BlockEntities {
			currentBlockEntities[dimBlockEntities.Dim][blockEntity.Name] = true
		}
	}
	for dimension, blockEntityList := range mu.loadedBlockEntities {
		if currentBlockEntities[dimension] == nil {
			mu.blockEntities.Delete(prometheus.Labels{"dimension": dimension})
		}
		for _, blockEntity := range blockEntityList {
			if !currentBlockEntities[dimension][blockEntity] {
				mu.blockEntities.DeletePartialMatch(prometheus.Labels{"name": blockEntity, "dimension": dimension})
			}
		}
	}
	mu.loadedBlockEntities = make(map[string][]string)
	for _, dimBlockEntities := range metrics.DimBlockEntities {
		if mu.loadedBlockEntities[dimBlockEntities.Dim] == nil {
			mu.loadedBlockEntities[dimBlockEntities.Dim] = make([]string, 0)
		}
		for _, blockEntity := range dimBlockEntities.BlockEntities {
			mu.blockEntities.DeletePartialMatch(prometheus.Labels{"name": blockEntity.Name, "dimension": dimBlockEntities.Dim})
			mu.blockEntities.With(prometheus.Labels{"name": blockEntity.Name, "dimension": dimBlockEntities.Dim}).Set(blockEntity.Count)
			mu.loadedBlockEntities[dimBlockEntities.Dim] = append(mu.loadedBlockEntities[dimBlockEntities.Dim], blockEntity.Name)
		}
	}
}

func (mu *MetricUpdater) entitiesMetrics(metrics Response) {
	currentEntities := make(map[string]map[string]bool)
	for _, dimEntities := range metrics.DimEntities {
		if currentEntities[dimEntities.Dim] == nil {
			currentEntities[dimEntities.Dim] = make(map[string]bool)
		}
		for _, entity := range dimEntities.Entities {
			currentEntities[dimEntities.Dim][entity.Name] = true
		}
	}
	for dimension, entityList := range mu.loadedEntities {
		if currentEntities[dimension] == nil {
			mu.entities.Delete(prometheus.Labels{"dimension": dimension})
		}
		for _, entity := range entityList {
			if !currentEntities[dimension][entity] {
				mu.entities.DeletePartialMatch(prometheus.Labels{"name": entity, "dimension": dimension})
			}
		}
	}
	mu.loadedEntities = make(map[string][]string)
	for _, dimEntities := range metrics.DimEntities {
		if mu.loadedEntities[dimEntities.Dim] == nil {
			mu.loadedEntities[dimEntities.Dim] = make([]string, 0)
		}
		for _, entity := range dimEntities.Entities {
			mu.entities.DeletePartialMatch(prometheus.Labels{"name": entity.Name, "dimension": dimEntities.Dim})
			mu.entities.With(prometheus.Labels{"name": entity.Name, "dimension": dimEntities.Dim}).Set(entity.Count)
			mu.loadedEntities[dimEntities.Dim] = append(mu.loadedEntities[dimEntities.Dim], entity.Name)
		}
	}
}
