package stockfish

import (
	"testing"

	"github.com/loganjspears/chess"
)

func TestStockfish(t *testing.T) {
	if _, err := Move(chess.NewGame(), 5, "."); err != nil {
		t.Fatal(err)
	}
}
