package view

type NewUserRequest struct {
	UserName string `json:"username"`
}

type UserRequest struct {
	UserID string `json:"user"`
}

type NewUserResponse struct {
	ID string `json:"id"`
}
