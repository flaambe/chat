package view

type User struct {
	ID       string `json:"id"`
	UserName string `json:"name"`
}

type Chat struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Users     []User `json:"users"`
	CreatedAt string `json:"created_at"`
}

type NewChatRequest struct {
	Name    string   `json:"name"`
	UsersID []string `json:"users"`
}

type ChatsRequest struct {
	UserID string `json:"user"`
}

type ChatRequest struct {
	СhatID string `json:"chat"`
}

type NewChatResponse struct {
	ID string `json:"id"`
}

type ChatsResponse []Chat
