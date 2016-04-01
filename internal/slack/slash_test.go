package slack

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/loganjspears/chess"
	"github.com/loganjspears/slackchess/internal/chessutil"
)

type SlashTest struct {
	SlashCommands []*SlashCmd
	Responses     []*Response
}

var (
	validSlashTests = []SlashTest{
		// help
		{
			SlashCommands: []*SlashCmd{newSlashCmd("test", "logan", "help")},
			Responses:     []*Response{helpResponse},
		},
		// invalid command
		{
			SlashCommands: []*SlashCmd{newSlashCmd("test", "logan", "unknown")},
			Responses:     []*Response{unknownResponse},
		},
		// board
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("test", "logan", "play magnus"),
				newSlashCmd("test", "logan", "board"),
			},
			Responses: []*Response{
				boardResponse(testGame("logan", "magnus")),
				boardResponse(testGame("logan", "magnus")),
			},
		},
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("board2", "logan", "board"),
			},
			Responses: []*Response{
				noGameResponse,
			},
		},
		// pgn
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("test", "logan", "play magnus"),
				newSlashCmd("test", "logan", "pgn"),
			},
			Responses: []*Response{
				boardResponse(testGame("logan", "magnus")),
				pgnResponse(testGame("logan", "magnus")),
			},
		},
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("board2", "logan", "board"),
			},
			Responses: []*Response{
				noGameResponse,
			},
		},
		// play and move
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("test", "logan", "play magnus"),
				newSlashCmd("test", "logan", "move e4"),
				newSlashCmd("test", "magnus", "move e5"),
			},
			Responses: []*Response{
				boardResponse(testGame("logan", "magnus")),
				boardResponse(testGame("logan", "magnus", "e4")),
				boardResponse(testGame("logan", "magnus", "e4", "e5")),
			},
		},
		// draw
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("test", "logan", "play magnus"),
				newSlashCmd("test", "logan", "draw offer"),
				newSlashCmd("test", "magnus", "draw accept"),
			},
			Responses: []*Response{
				boardResponse(testGame("logan", "magnus")),
				&Response{
					ResponseType: responseTypeInChannel,
					Text:         "logan offers a draw.",
				},
				boardResponse(testGameConfig("logan", "magnus", func(g *chess.Game) {
					g.Draw(chess.DrawOffer)
				})),
			},
		},
		// resign
		{
			SlashCommands: []*SlashCmd{
				newSlashCmd("test", "logan", "play magnus"),
				newSlashCmd("test", "logan", "move e4"),
				newSlashCmd("test", "magnus", "resign"),
			},
			Responses: []*Response{
				boardResponse(testGame("logan", "magnus")),
				boardResponse(testGame("logan", "magnus", "e4")),
				boardResponse(testGameConfig("logan", "magnus", func(g *chess.Game) {
					g.MoveAlg("e4")
					g.Resign(chess.Black)
				})),
			},
		},
	}
)

func testGameConfig(whitePlayer, blackPlayer string, f func(g *chess.Game)) *chess.Game {
	g := chessutil.NewGame(whitePlayer, blackPlayer)
	f(g)
	return g
}

func testGame(whitePlayer, blackPlayer string, moves ...string) *chess.Game {
	g := chessutil.NewGame(whitePlayer, blackPlayer)
	for _, m := range moves {
		g.MoveAlg(m)
	}
	return g
}

func newSlashCmd(ch, userName, text string) *SlashCmd {
	return &SlashCmd{
		TeamID:    "test",
		ChannelID: ch,
		UserName:  userName,
		Text:      text,
	}
}

func TestValidSlashTests(t *testing.T) {
	for _, test := range validSlashTests {
		fileNames := map[string]bool{}
		for i, s := range test.SlashCommands {
			expected := test.Responses[i]
			actual := s.Response()
			expectedB, _ := json.Marshal(expected)
			actualB, _ := json.Marshal(actual)
			if string(expectedB) != string(actualB) {
				t.Fatalf("expected %s but got %s", string(expectedB), string(actualB))
			}
			fileNames[s.GameFileName()] = true
		}
		for fileName := range fileNames {
			os.Remove(fileName)
		}
	}
}

func TestPlayStockfish(t *testing.T) {
	SetStockfishPath("../stockfish")
	cmd := newSlashCmd("test", "logan", "play slackbot black")
	resp := cmd.Response()
	g, err := cmd.Game()
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Moves()) != 1 {
		b, _ := json.Marshal(resp)
		t.Fatal(string(b))
	}
}
