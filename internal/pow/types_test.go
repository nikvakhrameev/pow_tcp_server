package pow_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
)

func TestChallenge_GetDifficultyString(t *testing.T) {
	ch := pow.Challenge{
		Data:       "test",
		Difficulty: 10,
	}

	require.Equal(t, "0000000000", ch.GetDifficultyString())

	ch = pow.Challenge{
		Data:       "test",
		Difficulty: 1,
	}

	require.Equal(t, "0", ch.GetDifficultyString())
}
