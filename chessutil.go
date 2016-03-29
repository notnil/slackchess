package main

import (
	"strings"
	"time"

	"github.com/loganjspears/chess"
)

func newGame(whitePlayer, blackPlayer string) *chess.Game {
	tagPairs := []*chess.TagPair{
		{Key: "Site", Value: "Slack /chess"},
		{Key: "Date", Value: time.Now().Format("2006.01.02")}, // TODO figure out time zone
		{Key: "White", Value: whitePlayer},
		{Key: "Black", Value: blackPlayer},
	}
	return chess.NewGame(chess.TagPairs(tagPairs))
}

func addDrawOffer(g *chess.Game, c chess.Color) *chess.Game {
	return addTag(g, "Draw Offer", c.String())
}

func drawOfferColor(g *chess.Game) chess.Color {
	v := tagPairValueForKey(g, "Draw Offer")
	switch v {
	case chess.White.String():
		return chess.White
	case chess.Black.String():
		return chess.Black
	}
	return chess.NoColor
}

func removeDrawOffer(g *chess.Game) *chess.Game {
	return removeTag(g, "Draw Offer")
}

func playerToMove(g *chess.Game) string {
	if g.Position().Turn() == chess.White {
		return tagPairValueForKey(g, "White")
	}
	return tagPairValueForKey(g, "Black")
}

func colorOfPlayer(g *chess.Game, player string) chess.Color {
	if tagPairValueForKey(g, "White") == player {
		return chess.White
	} else if tagPairValueForKey(g, "Black") == player {
		return chess.Black
	}
	return chess.NoColor
}

func addTag(g *chess.Game, key, value string) *chess.Game {
	tagPairs := g.TagPairs()
	tagPairs = append(tagPairs, &chess.TagPair{Key: key, Value: value})
	pgnStr := g.String()
	pgn, _ := chess.PGN(strings.NewReader(pgnStr))
	return chess.NewGame(pgn, chess.TagPairs(tagPairs))
}

func removeTag(g *chess.Game, key string) *chess.Game {
	tagPairs := g.TagPairs()
	cp := []*chess.TagPair{}
	for _, tagPair := range tagPairs {
		if tagPair.Key != key {
			cp = append(cp, tagPair)
		}
	}
	pgnStr := g.String()
	pgn, _ := chess.PGN(strings.NewReader(pgnStr))
	return chess.NewGame(pgn, chess.TagPairs(cp))
}

func tagPairValueForKey(g *chess.Game, key string) string {
	for _, tagPair := range g.TagPairs() {
		if tagPair.Key == key {
			return tagPair.Value
		}
	}
	return ""
}

func lastMoveText(g *chess.Game) string {
	pgn := g.String()
	parts := strings.Split(pgn, " ")
	if len(parts) <= 2 {
		return ""
	}
	return parts[len(parts)-3]
}
