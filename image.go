package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/loganjspears/chess"
	"github.com/loganjspears/chessimg"
)

func writeImage(g *chess.Game, w io.Writer) error {
	// create temp svg file to be used by rsvg-convert
	fileName := fmt.Sprint(time.Now().UnixNano())
	tempSVG, err := os.Create(fileName + ".svg")
	if err != nil {
		return errors.New("could not create svg file " + err.Error())
	}
	defer tempSVG.Close()
	defer os.Remove(fileName + ".svg")
	if err := chessimg.New(tempSVG).EncodeSVG(g.FEN()); err != nil {
		return errors.New("could not write to svg file " + err.Error())
	}

	// create temp png file using rsvg-convert
	// rsvg-convert -h 32 icon.svg > icon-32.png
	if err := exec.Command("rsvg-convert", "-h", "300", fileName+".svg", "-o", fileName+".png").Run(); err != nil {
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

// func boardImgHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "", http.StatusNotFound)
// 		return
// 	}
// 	log.Println("board img handler - %s", r.URL.Path)
// 	path := strings.TrimPrefix(r.URL.Path, "/board/")
// 	path = strings.TrimSuffix(path, ".png")
// 	path = path + " w KQkq - 0 1"

// 	fen, err := chess.FEN(strings.NewReader(path))
// 	if err != nil {
// 		http.Error(w, "could not parse fen "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	g := chess.NewGame(fen)
// 	log.Println("game created")

// 	fileName := fmt.Sprint(time.Now().UnixNano())
// 	tempSVG, err := os.Create(fileName + ".svg")
// 	if err != nil {
// 		http.Error(w, "could not create svg file "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("temp svg created")

// 	defer tempSVG.Close()
// 	defer os.Remove(fileName + ".svg")

// 	if err := chessimg.New(tempSVG).EncodeSVG(g.State().Board()); err != nil {
// 		http.Error(w, "could not write to svg file"+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("temp svg written")

// 	// rsvg-convert -h 32 icon.svg > icon-32.png
// 	if err := exec.Command("rsvg-convert", "-h", "300", fileName+".svg", "-o", fileName+".png").Run(); err != nil {
// 		http.Error(w, "could not use rsvg-convert "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("temp png created")
// 	tempPGN, err := os.Open(fileName + ".png")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("temp png opened")
// 	defer tempPGN.Close()
// 	defer os.Remove(fileName + ".png")

// 	w.Header().Set("Content-Type", "image/png")
// 	w.Header().Set("Cache-Control", "max-age=31536000")
// 	if _, err := io.Copy(w, tempPGN); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Println("copied png to output")
// }
