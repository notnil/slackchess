package chessutil

import (
	"testing"

	"github.com/loganjspears/chess"
)

func TestNewGame(t *testing.T) {
	g := NewGame("p1", "p2")
	w := tagPairValueForKey(g, "White")
	if w != "p1" {
		t.Fatalf("expected tag pair value %s but got %s", "p1", w)
	}

	b := tagPairValueForKey(g, "Black")
	if b != "p2" {
		t.Fatalf("expected tag pair value %s but got %s", "p2", b)
	}
}

func TestDrawOffer(t *testing.T) {
	g := NewGame("p1", "p2")
	g = AddDrawOffer(g, chess.White)
	actual := DrawOfferColor(g)
	if actual != chess.White {
		t.Fatalf("expected draw offer %s but got %s", chess.White, actual)
	}
	g = RemoveDrawOffer(g)
	actual = DrawOfferColor(g)
	if actual != chess.NoColor {
		t.Fatalf("expected draw offer %s but got %s", chess.NoColor, actual)
	}
}

func TestPlayerToMove(t *testing.T) {
	g := NewGame("p1", "p2")
	actual := PlayerToMove(g)
	if actual != "p1" {
		t.Fatalf("expected player to move %s but got %s", "p1", actual)
	}
	if err := g.MoveAlg("e4"); err != nil {
		t.Fatal(err)
	}
	actual = PlayerToMove(g)
	if actual != "p2" {
		t.Fatalf("expected player to move %s but got %s", "p2", actual)
	}
}

func TestColorOfPlayer(t *testing.T) {
	g := NewGame("p1", "p2")
	actual := ColorOfPlayer(g, "p1")
	if actual != chess.White {
		t.Fatalf("expected player to be %s but got %s", chess.White, actual)
	}
	actual = ColorOfPlayer(g, "p2")
	if actual != chess.Black {
		t.Fatalf("expected player to be %s but got %s", chess.Black, actual)
	}
	actual = ColorOfPlayer(g, "p3")
	if actual != chess.NoColor {
		t.Fatalf("expected player to be %s but got %s", chess.NoColor, actual)
	}
}

func TestPlayerForColor(t *testing.T) {
	g := NewGame("p1", "p2")
	actual := PlayerForColor(g, chess.White)
	if actual != "p1" {
		t.Fatalf("expected player to be %s but got %s", "p1", actual)
	}
	actual = PlayerForColor(g, chess.Black)
	if actual != "p2" {
		t.Fatalf("expected player to be %s but got %s", "p2", actual)
	}
}
