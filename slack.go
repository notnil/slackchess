package main

type SlashCmd struct {
	Token       string `schema:"token"`
	TeamID      string `schema:"team_id"`
	TeamDomain  string `schema:"team_domain"`
	ChannelID   string `schema:"channel_id"`
	ChannelName string `schema:"channel_name"`
	UserID      string `schema:"user_id"`
	UserName    string `schema:"user_name"`
	Command     string `schema:"command"`
	Text        string `schema:"text"`
	ResponseURL string `schema:"response_url"`
}

type Attachment struct {
	ResponseType string `json:"response_type,omitempty"`
	Text         string `json:"text,omitempty"`
	Attachments  []struct {
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
	} `json:"attachments"`
}

const (
	responseTypeEphemeral = "ephemeral"
	responseTypeInChannel = "in_channel"

	helpResponseJSON = `{
    "text": "The chess slash command adds chess playing capabilities to slack.  Here is the list of commands:\n*/chess help* - this help screen\n*/chess play* - 'chess play @magnus' will start a game against the other player in this channel.  There can only be one game per channel and starting a new game will end any in progress.\n*/chess board* - will show the board of the current game\n*/chess move* - 'chess move e4' will move the player using the given Algebraic Notation\n*/chess resign* - resigns the current game\n*/chess draw offer* - offers a draw to other player\n*/chess draw accept* - accepts the draw offer\n*/chess draw reject* - rejects the draw offer (also moving will reject a draw offer)"
    }`
)
