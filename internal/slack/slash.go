package slack

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/loganjspears/chess"
	"github.com/loganjspears/slackchess/internal/chessutil"
	"github.com/loganjspears/slackchess/internal/stockfish"
)

var (
	baseURL       = ""
	stockfishPath = ""
)

// SetBaseURL sets the baseURL used in the board image URL embedded in slack attachments.
func SetBaseURL(url string) {
	baseURL = url
}

// SetStockfishPath sets the path to the stockfish folder
func SetStockfishPath(path string) {
	stockfishPath = path
}

// SlashCmd represents a slack "Slash Command".  You can read more about
// slash commands here: https://api.slack.com/slash-commands
type SlashCmd struct {
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	TeamDomain  string `schema:"team_domain"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

// Response returns the response to the given command.
func (s *SlashCmd) Response() *Response {
	cmd := userEntryFromText(s.Text)
	switch cmd.Type {
	case cmdUnknown:
		return unknownResponse
	case cmdHelp:
		return helpResponse
	case cmdPlay:
		var g *chess.Game
		switch cmd.Args[1] {
		case "", "white":
			g = chessutil.NewGame(s.UserName, cmd.Args[0])
		case "black":
			g = chessutil.NewGame(cmd.Args[0], s.UserName)
		case "random":
			rand.Seed(time.Now().UnixNano())
			r := rand.Intn(2)
			if r == 0 {
				g = chessutil.NewGame(s.UserName, cmd.Args[0])
			} else {
				g = chessutil.NewGame(cmd.Args[0], s.UserName)
			}
		}
		if err := checkForStockfishMove(g); err != nil {
			return errorResponse(err)
		}
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	}
	g, _ := s.Game()
	if g == nil {
		return noGameResponse
	}
	c := chessutil.ColorOfPlayer(g, s.UserName)
	if c == chess.NoColor {
		return notInGameResponse
	}
	switch cmd.Type {
	case cmdMove:
		player := chessutil.PlayerToMove(g)
		if player != s.UserName {
			return outOfTurnResponse
		}
		if err := g.MoveAlg(cmd.Args[0]); err != nil {
			return invalidMoveResponse
		}
		g = chessutil.RemoveDrawOffer(g)
		if err := checkForStockfishMove(g); err != nil {
			return errorResponse(err)
		}
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case cmdBoard:
		return boardResponse(g)
	case cmdPGN:
		return pgnResponse(g)
	case cmdResign:
		g.Resign(c)
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case cmdDrawOffer:
		g = chessutil.AddDrawOffer(g, c)
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return &Response{
			ResponseType: responseTypeInChannel,
			Text:         s.UserName + " offers a draw.",
		}
	case cmdDrawAccept:
		drawColor := chessutil.DrawOfferColor(g)
		if drawColor != c.Other() {
			return noDrawOfferResponse
		}
		g = chessutil.RemoveDrawOffer(g)
		if err := g.Draw(chess.DrawOffer); err != nil {
			return errorResponse(err)
		}
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case cmdDrawReject:
		g = chessutil.RemoveDrawOffer(g)
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return &Response{
			ResponseType: responseTypeInChannel,
			Text:         s.UserName + " rejects a draw.",
		}
	}
	return unknownResponse
}

// Game returns the game associated with the slash command.
// Games are stored on the file system and an error is returned
// if there is an IO error or if the game doesn't exist.
func (s *SlashCmd) Game() (*chess.Game, error) {
	chess.NewGame()
	// open pgn file if it exists
	f, err := os.Open(s.GameFileName())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	pgn, err := chess.PGN(f)
	if err != nil {
		return nil, err
	}
	return chess.NewGame(pgn), nil
}

// SaveGame saves the game to disk.  An error is returned if there is an IO error.
func (s *SlashCmd) SaveGame(game *chess.Game) error {
	pgn := []byte(game.String())
	return ioutil.WriteFile(s.GameFileName(), pgn, 0666)
}

// GameFileName returns the file name of where the game would be stored (even if there
// is no game actually saved).
func (s *SlashCmd) GameFileName() string {
	return fmt.Sprintf("%s_%s.pgn", s.TeamID, s.ChannelID)
}

// commandType represents a command supported by the slash command
type commandType int

const (
	cmdUnknown commandType = iota
	cmdHelp
	cmdPlay
	cmdBoard
	cmdPGN
	cmdMove
	cmdResign
	cmdDrawOffer
	cmdDrawAccept
	cmdDrawReject
)

// userEntry is a structure result of SlashCmd's text field
type userEntry struct {
	Type commandType
	Args []string
}

// userEntryFromText takes a text argument and returns a formatted
// userEntry.  If the text can't be parsed then returned userEntry's
// type will be cmdUnknown.
func userEntryFromText(text string) userEntry {
	parts := strings.Split(text, " ")
	if len(parts) == 1 {
		switch parts[0] {
		case "help":
			return userEntry{Type: cmdHelp, Args: []string{}}
		case "board":
			return userEntry{Type: cmdBoard, Args: []string{}}
		case "pgn":
			return userEntry{Type: cmdPGN, Args: []string{}}
		case "resign":
			return userEntry{Type: cmdResign, Args: []string{}}
		}
	} else if len(parts) == 2 && parts[0] == "draw" {
		switch parts[1] {
		case "offer":
			return userEntry{Type: cmdDrawOffer, Args: []string{}}
		case "accept":
			return userEntry{Type: cmdDrawAccept, Args: []string{}}
		case "reject":
			return userEntry{Type: cmdDrawReject, Args: []string{}}
		}
	} else if len(parts) == 2 && parts[0] == "move" {
		return userEntry{Type: cmdMove, Args: []string{parts[1]}}
	} else if (len(parts) == 2 || len(parts) == 3) && parts[0] == "play" {
		username := strings.Replace(parts[1], "@", "", -1)
		colorTxt := ""
		if len(parts) == 3 {
			colorTxt = parts[2]
		}
		return userEntry{Type: cmdPlay, Args: []string{username, colorTxt}}
	}
	return userEntry{Type: cmdUnknown, Args: []string{}}
}

func checkForStockfishMove(g *chess.Game) error {
	c := g.Position().Turn()
	isBot, skillLvl := chessutil.BotForColor(g, c)
	if !isBot {
		return nil
	}
	move, err := stockfish.Move(g, skillLvl, stockfishPath)
	if err != nil {
		return err
	}
	return g.Move(move)
}
