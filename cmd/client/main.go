package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	"github.com/nikvakhrameev/pow_tcp_server/lib/client"
)

const appName = "POW"

func main() {
	cfg := new(Config)
	cfg.fromEnv(appName)

	powSolver := pow.NewChallenger()

	ctx := context.Background()

	logHandler := slog.NewTextHandler(os.Stdout, new(slog.HandlerOptions))
	logger := slog.New(logHandler)

	cli := client.NewClient(cfg.Client, powSolver, logHandler)
	res, err := cli.GetWordOfWisdom(ctx)
	if err != nil {
		logger.Error("get word of wisdom error", "err", err)
		os.Exit(1)
	}

	logger.Info("got word of wisdom", "res", res)
}

type Config struct {
	Client client.Config `envconfig:"CLIENT"`
}

func (c *Config) fromEnv(prefix string) {
	envconfig.MustProcess(prefix, c)
}
