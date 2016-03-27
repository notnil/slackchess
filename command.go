package main

import "strings"

type CommandType int

const (
	UnknownCommand CommandType = iota
	Help
	Play
	Board
	PGN
	Move
	Resign
	DrawOffer
	DrawAccept
	DrawReject
)

type Command struct {
	Type CommandType
	Args []string
}

func CommandFromText(text string) Command {
	parts := strings.Split(text, " ")
	if len(parts) == 0 {
		return Command{Type: UnknownCommand, Args: []string{}}
	} else if len(parts) == 1 {
		switch parts[0] {
		case "help":
			return Command{Type: Help, Args: []string{}}
		case "board":
			return Command{Type: Board, Args: []string{}}
		case "pgn":
			return Command{Type: PGN, Args: []string{}}
		case "resign":
			return Command{Type: Resign, Args: []string{}}
		}
	} else if len(parts) == 2 && parts[0] == "draw" {
		switch parts[1] {
		case "offer":
			return Command{Type: DrawOffer, Args: []string{}}
		case "accept":
			return Command{Type: DrawAccept, Args: []string{}}
		case "reject":
			return Command{Type: DrawReject, Args: []string{}}
		}
	} else if len(parts) == 2 && parts[0] == "move" {
		return Command{Type: Move, Args: []string{parts[1]}}
	} else if len(parts) == 2 && parts[0] == "play" {
		return Command{Type: Play, Args: []string{parts[1]}}
	}
	return Command{Type: UnknownCommand, Args: []string{}}
}
