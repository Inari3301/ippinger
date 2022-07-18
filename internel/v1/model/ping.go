package model

import "time"

type PingResult struct {
	IP       string        `json:"ip"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error"`
}

type PingByCsvResult struct {
	PingResults []PingResult `json:"ping_results"`
}
