package main

import (
	"os"
	"testing"

	"github.com/loganjspears/chess"
)

func TestImage(t *testing.T) {
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := writeImage(chess.NewGame(), f); err != nil {
		t.Fatal(err)
	}
}
