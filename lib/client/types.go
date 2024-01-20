package client

import (
	"context"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
)

type Config struct {
	ServerUrl string `envconfig:"SERVER_URL" default:"localhost:8085"`
}

type PowChallengeSolver interface {
	SolvePowChallenge(ctx context.Context, pow pow.Challenge) (uint64, error)
}
