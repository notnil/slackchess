package imageutil

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/loganjspears/chess"
)

const (
	// md5 hash of test.png
	expectedMD5 = "6f711cc83010cc0694943171ea3c8518"
)

func TestImage(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	g := chess.NewGame()
	if err := g.MoveAlg("e4"); err != nil {
		t.Fatal(err)
	}
	p := g.Position()
	if err := WritePNG(buf, p, chess.E2, chess.E4); err != nil {
		t.Fatal(err)
	}
	actualMD5 := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	if expectedMD5 != actualMD5 {
		t.Fatalf("expected %s md5 hash but got %s", expectedMD5, actualMD5)
	}
}
