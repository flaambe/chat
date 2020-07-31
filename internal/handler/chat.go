package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/flaambe/avito/internal/view"
)

type ChatService interface {
	AddUser(user view.NewUserRequest) (view.NewUserResponse, error)
	AddChat(chat view.NewChatRequest) (view.NewChatResponse, error)
	AddMessage(message view.NewMessageRequest) (view.NewMessageResponse, error)
	GetChats(chats view.ChatsRequest) (view.ChatsResponse, error)
	GetMessages(messages view.MessagesRequest) (view.MessagesResponse, error)
}

type ChatHandler struct {
	chatService ChatService
}

func NewChatHandler(s ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: s,
	}
}

func (c *ChatHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	var body view.NewUserRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	user, err := c.chatService.AddUser(body)
	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ChatHandler) AddChat(w http.ResponseWriter, r *http.Request) {
	var body view.NewChatRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	user, err := c.chatService.AddChat(body)
	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ChatHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	var body view.NewMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	user, err := c.chatService.AddMessage(body)
	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ChatHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	var body view.ChatsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	user, err := c.chatService.GetChats(body)
	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var body view.MessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	user, err := c.chatService.GetMessages(body)
	if err != nil {
		log.Println(err)
	}

	response, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}
