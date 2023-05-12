package minecraft

type Response struct {
	Version          string                `json:"version"`
	Mspt             float64               `json:"mspt"`
	Day              float64               `json:"day"`
	Uptime           string                `json:"time_started"`
	Tps              TPS                   `json:"tps"`
	Players          []Player              `json:"players"`
	Ram              RAM                   `json:"ram"`
	DimEntities      []EntitiesPerDim      `json:"entities"`
	DimBlockEntities []BlockEntitiesPerDim `json:"block_entities"`
	DimChunks        []ChunksPerDim        `json:"chunks"`
	Dimensions       []string              `json:"dimensions"`
}

type TPS struct {
	FiveSec   float64 `json:"5s"`
	ThirtySec float64 `json:"30s"`
	OneMin    float64 `json:"1m"`
}

type Player struct {
	Name string  `json:"name"`
	UUID string  `json:"uuid"`
	Dim  string  `json:"dim"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	Z    float64 `json:"z"`
}

type RAM struct {
	Used float64 `json:"used"`
	Max  float64 `json:"max"`
}

type EntitiesPerDim struct {
	Dim      string         `json:"dim"`
	Entities []GenericCount `json:"entities"`
}

type BlockEntitiesPerDim struct {
	Dim           string         `json:"dim"`
	BlockEntities []GenericCount `json:"block_entities"`
}

type GenericCount struct {
	Name  string  `json:"name"`
	Count float64 `json:"count"`
}

type ChunksPerDim struct {
	Dim   string  `json:"dim"`
	Count float64 `json:"count"`
}
