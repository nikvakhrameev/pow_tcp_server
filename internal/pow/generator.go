package pow

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

type Challenger struct {
	difficultyGetter    DifficultyGetter
	randomDataGenerator RandomDataGetter
	hasher              Hasher
}

func NewChallenger(
	difficultyGetter DifficultyGetter,
	randomDataGenerator RandomDataGetter,
	hasher Hasher,
) *Challenger {
	return &Challenger{
		difficultyGetter:    difficultyGetter,
		randomDataGenerator: randomDataGenerator,
		hasher:              hasher,
	}
}

func (c *Challenger) GenerateChallenge() (Challenge, error) {
	data, err := c.randomDataGenerator.GetRandomDataBytes()
	if err != nil {
		return Challenge{}, fmt.Errorf("generate random data bytes error: %w", err)
	}
	return Challenge{
		Data:       hex.EncodeToString(data),
		Difficulty: c.difficultyGetter.GetDifficulty(),
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

	return c.validateSolution(data, challenge.GetDifficultyString()), nil
}

func (c *Challenger) SolvePowChallenge(ctx context.Context, challenge Challenge) (uint64, error) {
	data, err := hex.DecodeString(challenge.Data)
	if err != nil {
		return 0, fmt.Errorf("decode hex from string %v error: %w", challenge.Data, err)
	}

	nonceBytes := make([]byte, 8)
	data = append(data, nonceBytes...)
	data = data[:len(data)-8]

	difficultyString := challenge.GetDifficultyString()

	for i := uint64(0); i < math.MaxUint64; i++ {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}

		binary.LittleEndian.PutUint64(nonceBytes, i)
		if c.validateSolution(append(data, nonceBytes...), difficultyString) {
			return i, nil
		}
	}

	return 0, fmt.Errorf("no solution error")
}

func (c *Challenger) validateSolution(dataWithNonce []byte, difficulty string) bool {
	hash := c.hasher.HashData(dataWithNonce)
	return strings.HasPrefix(hex.EncodeToString(hash[:]), difficulty)
}
