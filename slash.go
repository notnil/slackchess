package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/loganjspears/chess"
)

// SlashCmd is a slash command.  Here is example data:
// token=gIkuvaNzQIHg97ATvDxqgjtO
// team_id=T0001
// team_domain=example
// channel_id=C2147483705
// channel_name=test
// user_id=U2147483697
// user_name=Steve
// command=/weather
// text=94070
// response_url=https://hooks.slack.com/commands/1234/5678
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

func (s *SlashCmd) Response() *Response {
	cmd := CommandFromText(s.Text)
	switch cmd.Type {
	case UnknownCommand:
		return unknownResponse
	case Help:
		return helpResponse
	case Play:
		g := newGame(s.UserName, cmd.Args[0])
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	}
	g, _ := s.Game()
	if g == nil {
		return noGameResponse
	}
	c := colorOfPlayer(g, s.UserName)
	if c == chess.NoColor {
		return notInGameResponse
	}
	switch cmd.Type {
	case Move:
		player := playerToMove(g)
		if player != s.UserName {
			return outOfTurnResponse
		}
		if err := g.MoveAlg(cmd.Args[0]); err != nil {
			return invalidMoveResponse
		}
		g = removeTag(g, "Draw Offer")
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case Board:
		return boardResponse(g)
	case PGN:
		return pgnResponse(g)
	case Resign:
		g.Resign(c)
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case DrawOffer:
		g = addDrawOffer(g, c)
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return &Response{
			ResponseType: responseTypeInChannel,
			Text:         s.UserName + " offers a draw.",
		}
	case DrawAccept:
		drawColor := drawOfferColor(g)
		if drawColor != c.Other() {
			return noDrawOfferResponse
		}
		g = removeTag(g, "Draw Offer")
		if err := g.Draw(chess.DrawOffer); err != nil {
			return errorResponse(err)
		}
		if err := s.SaveGame(g); err != nil {
			return errorResponse(err)
		}
		return boardResponse(g)
	case DrawReject:
		g = removeTag(g, "Draw Offer")
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

func (s *SlashCmd) SaveGame(game *chess.Game) error {
	pgn := []byte(game.String())
	return ioutil.WriteFile(s.GameFileName(), pgn, 0666)
}

func (s *SlashCmd) GameFileName() string {
	return fmt.Sprintf("%s_%s.pgn", s.TeamID, s.ChannelID)
}
