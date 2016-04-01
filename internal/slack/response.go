package slack

import (
	"fmt"

	"github.com/loganjspears/chess"
	"github.com/loganjspears/slackchess/internal/chessutil"
)

// Response is a slash command response.  More information can
// be found here: https://api.slack.com/slash-commands#responding_to_a_command
type Response struct {
	ResponseType string        `json:"response_type,omitempty"`
	Text         string        `json:"text,omitempty"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
}

// Attachment is a slack attachment.  Formatting and other information can be
// found here: https://api.slack.com/docs/attachments
type Attachment struct {
	Fallback   string `json:"fallback,omitempty"`
	Color      string `json:"color,omitempty"`
	Pretext    string `json:"pretext,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title      string `json:"title,omitempty"`
	TitleLink  string `json:"title_link,omitempty"`
	Text       string `json:"text,omitempty"`
	Fields     []struct {
		Title string `json:"title,omitempty"`
		Value string `json:"value,omitempty"`
		Short bool   `json:"short,omitempty"`
	} `json:"fields,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`
}

// boardResponse returns a response for the current state of the board.
func boardResponse(g *chess.Game) *Response {
	whitePlayer := chessutil.PlayerForColor(g, chess.White)
	blackPlayer := chessutil.PlayerForColor(g, chess.Black)
	moveNum := (len(g.Moves()) / 2) + 1
	color := "#f2f2f2"
	if g.Position().Turn() == chess.Black {
		color = "#000000"
	}
	resp := &Response{
		ResponseType: responseTypeInChannel,
		Attachments: []*Attachment{{
			ImageURL: imageURL(g),
			Color:    color,
		}},
	}
	outcome := g.Outcome()
	switch outcome {
	case chess.NoOutcome:
		title := "White to move"
		if g.Position().Turn() == chess.Black {
			title = "Black to move"
		}
		resp.Attachments[0].Text = fmt.Sprintf("@%s vs @%s on move #%d", whitePlayer, blackPlayer, moveNum)
		resp.Attachments[0].Title = title
	case chess.WhiteWon:
		resp.Attachments[0].Title = "White won"
		resp.Attachments[0].Text = fmt.Sprintf("🏆 @%s defeated @%s by %s on move #%d", whitePlayer, blackPlayer, g.Method().String(), moveNum)
		resp.Attachments = append(resp.Attachments, &Attachment{
			Text: g.String(),
		})
	case chess.BlackWon:
		resp.Attachments[0].Title = "Black won"
		resp.Attachments[0].Text = fmt.Sprintf("🏆 @%s defeated @%s by %s on move #%d", blackPlayer, whitePlayer, g.Method().String(), moveNum)
		resp.Attachments = append(resp.Attachments, &Attachment{
			Text: g.String(),
		})
	case chess.Draw:
		resp.Attachments[0].Title = "Draw"
		resp.Attachments[0].Text = fmt.Sprintf("Draw by %s on move #%d", g.Method().String(), moveNum)
		resp.Attachments = append(resp.Attachments, &Attachment{
			Text: g.String(),
		})
	}
	resp.Attachments[0].Fallback = resp.Attachments[0].Title
	return resp
}

// imageURL returns a string in the format:
// http://104.196.27.70/board/rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR.png?markSquares=e2,e4
func imageURL(g *chess.Game) string {
	queryStr := ""
	moves := g.Moves()
	if len(moves) > 0 {
		m := moves[len(moves)-1]
		queryStr = fmt.Sprintf("?markSquares=%s,%s", m.S1().String(), m.S2().String())
	}
	return fmt.Sprintf("%s/board/%s.png%s", baseURL, g.Position().Board().String(), queryStr)
}

var (
	helpResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         helpText,
	}
	unknownResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         unknownText,
	}
	noGameResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         noGameText,
	}
	outOfTurnResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         outOfTurnText,
	}
	notInGameResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         notInGameText,
	}
	invalidMoveResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         invalidMoveText,
	}
	noDrawOfferResponse = &Response{
		ResponseType: responseTypeEphemeral,
		Text:         noDrawToAccept,
	}
)

func errorResponse(err error) *Response {
	text := "The server encountered an error - " + err.Error()
	return &Response{
		ResponseType: responseTypeEphemeral,
		Text:         text,
	}
}

func pgnResponse(g *chess.Game) *Response {
	return &Response{
		ResponseType: responseTypeEphemeral,
		Text:         g.String(),
	}
}

const (
	responseTypeEphemeral = "ephemeral"
	responseTypeInChannel = "in_channel"

	unknownText     = "The command given couldn't be processed.  Use /chess help to view available commands."
	noGameText      = "This command requires an ongoing game in the channel.  Use /chess help to view available commands."
	outOfTurnText   = "Its not your turn!"
	notInGameText   = "Your not in the current game!  Use /chess help to view available commands."
	invalidMoveText = "Invalid move!  Please use Algebraic Notation: https://en.wikipedia.org/wiki/Algebraic_notation_(chess)"
	noDrawToAccept  = "There is no draw offer to accept!"

	helpText = "The chess slash command adds chess playing capabilities to slack.  Here is the list of commands:\n*/chess help* - this help screen\n*/chess play* - '/chess play @magnus' will start a game against the other player in this channel.  There can only be one game per channel and starting a new game will end any in progress.  '/chess play slackbot' and '/chess play slackbot:10' will play against the stockfish engine with an optional skill level.  The skill level starts at 0 and goes up to 20. '/chess play @magnus black' will start a game as black.  '/chess play @magnus random' will start a game as random color.\n*/chess board* - will show the board of the current game\n*/chess pgn* - will show the PGN of the current game\n*/chess move* - 'chess move e4' will move the player using the given Algebraic Notation\n*/chess resign* - resigns the current game\n*/chess draw offer* - offers a draw to other player\n*/chess draw accept* - accepts the draw offer\n*/chess draw reject* - rejects the draw offer (also moving will reject a draw offer)"
)
