package service_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/flaambe/avito/internal/errs"
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

	userErrRepoMock := new(mocks.UserRepository)
	userErrRepoMock.On("InsertUser", userModel.UserName).Return("-1", errors.New("internal db error"))

	testObj = service.NewChatService(userErrRepoMock, chatRepoMock, messageRepoMock)
	userResponse, err = testObj.AddUser(userRequest)
	assert.Error(err)
	var responseError *errs.ResponseError
	if errors.As(err, &responseError) {
		assert.Equal(500, responseError.Status)
	}
	assert.Empty(userResponse.ID)
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

	chatErrRequest := view.NewChatRequest{
		Name:    chatModel.Name,
		UsersID: []string{"incorrect id"},
	}
	userErrRepoMock := new(mocks.UserRepository)
	userErrRepoMock.On("FindUserByID", "incorrect id").Return(model.User{}, errors.New("incorrect id"))
	testObj = service.NewChatService(userErrRepoMock, chatRepoMock, messageRepoMock)
	chatResponse, err = testObj.AddChat(chatErrRequest)
	assert.Error(err)
	var responseError *errs.ResponseError
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Equal("", chatResponse.ID)

	chatErrRepoMock := new(mocks.ChatRepository)
	chatErrRepoMock.On("InsertChat", chatModel.Name, chatModel.Users).Return("", errors.New("internal db error"))
	testObj = service.NewChatService(userRepoMock, chatErrRepoMock, messageRepoMock)
	chatResponse, err = testObj.AddChat(chatRequest)
	assert.Error(err)
	if errors.As(err, &responseError) {
		assert.Equal(500, responseError.Status)
	}
	assert.Empty(chatResponse.ID)
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

	newMessageErrRequest := view.NewMessageRequest{
		ChatID: "incorrect id",
	}
	chatErrRepoMock := new(mocks.ChatRepository)
	chatErrRepoMock.On("FindChatByID", "incorrect id").Return(model.Chat{}, errors.New("incorrect id"))
	testObj = service.NewChatService(userRepoMock, chatErrRepoMock, messageRepoMock)
	messageResponse, err = testObj.AddMessage(newMessageErrRequest)
	assert.Error(err)
	var responseError *errs.ResponseError
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(messageResponse)

	newMessageErrRequest = view.NewMessageRequest{
		ChatID: chatModel.ID.Hex(),
		UserID: "incorrect id",
	}
	userErrRepoMock := new(mocks.UserRepository)
	userErrRepoMock.On("FindUserByID", "incorrect id").Return(model.User{}, errors.New("incorrect id"))
	testObj = service.NewChatService(userErrRepoMock, chatRepoMock, messageRepoMock)
	messageResponse, err = testObj.AddMessage(newMessageErrRequest)
	assert.Error(err)
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(messageResponse)

	newMessageErrRequest = view.NewMessageRequest{
		ChatID: chatModel.ID.Hex(),
		UserID: userModel.ID.Hex(),
		Text:   messageModel.Text,
	}
	messageErrRepoMock := new(mocks.MessageRepository)
	messageErrRepoMock.On("InsertMessage", chatModel, userModel, messageModel.Text).Return("", errors.New("internal db error"))
	testObj = service.NewChatService(userRepoMock, chatRepoMock, messageErrRepoMock)
	messageResponse, err = testObj.AddMessage(newMessageErrRequest)
	assert.Error(err)
	if errors.As(err, &responseError) {
		assert.Equal(500, responseError.Status)
	}
	assert.Empty(messageResponse)
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

	chatsErrRequest := view.ChatsRequest{
		UserID: "incorrect id",
	}
	userErrRepoMock := new(mocks.UserRepository)
	userErrRepoMock.On("FindUserByID", "incorrect id").Return(model.User{}, errors.New("incorrect id"))
	testObj = service.NewChatService(userErrRepoMock, chatRepoMock, messageRepoMock)
	chatsResponse, err = testObj.GetChats(chatsErrRequest)
	assert.Error(err)
	var responseError *errs.ResponseError
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(chatsResponse)

	chatErrRepoMock := new(mocks.ChatRepository)
	chatErrRepoMock.On("FindChats", userModel).Return([]model.Chat{}, errors.New("internal db error"))
	testObj = service.NewChatService(userRepoMock, chatErrRepoMock, messageRepoMock)
	chatsResponse, err = testObj.GetChats(chatsRequest)
	assert.Error(err)
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(chatsResponse)
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

	messagesErrRequest := view.MessagesRequest{
		СhatID: "incorrect id",
	}
	chatErrRepoMock := new(mocks.ChatRepository)
	chatErrRepoMock.On("FindChatByID", "incorrect id").Return(model.Chat{}, errors.New("incorrect id"))
	testObj = service.NewChatService(userRepoMock, chatErrRepoMock, messageRepoMock)
	messagesResponse, err = testObj.GetMessages(messagesErrRequest)
	assert.Error(err)
	var responseError *errs.ResponseError
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(messagesResponse)

	messageErrRepoMock := new(mocks.MessageRepository)
	messageErrRepoMock.On("FindMessages", chatModel).Return([]model.Message{}, errors.New("internal db error"))
	testObj = service.NewChatService(userRepoMock, chatRepoMock, messageErrRepoMock)
	messagesResponse, err = testObj.GetMessages(messagesRequest)
	assert.Error(err)
	if errors.As(err, &responseError) {
		assert.Equal(404, responseError.Status)
	}
	assert.Empty(messagesResponse)
}
