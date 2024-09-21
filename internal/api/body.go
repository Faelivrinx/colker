package api

import (
	"encoding/json"
)

type BodyProvider interface {
	Provide(message string) ([]byte, error)
}

type MsTeamsProvider struct{}

func (p *MsTeamsProvider) Provide(message string) ([]byte, error) {
	payload, err := json.Marshal(New(message))
	return payload, err
}
