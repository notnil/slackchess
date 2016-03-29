package main

import (
	"errors"
	"fmt"
	"image/color"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/loganjspears/chess"
	"github.com/loganjspears/chessimg"
)

func writeImage(w io.Writer, g *chess.Game, sqs ...chess.Square) error {
	// create temp svg file to be used by rsvg-convert
	fileName := fmt.Sprint(time.Now().UnixNano())
	tempSVG, err := os.Create(fileName + ".svg")
	if err != nil {
		return errors.New("could not create svg file " + err.Error())
	}
	defer tempSVG.Close()
	defer os.Remove(fileName + ".svg")
	mark := chessimg.MarkSquares(color.RGBA{R: 255, G: 255, B: 0, A: 1}, sqs...)
	if err := chessimg.New(tempSVG, mark).EncodeSVG(g.FEN()); err != nil {
		return errors.New("could not write to svg file " + err.Error())
	}

	// create temp png file using rsvg-convert
	// rsvg-convert -h 32 icon.svg > icon-32.png
	if err := exec.Command("rsvg-convert", "-h", "500", fileName+".svg", "-o", fileName+".png").Run(); err != nil {
		return errors.New("could not use rsvg-convert " + err.Error())
	}
	tempPGN, err := os.Open(fileName + ".png")
	if err != nil {
		return errors.New("could not open png file " + err.Error())
	}
	defer tempPGN.Close()
	defer os.Remove(fileName + ".png")

	if _, err := io.Copy(w, tempPGN); err != nil {
		return errors.New("could not copy png to writer " + err.Error())
	}
	return nil
}
