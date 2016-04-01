package chessutil

import (
	"strconv"
	"strings"
	"time"

	"github.com/loganjspears/chess"
)

// NewGame returns a game at the starting position with added tag pairs.
func NewGame(whitePlayer, blackPlayer string) *chess.Game {
	tagPairs := []*chess.TagPair{
		{Key: "Event", Value: "Casual game"},
		{Key: "Date", Value: time.Now().Format("2006.01.02")}, // TODO figure out time zone
		{Key: "White", Value: whitePlayer},
		{Key: "Black", Value: blackPlayer},
		{Key: "TimeControl", Value: "-"},
		{Key: "Annotator", Value: "slackchess"},
	}
	return chess.NewGame(chess.TagPairs(tagPairs))
}

func BotForColor(g *chess.Game, c chess.Color) (isBot bool, skillLvl int) {
	player := PlayerForColor(g, c)
	parts := strings.Split(player, ":")
	isBot = player == "slackbot" || (len(parts) > 0 && parts[0] == "slackbot")
	if isBot && len(parts) == 1 {
		return true, 10
	} else if isBot && len(parts) == 2 {
		i, err := strconv.Atoi(parts[1])
		if err == nil && (i >= 0 && i <= 20) {
			return true, i
		}
	}
	return false, 0
}

// AddDrawOffer adds the "Draw Offer" tag pair showing a pending draw offer.
func AddDrawOffer(g *chess.Game, c chess.Color) *chess.Game {
	return addTag(g, "Draw Offer", c.String())
}

// RemoveDrawOffer removes the "Draw Offer" tag pair showing a pending draw offer.
func RemoveDrawOffer(g *chess.Game) *chess.Game {
	return removeTag(g, "Draw Offer")
}

// DrawOfferColor reads the "Draw Offer" tag pair and returns the color
// that made the offer.  If no offer has been made, chess.NoColor will
// be returned.
func DrawOfferColor(g *chess.Game) chess.Color {
	v := tagPairValueForKey(g, "Draw Offer")
	switch v {
	case chess.White.String():
		return chess.White
	case chess.Black.String():
		return chess.Black
	}
	return chess.NoColor
}

// PlayerToMove returns the name of the current player to move from the
// "White" or "Black" tag pair.
func PlayerToMove(g *chess.Game) string {
	if g.Position().Turn() == chess.White {
		return tagPairValueForKey(g, "White")
	}
	return tagPairValueForKey(g, "Black")
}

// ColorOfPlayer returns the chess.Color of the given player.  If the
// player isn't in the "White" or "Black" tag pair then chess.NoColor
// is returned.
func ColorOfPlayer(g *chess.Game, player string) chess.Color {
	if tagPairValueForKey(g, "White") == player {
		return chess.White
	} else if tagPairValueForKey(g, "Black") == player {
		return chess.Black
	}
	return chess.NoColor
}

// PlayerForColor returns the player for the give chess.Color.
func PlayerForColor(g *chess.Game, c chess.Color) string {
	switch c {
	case chess.White:
		return tagPairValueForKey(g, "White")
	case chess.Black:
		return tagPairValueForKey(g, "Black")
	}
	return ""
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
