package server

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nikvakhrameev/pow_tcp_server/internal/pow"
	mocks "github.com/nikvakhrameev/pow_tcp_server/mocks/internal_/server"
)

// TODO: client side tests should be implemented in server side style

type srvTestCase struct {
	Name                       string
	GenerateChallengeError     error
	GeneratedChallenge         pow.Challenge
	ClientChallengeSolutionRaw string
	ChallengeSolutionCorrect   *bool
	ClientSolutionNonce        uint64
	CheckSolutionError         error
	Quote                      *WordOfWisdom
	HandleConnErrExpected      bool
}

func TestServer_HandleConnection(t *testing.T) {
	testCases := []srvTestCase{
		{
			Name:                   "solved_pow_no_errors",
			GenerateChallengeError: nil,
			GeneratedChallenge: pow.Challenge{
				Data:       "test_data",
				Difficulty: 10,
			},
			ClientChallengeSolutionRaw: `{"nonce":10}`,
			ClientSolutionNonce:        10,
			ChallengeSolutionCorrect: func() *bool {
				t := true
				return &t
			}(),
			CheckSolutionError:    nil,
			Quote:                 &WordOfWisdom{Text: "test quote"},
			HandleConnErrExpected: false,
		},
		{
			Name:                   "incorrect_solution",
			GenerateChallengeError: nil,
			GeneratedChallenge: pow.Challenge{
				Data:       "test_data",
				Difficulty: 10,
			},
			ClientChallengeSolutionRaw: `{"nonce":20}`,
			ClientSolutionNonce:        20,
			ChallengeSolutionCorrect: func() *bool {
				t := false
				return &t
			}(),
			CheckSolutionError:    nil,
			Quote:                 nil,
			HandleConnErrExpected: false,
		},
		{
			Name:                     "create_challenge_error",
			GenerateChallengeError:   errors.New("test"),
			ChallengeSolutionCorrect: nil,
			CheckSolutionError:       nil,
			Quote:                    nil,
			HandleConnErrExpected:    true,
		},
		{
			Name:                   "invalid_solution_data_from_client_error",
			GenerateChallengeError: nil,
			GeneratedChallenge: pow.Challenge{
				Data:       "test_data",
				Difficulty: 10,
			},
			ClientChallengeSolutionRaw: `invalid data`,
			ChallengeSolutionCorrect:   nil,
			HandleConnErrExpected:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			srv, mockWisdomQuotes, mockDdosProtector := makeServerWithMocks(t)

			mockDdosProtector.On("GenerateChallenge").
				Return(tc.GeneratedChallenge, tc.GenerateChallengeError).Once()

			if tc.ChallengeSolutionCorrect != nil {
				mockDdosProtector.On("CheckSolution", tc.GeneratedChallenge, tc.ClientSolutionNonce).
					Return(*tc.ChallengeSolutionCorrect, tc.CheckSolutionError).Once()
			}

			if tc.Quote != nil {
				mockWisdomQuotes.On("GetWisdomQuote").Return(tc.Quote.Text).Once()
			}

			srvConn, cliConn := net.Pipe()

			cliExitChan := make(chan struct{})

			go func(tc srvTestCase) {
				defer close(cliExitChan)

				var pc PowChallenge
				err := json.NewDecoder(cliConn).Decode(&pc)
				if tc.GenerateChallengeError != nil {
					require.Error(t, err)
					return
				} else {
					require.NoError(t, err)
				}

				require.Equal(t, PowChallenge(tc.GeneratedChallenge), pc)

				if len(tc.ClientChallengeSolutionRaw) > 0 {
					_, err := io.WriteString(cliConn, tc.ClientChallengeSolutionRaw)
					require.NoError(t, err)
				}

				if tc.ChallengeSolutionCorrect != nil && *tc.ChallengeSolutionCorrect && tc.Quote != nil {
					var wow WordOfWisdom
					err := json.NewDecoder(cliConn).Decode(&wow)
					if tc.HandleConnErrExpected {
						require.Error(t, err)
						return
					} else {
						require.NoError(t, err)
					}
					require.Equal(t, *tc.Quote, wow)
				}
			}(tc)

			err := srv.handleConnection(srvConn)
			if tc.HandleConnErrExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			<-cliExitChan
		})
	}
}

func makeServerWithMocks(t *testing.T) (*Server, *mocks.WisdomQuotesGetter, *mocks.DdosProtector) {
	mockDdosProtector := mocks.NewDdosProtector(t)
	mockQuotesGetter := mocks.NewWisdomQuotesGetter(t)
	return NewServer(
		Config{},
		mockDdosProtector,
		mockQuotesGetter,
		slog.NewTextHandler(io.Discard, new(slog.HandlerOptions)),
	), mockQuotesGetter, mockDdosProtector
}
