package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"time"
)

type Server struct {
	logger        *slog.Logger
	cfg           Config
	ddosProtector DdosProtector
}

func NewServer(cfg Config, protector DdosProtector, logger slog.Handler) *Server {
	return &Server{
		cfg:           cfg,
		ddosProtector: protector,
		logger:        slog.New(logger.WithGroup("server")),
	}
}

func (s *Server) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Port)
	if err != nil {
		return fmt.Errorf("listen for tcp on %v error: %w", s.cfg.Port, err)
	}

	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			s.logger.Error("close listener error", "err", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return context.Canceled
			}
			return fmt.Errorf("accept new connection error: %w", err)
		}

		go func() {
			if err := s.handleConnection(conn); err != nil {
				s.logger.Error("handle connection error", "err", err)
			}
		}()
	}
}

func (s *Server) handleConnection(conn net.Conn) error {
	s.logger.Info("got new connection")

	if s.cfg.HandleConnectionTimeout != 0 {
		if err := conn.SetDeadline(time.Now().Add(s.cfg.HandleConnectionTimeout)); err != nil {
			return fmt.Errorf("set connection deadline error: %w", err)
		}
	}
	defer conn.Close()

	ok, err := s.verifyConnection(conn)
	if err != nil {
		return fmt.Errorf("verify connection error: %w", err)
	}
	if !ok {
		s.logger.Warn("connection wasn't verified, reject connection")
		return nil
	}

	if err := json.NewEncoder(conn).Encode(WordOfWisdom{
		Text: "Good people are good because they've come to wisdom through failure."},
	); err != nil {
		return fmt.Errorf("write word of wisdom to connection error: %w", err)
	}

	return nil
}

const maxNonceReadBytes = 32

func (s *Server) verifyConnection(conn net.Conn) (bool, error) {
	pow, err := s.ddosProtector.GenerateChallenge()
	if err != nil {
		return false, fmt.Errorf("generate solution error: %w", err)
	}

	logger := s.logger.With("data", pow.Data, "difficulty", pow.Difficulty)

	logger.Info("pow challenge generated")

	if err := json.NewEncoder(conn).Encode(PowChallenge(pow)); err != nil {
		return false, fmt.Errorf("encode pow challenge error: %w", err)
	}

	var powSolution PowChallengeSolution
	if err := json.NewDecoder(io.LimitReader(conn, maxNonceReadBytes)).Decode(&powSolution); err != nil {
		return false, fmt.Errorf("decode pos challenge solution error: %w", err)
	}

	logger = logger.With("solution_nonce", powSolution.Nonce)

	logger.Info("got pow challenge solution")

	ok, err := s.ddosProtector.CheckSolution(pow, powSolution.Nonce)
	if err != nil {
		return false, fmt.Errorf("check solution error: %w", err)
	}

	return ok, nil
}
