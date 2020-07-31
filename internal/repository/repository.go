package repository

import (
	"github.com/flaambe/avito/internal/model"
)

type UserRepository interface {
	FindUserByID(id string) (model.User, error)
	InsertUser(name string) (string, error)
}

type ChatRepository interface {
	FindChatByID(id string) (model.Chat, error)
	FindChats(user model.User) ([]model.Chat, error)
	InsertChat(name string, users []model.User) (string, error)
}

type MessageRepository interface {
	FindMessages(chat model.Chat) ([]model.Message, error)
	InsertMessage(chat model.Chat, user model.User, text string) (string, error)
}
