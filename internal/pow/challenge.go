package pow

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

const (
	powDifficulty = 7
)

var powDifficultyString = strings.Repeat("0", powDifficulty)

type Challenger struct{}

func NewChallenger() *Challenger {
	return new(Challenger)
}

func (c *Challenger) GenerateChallenge() (Challenge, error) {
	data, err := c.generateRandomDataBytes()
	if err != nil {
		return Challenge{}, fmt.Errorf("generate random data bytes error: %w", err)
	}
	return Challenge{
		Data:       hex.EncodeToString(data),
		Difficulty: powDifficulty,
	}, nil
}

func (c *Challenger) CheckSolution(challenge Challenge, nonce uint64) (bool, error) {
	data, err := hex.DecodeString(challenge.Data)
	if err != nil {
		return false, fmt.Errorf("decode hex from string %v error: %w", challenge.Data, err)
	}

	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)

	data = append(data, nonceBytes...)

	return c.validateSolution(data), nil
}

func (c *Challenger) SolvePowChallenge(ctx context.Context, challenge Challenge) (uint64, error) {
	data, err := hex.DecodeString(challenge.Data)
	if err != nil {
		return 0, fmt.Errorf("decode hex from string %v error: %w", challenge.Data, err)
	}

	nonceBytes := make([]byte, 8)
	data = append(data, nonceBytes...)
	data = data[:len(data)-8]

	for i := uint64(0); i < math.MaxUint64; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		binary.LittleEndian.PutUint64(nonceBytes, i)
		if c.validateSolution(append(data, nonceBytes...)) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("no solution error")
}

func (c *Challenger) validateSolution(dataWithNonce []byte) bool {
	hash := sha256.Sum256(dataWithNonce)
	return strings.HasPrefix(hex.EncodeToString(hash[:]), powDifficultyString)
}

func (c *Challenger) generateRandomDataBytes() ([]byte, error) {
	b := make([]byte, sha256.Size)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read random bytes error: %w", err)
	}
	return b, nil
}
