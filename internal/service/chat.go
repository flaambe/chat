package service

import (
	"github.com/flaambe/avito/internal/errs"
	"github.com/flaambe/avito/internal/model"
	"github.com/flaambe/avito/internal/view"
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

type ChatService struct {
	userRepo    UserRepository
	chatRepo    ChatRepository
	messageRepo MessageRepository
}

func NewChatService(u UserRepository, c ChatRepository, m MessageRepository) *ChatService {
	return &ChatService{u, c, m}
}

func (c *ChatService) AddUser(user view.NewUserRequest) (view.NewUserResponse, error) {
	userId, err := c.userRepo.InsertUser(user.UserName)
	if err != nil {
		return view.NewUserResponse{}, errs.New(500, "internal server error", err)
	}

	return view.NewUserResponse{ID: userId}, nil
}

func (c *ChatService) AddChat(chat view.NewChatRequest) (view.NewChatResponse, error) {
	var usersModel []model.User

	for _, userID := range chat.UsersID {
		user, err := c.userRepo.FindUserByID(userID)
		if err != nil {
			return view.NewChatResponse{}, errs.New(404, "user not found", err)
		}

		userModel := model.User{
			ID:       user.ID,
			UserName: user.UserName,
		}
		usersModel = append(usersModel, userModel)
	}

	chatId, err := c.chatRepo.InsertChat(chat.Name, usersModel)
	if err != nil {
		return view.NewChatResponse{}, errs.New(500, "internal server error", err)
	}

	return view.NewChatResponse{ID: chatId}, nil
}

func (c *ChatService) AddMessage(message view.NewMessageRequest) (view.NewMessageResponse, error) {
	chat, err := c.chatRepo.FindChatByID(message.ChatID)
	if err != nil {
		return view.NewMessageResponse{}, errs.New(404, "chat not found", err)
	}

	user, err := c.userRepo.FindUserByID(message.UserID)
	if err != nil {
		return view.NewMessageResponse{}, errs.New(404, "user not found", err)
	}

	messageId, err := c.messageRepo.InsertMessage(chat, user, message.Text)
	if err != nil {
		return view.NewMessageResponse{}, errs.New(500, "internal server error", err)
	}

	return view.NewMessageResponse{ID: messageId}, nil
}

func (c *ChatService) GetChats(chats view.ChatsRequest) (view.ChatsResponse, error) {
	var chatsView []view.Chat

	user, err := c.userRepo.FindUserByID(chats.UserID)
	if err != nil {
		return view.ChatsResponse{}, errs.New(404, "user not found", err)
	}

	chatsModel, err := c.chatRepo.FindChats(user)
	if err != nil {
		return view.ChatsResponse{}, errs.New(404, "chats not found", err)
	}

	for _, chatModel := range chatsModel {
		var users []view.User
		for _, user := range chatModel.Users {
			userView := view.User{
				ID:       user.ID.Hex(),
				UserName: user.UserName,
			}

			users = append(users, userView)
		}

		chatView := view.Chat{
			ID:        chatModel.ID.Hex(),
			Name:      chatModel.Name,
			Users:     users,
			CreatedAt: chatModel.CreatedAt.Time().String(),
		}

		chatsView = append(chatsView, chatView)
	}

	return chatsView, nil
}

func (c *ChatService) GetMessages(chat view.MessagesRequest) (view.MessagesResponse, error) {
	var messagesView []view.Message

	chatModel, err := c.chatRepo.FindChatByID(chat.СhatID)
	if err != nil {
		return view.MessagesResponse{}, errs.New(404, "chat not found", err)
	}

	messagesModel, err := c.messageRepo.FindMessages(chatModel)
	if err != nil {
		return view.MessagesResponse{}, errs.New(404, "messages not found", err)
	}

	for _, messageModel := range messagesModel {
		messageView := view.Message{
			ID:        messageModel.ID.Hex(),
			ChatID:    messageModel.Chat.Hex(),
			AuthorID:  messageModel.Author.Hex(),
			Text:      messageModel.Text,
			CreatedAt: messageModel.CreatedAt.Time().String(),
		}

		messagesView = append(messagesView, messageView)
	}

	return messagesView, nil
}
