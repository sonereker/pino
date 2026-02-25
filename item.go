package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Pino struct {
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	Summary   string    `json:"summary"`
	Prompt    string    `json:"prompt"`
	Plan      string    `json:"plan,omitempty"`
}

func (p Pino) Title() string       { return p.Summary }
func (p Pino) FilterValue() string { return p.Summary + " " + p.Prompt + " " + p.Plan }

func (p Pino) Description() string {
	return fmt.Sprintf("%s  %s",
		p.CreatedAt.Format("Jan 02, 2006 15:04"),
		truncate(p.Prompt, 60),
	)
}

func (p Pino) MarshalJSON() ([]byte, error) {
	type Alias Pino
	return json.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     Alias(p),
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
	})
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
