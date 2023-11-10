package service

type Telemetry struct {
	Nonce  string         `json:"nonce"`
	Status map[string]any `json:"status"`
}

type TelemetryPing struct {
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
}

type TelemetryPong struct {
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
	Owl       string `json:"owl"`
}
