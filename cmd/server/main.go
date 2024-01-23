package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	"github.com/nikvakhrameev/pow_tcp_server/internal/server"
	"github.com/nikvakhrameev/pow_tcp_server/internal/wisdom"
)

const appName = "POW"

func main() {
	cfg := new(Config)
	cfg.fromEnv(appName)

	powChallenger := pow.NewChallenger(
		pow.NewDifficultyStorage(),
		pow.NewRandomDataGenerator(sha256.Size),
		pow.NewSha256Hasher(),
	)

	ctx, cancel := context.WithCancel(context.Background())

	logHandler := slog.NewTextHandler(os.Stdout, new(slog.HandlerOptions))
	logger := slog.New(logHandler)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case s := <-sigCh:
			logger.Warn("signal received, stopping", "signal", s)
			cancel()
		case <-ctx.Done():
			return
		}
	}()

	quotesStorage := wisdom.NewQuotesStorage()

	srv := server.NewServer(cfg.Server, powChallenger, quotesStorage, logHandler)

	logger.Info("run server")

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("run server error", "err", err)
		os.Exit(1)
	}
}

type Config struct {
	Server server.Config `envconfig:"SERVER"`
}

func (c *Config) fromEnv(prefix string) {
	envconfig.MustProcess(prefix, c)
}
