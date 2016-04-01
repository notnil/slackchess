package imageutil

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"runtime"
	"testing"

	"github.com/loganjspears/chess"
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
	actual := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	expected := getExpectedMD5Hash()
	if expected != actual {
		t.Fatalf("expected %s md5 hash but got %s", expected, actual)
	}
}

func getExpectedMD5Hash() string {
	if runtime.GOOS == "darwin" {
		return "6f711cc83010cc0694943171ea3c8518"
	}
	return "954a16709a5ee00d5f34bec0d01a3fcc"
}
