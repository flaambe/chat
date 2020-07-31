package view

type NewMessageRequest struct {
	ChatID string `json:"chat"`
	UserID string `json:"author"`
	Text   string `json:"text"`
}

type NewMessageResponse struct {
	ID string `json:"id"`
}

type Message struct {
	ID        string `json:"id"`
	ChatID    string `json:"chat"`
	AuthorID  string `json:"author"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type MessagesResponse []Message
