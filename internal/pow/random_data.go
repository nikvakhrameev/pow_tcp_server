package pow

import (
	"crypto/rand"
	"fmt"
)

type RandomDataGenerator struct {
	bytesCount int
}

func NewRandomDataGenerator(bytesCount int) RandomDataGenerator {
	return RandomDataGenerator{bytesCount: bytesCount}
}

func (rdg RandomDataGenerator) GetRandomDataBytes() ([]byte, error) {
	b := make([]byte, rdg.bytesCount)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("read random bytes error: %w", err)
	}
	return b, nil
}
