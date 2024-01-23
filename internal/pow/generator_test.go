package pow_test

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"slices"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	mocks "github.com/nikvakhrameev/pow_tcp_server/mocks/internal_/pow"
)

func TestGenerator_GenerateChallenge(t *testing.T) {
	challenger, _, difficultyGetter, randomDataGetter := makeGeneratorWithMocks(t)

	const (
		testRandomData = "test_data"
		testDifficulty = 10
	)

	difficultyGetter.On("GetDifficulty").Return(testDifficulty).Once()
	randomDataGetter.On("GetRandomDataBytes").Return([]byte(testRandomData), nil).Once()

	challenge, err := challenger.GenerateChallenge()
	require.NoError(t, err)
	require.Equal(t, testDifficulty, challenge.Difficulty)
	require.Equal(t, hex.EncodeToString([]byte(testRandomData)), challenge.Data)
}

func TestGenerator_GenerateChallengeError(t *testing.T) {
	challenger, _, _, randomDataGetter := makeGeneratorWithMocks(t)

	testErr := errors.New("test error")
	randomDataGetter.On("GetRandomDataBytes").Return(nil, testErr).Once()

	_, err := challenger.GenerateChallenge()
	require.ErrorIs(t, err, testErr)
}

func TestGenerator_CheckSolution(t *testing.T) {
	challenger, hasher, _, _ := makeGeneratorWithMocks(t)

	const (
		testChallengeData = "48656c6c6f20476f7068657221"
		testDifficulty    = 3

		testCorrectNonceSolution = 10
		testCorrectHexStringHash = "00056c6c6f20476f7068657221"

		testIncorrectNonceSolution = 5
		testIncorrectHexStringHash = "12156c6c6f20476f7068657221"
	)

	testDataHexBytes, err := hex.DecodeString(testChallengeData)
	require.NoError(t, err)

	correctNonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(correctNonceBytes, testCorrectNonceSolution)
	correctHashResult, err := hex.DecodeString(testCorrectHexStringHash)
	require.NoError(t, err)

	incorrectNonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(incorrectNonceBytes, testIncorrectNonceSolution)
	incorrectHashResult, err := hex.DecodeString(testIncorrectHexStringHash)
	require.NoError(t, err)

	hasher.On(
		"HashData",
		mock.MatchedBy(func(data []byte) bool {
			return slices.Equal(data, append(testDataHexBytes, correctNonceBytes...))
		}),
	).Return(correctHashResult).Once()

	hasher.On(
		"HashData",
		mock.MatchedBy(func(data []byte) bool {
			return slices.Equal(data, append(testDataHexBytes, incorrectNonceBytes...))
		}),
	).Return(incorrectHashResult).Once()

	ok, err := challenger.CheckSolution(pow.Challenge{
		Data:       testChallengeData,
		Difficulty: testDifficulty,
	}, testCorrectNonceSolution)
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = challenger.CheckSolution(pow.Challenge{
		Data:       testChallengeData,
		Difficulty: testDifficulty,
	}, testIncorrectNonceSolution)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestGenerator_SolvePowChallenge(t *testing.T) {
	challenger, hasher, _, _ := makeGeneratorWithMocks(t)

	const (
		testChallengeData = "48656c6c6f20476f7068657221"
		testDifficulty    = 3

		testIncorrectNonceSolution = uint64(0)
		testIncorrectHexStringHash = "12156c6c6f20476f7068657221"

		testCorrectNonceSolution = uint64(1)
		testCorrectHexStringHash = "00056c6c6f20476f7068657221"
	)

	testDataHexBytes, err := hex.DecodeString(testChallengeData)
	require.NoError(t, err)

	correctNonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(correctNonceBytes, testCorrectNonceSolution)
	correctHashResult, err := hex.DecodeString(testCorrectHexStringHash)
	require.NoError(t, err)

	incorrectNonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(incorrectNonceBytes, testIncorrectNonceSolution)
	incorrectHashResult, err := hex.DecodeString(testIncorrectHexStringHash)
	require.NoError(t, err)

	hasher.On(
		"HashData",
		mock.MatchedBy(func(data []byte) bool {
			return slices.Equal(data, append(testDataHexBytes, correctNonceBytes...))
		}),
	).Return(correctHashResult).Once()

	hasher.On(
		"HashData",
		mock.MatchedBy(func(data []byte) bool {
			return slices.Equal(data, append(testDataHexBytes, incorrectNonceBytes...))
		}),
	).Return(incorrectHashResult).Once()

	resNonce, err := challenger.SolvePowChallenge(context.Background(), pow.Challenge{
		Data:       testChallengeData,
		Difficulty: testDifficulty,
	})
	require.NoError(t, err)
	require.Equal(t, testCorrectNonceSolution, resNonce)
}

func TestGenerator_SolvePowChallengeCtxErr(t *testing.T) {
	challenger, _, _, _ := makeGeneratorWithMocks(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := challenger.SolvePowChallenge(ctx, pow.Challenge{
		Data:       "48656c6c6f20476f7068657221",
		Difficulty: 1,
	})
	require.ErrorIs(t, err, context.Canceled)
}

func makeGeneratorWithMocks(t *testing.T) (
	*pow.Challenger,
	*mocks.Hasher,
	*mocks.DifficultyGetter,
	*mocks.RandomDataGetter,
) {
	mockHasher := mocks.NewHasher(t)
	mockDifficultyGetter := mocks.NewDifficultyGetter(t)
	mockRandDataGetter := mocks.NewRandomDataGetter(t)
	return pow.NewChallenger(mockDifficultyGetter, mockRandDataGetter, mockHasher),
		mockHasher,
		mockDifficultyGetter,
		mockRandDataGetter
}
