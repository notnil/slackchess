# chessimg
[![GoDoc](https://godoc.org/github.com/loganjspears/chessimg?status.svg)](https://godoc.org/github.com/loganjspears/chessimg)
[![Build Status](https://drone.io/github.com/loganjspears/chessimg/status.png)](https://drone.io/github.com/loganjspears/chessimg/latest)
[![Coverage Status](https://coveralls.io/repos/github/loganjspears/chessimg/badge.svg?branch=master)](https://coveralls.io/github/loganjspears/chessimg?branch=master)
[![Go Report Card](http://goreportcard.com/badge/loganjspears/chessimg)](http://goreportcard.com/report/loganjspears/chessimg)

### Code Example

```go
// create file
f, err := os.Create("example.svg")
if err != nil {
    log.Fatal(err)
}
// write image of position and marked squares to file
fenStr := "rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1"
mark := chessimg.MarkSquares(color.RGBA{255, 255, 0, 1}, chess.D2, chess.D4)
if err := chessimg.New(f, mark).EncodeSVG(fenStr); err != nil {
	log.Fatal(err)
}
```

### Resulting Image

![rnbqkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBQKBNR b KQkq - 0 1](/example.png)
 
