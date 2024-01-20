package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	"github.com/nikvakhrameev/pow_tcp_server/internal/server"
)

type Client struct {
	cfg       Config
	powSolver PowChallengeSolver
	logger    *slog.Logger
}

func NewClient(cfg Config, powSolver PowChallengeSolver, logger slog.Handler) *Client {
	return &Client{
		cfg:       cfg,
		powSolver: powSolver,
		logger:    slog.New(logger.WithGroup("client")),
	}
}

func (c *Client) GetWordOfWisdom(ctx context.Context) (string, error) {
	conn, err := net.Dial("tcp", c.cfg.ServerUrl)
	if err != nil {
		return "", fmt.Errorf("dial with server error: %w", err)
	}
	defer conn.Close()

	var pc server.PowChallenge
	if err := json.NewDecoder(conn).Decode(&pc); err != nil {
		return "", fmt.Errorf("decode server pow challenge error: %w", err)
	}

	logger := c.logger.With("pow_data", pc.Data, "pow_difficulty", pc.Difficulty)
	logger.Info("got pow challenge")

	nonce, err := c.powSolver.SolvePowChallenge(ctx, pow.Challenge(pc))
	if err != nil {
		return "", fmt.Errorf("solve pow challenge error: %w", err)
	}

	logger.Info("challenge solved", "nonce", nonce)

	if err := json.NewEncoder(conn).Encode(
		server.PowChallengeSolution{Nonce: nonce},
	); err != nil {
		return "", fmt.Errorf("encode pow challenge solution errror: %w", err)
	}

	var res server.WordOfWisdom
	if err := json.NewDecoder(conn).Decode(&res); err != nil {
		return "", fmt.Errorf("read word of wisdom error: %w", err)
	}

	return res.Text, nil
}
