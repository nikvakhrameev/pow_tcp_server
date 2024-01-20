package server

import (
	"encoding/json"
	"time"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
)

type Config struct {
	Port                    string        `envconfig:"PORT" default:":8085"`
	HandleConnectionTimeout time.Duration `envconfig:"HANDLE_TIMEOUT" default:"10m"`
}

type DdosProtector interface {
	GenerateChallenge() (pow.Challenge, error)
	CheckSolution(challenge pow.Challenge, nonce uint64) (bool, error)
}

type PowChallenge struct {
	Data       string `json:"data"`
	Difficulty int    `json:"difficulty"`
}

func (pc PowChallenge) encode() ([]byte, error) {
	return json.Marshal(pc)
}

type PowChallengeSolution struct {
	Nonce uint64 `json:"nonce"`
}

type WordOfWisdom struct {
	Text string `json:"text"`
}
