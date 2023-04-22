package minecraft

type Response struct {
	Version string   `json:"version"`
	Mspt    float64  `json:"mspt"`
	Tps     TPS      `json:"tps"`
	Players []Player `json:"players"`
	Ram     RAM      `json:"ram"`
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
