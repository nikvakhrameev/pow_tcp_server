package pow

import "strings"

type Challenge struct {
	Data       string
	Difficulty int
}

func (ch Challenge) GetDifficultyString() string {
	return strings.Repeat("0", ch.Difficulty)
}

type Hasher interface {
	HashData(data []byte) []byte
}

type DifficultyGetter interface {
	GetDifficulty() int
}

type RandomDataGetter interface {
	GetRandomDataBytes() ([]byte, error)
}
