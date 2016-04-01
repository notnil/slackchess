package stockfish

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/loganjspears/chess"
)

// Move returns a move from the Stockfish chess engine.  skillLvl is the skill
// level of the engine and must be [0,20].  execPath should be the path to stockfish
// directory.  An error is returned if there an issue communicating with the stockfish executable.
func Move(g *chess.Game, skillLvl int, execPath string) (*chess.Move, error) {
	if skillLvl < 0 || skillLvl > 20 {
		return nil, errors.New("stockfish: skill level must be between 0 and 20")
	}
	buf := bytes.NewBuffer([]byte{})
	cmd := exec.Command(execPath+"/stockfish.sh", fmt.Sprint(skillLvl), g.FEN(), executable(execPath))
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, errors.New("stockfish: error occured when running stockfish command " + err.Error())
	}

	output := buf.String()
	moveTxt := parseOutput(output)
	if moveTxt == "" {
		return nil, errors.New("stockfish: couldn't parse stockfish output - " + output)
	}

	move := getMoveFromText(g, moveTxt)
	if move == nil {
		return nil, errors.New("stockfish: couldn't parse stockfish move - " + moveTxt)
	}
	return move, nil
}

func getMoveFromText(g *chess.Game, moveTxt string) *chess.Move {
	moveTxt = strings.Replace(moveTxt, "x", "", -1)
	isValidLength := (len(moveTxt) == 4 || len(moveTxt) == 5)
	if !isValidLength {
		return nil
	}
	s1Txt := moveTxt[0:2]
	s2Txt := moveTxt[2:4]
	promoTxt := ""
	if len(moveTxt) == 5 {
		promoTxt = moveTxt[4:5]
	}
	for _, m := range g.ValidMoves() {
		if m.S1().String() == s1Txt &&
			m.S2().String() == s2Txt &&
			promoTxt == m.Promo().String() {
			return m
		}
	}
	return nil
}

// searching for format: "bestmove e2e4 ponder e7e6"
func parseOutput(output string) string {
	output = strings.Replace(output, "\n", " ", -1)
	words := strings.Split(output, " ")
	next := false
	for _, word := range words {
		if next {
			return word
		}
		if word == "bestmove" {
			next = true
		}
	}
	return ""
}

func executable(path string) string {
	if runtime.GOOS == "darwin" {
		return path + "/stockfish-7-64-mac"
	}
	return path + "/stockfish-7-64-linux"
}
