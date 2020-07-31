package service_test

import (
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/flaambe/avito/internal/model"
	"github.com/flaambe/avito/internal/repository/mocks"
	"github.com/flaambe/avito/internal/service"
	"github.com/flaambe/avito/internal/view"

	"github.com/stretchr/testify/assert"
)

var (
	userModel    model.User
	chatModel    model.Chat
	messageModel model.Message

	userRepoMock    *mocks.UserRepository
	chatRepoMock    *mocks.ChatRepository
	messageRepoMock *mocks.MessageRepository
)

func TestMain(m *testing.M) {
	userModel = model.User{
		ID:       primitive.NewObjectID(),
		UserName: "Test",
	}
	chatModel = model.Chat{
		ID:    primitive.NewObjectID(),
		Name:  "test_chat",
		Users: []model.User{userModel},
	}
	messageModel = model.Message{
		ID:        primitive.NewObjectID(),
		Chat:      chatModel.ID,
		Author:    userModel.ID,
		Text:      "Test_text",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	userRepoMock = new(mocks.UserRepository)
	userRepoMock.On("FindUserByID", userModel.ID.Hex()).Return(userModel, nil)
	userRepoMock.On("InsertUser", userModel.UserName).Return(userModel.ID.Hex(), nil)

	chatRepoMock = new(mocks.ChatRepository)
	chatRepoMock.On("FindChatByID", chatModel.ID.Hex()).Return(chatModel, nil)
	chatRepoMock.On("FindChats", userModel).Return([]model.Chat{chatModel}, nil)
	chatRepoMock.On("InsertChat", chatModel.Name, chatModel.Users).Return(chatModel.ID.Hex(), nil)

	messageRepoMock = new(mocks.MessageRepository)
	messageRepoMock.On("FindMessages", chatModel).Return([]model.Message{messageModel}, nil)
	messageRepoMock.On("InsertMessage", chatModel, userModel, messageModel.Text).Return(messageModel.ID.Hex(), nil)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestAddUser(t *testing.T) {
	assert := assert.New(t)
	testObj := service.NewChatService(userRepoMock, chatRepoMock, messageRepoMock)

	userRequest := view.NewUserRequest{
		UserName: userModel.UserName,
	}

	userResponse, err := testObj.AddUser(userRequest)
	assert.NoError(err)
	assert.Equal(userModel.ID.Hex(), userResponse.ID)
}

func TestAddChat(t *testing.T) {
	assert := assert.New(t)
	testObj := service.NewChatService(userRepoMock, chatRepoMock, messageRepoMock)

	chatRequest := view.NewChatRequest{
		Name:    chatModel.Name,
		UsersID: []string{userModel.ID.Hex()},
	}

	chatResponse, err := testObj.AddChat(chatRequest)
	assert.NoError(err)
	assert.Equal(chatModel.ID.Hex(), chatResponse.ID)
}

func TestAddMessage(t *testing.T) {
	assert := assert.New(t)
	testObj := service.NewChatService(userRepoMock, chatRepoMock, messageRepoMock)

	messageRequest := view.NewMessageRequest{
		ChatID: chatModel.ID.Hex(),
		UserID: userModel.ID.Hex(),
		Text:   messageModel.Text,
	}

	messageResponse, err := testObj.AddMessage(messageRequest)
	assert.NoError(err)
	assert.Equal(messageModel.ID.Hex(), messageResponse.ID)
}

func TestGetChats(t *testing.T) {
	assert := assert.New(t)
	testObj := service.NewChatService(userRepoMock, chatRepoMock, messageRepoMock)

	chatsRequest := view.ChatsRequest{
		UserID: userModel.ID.Hex(),
	}

	chatsResponse, err := testObj.GetChats(chatsRequest)
	assert.NoError(err)

	var usersView []view.User

	for _, user := range chatsResponse[0].Users {
		usersView = append(usersView, view.User{ID: user.ID, UserName: user.UserName})
	}

	assert.Equal(chatModel.ID.Hex(), chatsResponse[0].ID)
	assert.Equal(chatModel.Name, chatsResponse[0].Name)
	assert.Equal(usersView, chatsResponse[0].Users)
	assert.Equal(chatModel.CreatedAt.Time().String(), chatsResponse[0].CreatedAt)
}

func TestGetMessages(t *testing.T) {
	assert := assert.New(t)
	testObj := service.NewChatService(userRepoMock, chatRepoMock, messageRepoMock)

	messagesRequest := view.MessagesRequest{
		СhatID: chatModel.ID.Hex(),
	}

	messagesResponse, err := testObj.GetMessages(messagesRequest)
	assert.NoError(err)

	assert.Equal(messageModel.ID.Hex(), messagesResponse[0].ID)
	assert.Equal(messageModel.Chat.Hex(), messagesResponse[0].ChatID)
	assert.Equal(messageModel.Author.Hex(), messagesResponse[0].AuthorID)
	assert.Equal(messageModel.Text, messagesResponse[0].Text)
	assert.Equal(messageModel.CreatedAt.Time().String(), messagesResponse[0].CreatedAt)
}
