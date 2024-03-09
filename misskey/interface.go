package misskey

type Note struct {
	Id        string
	CreatedAt string
	Text      string
	Cw        string
	// User
	UserId     string
	Visibility string
}

type NotesCreateRequestBody struct {
	I              string   `json:"i"`
	Text           string   `json:"text"`
	Visibility     string   `json:"visibility,omitempty"` // public, home, followers, specified
	VisibleUserIds []string `json:"visibleUserIds,omitempty"`
	Cw             string   `json:"cw,omitempty"`
	LocalOnly      bool     `json:"localOnly,omitempty"`
}

type NotesCreateResponse struct {
	CreatedNote Note
}
