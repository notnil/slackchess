# chessimg
[![GoDoc](https://godoc.org/github.com/loganjspears/chessimg?status.svg)](https://godoc.org/github.com/loganjspears/chessimg)
[![Build Status](https://drone.io/github.com/loganjspears/chessimg/status.png)](https://drone.io/github.com/loganjspears/chessimg/latest)
[![Coverage Status](https://coveralls.io/repos/github/loganjspears/chessimg/badge.svg?branch=master)](https://coveralls.io/github/loganjspears/chessimg?branch=master)
[![Go Report Card](http://goreportcard.com/badge/loganjspears/chessimg)](http://goreportcard.com/report/loganjspears/chessimg)

## Usage

```go
// populate buffer w/ SVG of the starting position
f, err := os.Create("actual.svg")
if err != nil {
    log.Fatal(err)
}
fenStr := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
if err := chessimg.New(f).EncodeSVG(fenStr); err != nil {
	log.Fatal(err)
}
// take a look at actual.svg in the repo to view the result
```

 
