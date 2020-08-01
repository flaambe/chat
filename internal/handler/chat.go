package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/flaambe/avito/internal/errs"
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
		respondWithError(w, http.StatusBadRequest, "json is invalid: "+err.Error())
		return
	}

	if body.UserName == "" {
		respondWithError(w, http.StatusBadRequest, "username not found")
		return
	}

	response, err := c.chatService.AddUser(body)
	if err != nil {
		var responseError *errs.ResponseError
		if errors.As(err, &responseError) {
			if responseError.Err != nil {
				log.Println(responseError.Err.Error())
			}

			respondWithError(w, responseError.Status, responseError.Message)
			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (c *ChatHandler) AddChat(w http.ResponseWriter, r *http.Request) {
	var body view.NewChatRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "json is invalid: "+err.Error())
		return
	}

	if body.Name == "" || len(body.UsersID) < 1 {
		respondWithError(w, http.StatusBadRequest, "name or users not found")
		return
	}

	response, err := c.chatService.AddChat(body)
	if err != nil {
		var responseError *errs.ResponseError
		if errors.As(err, &responseError) {
			if responseError.Err != nil {
				log.Println(responseError.Err.Error())
			}

			respondWithError(w, responseError.Status, responseError.Message)

			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (c *ChatHandler) AddMessage(w http.ResponseWriter, r *http.Request) {
	var body view.NewMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "json is invalid: "+err.Error())

		return
	}

	if body.Text == "" || body.ChatID == "" || body.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "chat, user or text not found")
		return
	}

	response, err := c.chatService.AddMessage(body)
	if err != nil {
		var responseError *errs.ResponseError
		if errors.As(err, &responseError) {
			if responseError.Err != nil {
				log.Println(responseError.Err.Error())
			}

			respondWithError(w, responseError.Status, responseError.Message)

			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func (c *ChatHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	var body view.ChatsRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "json is invalid: "+err.Error())

		return
	}

	response, err := c.chatService.GetChats(body)
	if err != nil {
		var responseError *errs.ResponseError
		if errors.As(err, &responseError) {
			if responseError.Err != nil {
				log.Println(responseError.Err.Error())
			}

			respondWithError(w, responseError.Status, responseError.Message)

			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func (c *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var body view.MessagesRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "json is invalid: "+err.Error())

		return
	}

	response, err := c.chatService.GetMessages(body)
	if err != nil {
		var responseError *errs.ResponseError
		if errors.As(err, &responseError) {
			if responseError.Err != nil {
				log.Println(responseError.Err.Error())
			}

			respondWithError(w, responseError.Status, responseError.Message)

			return
		}

		respondWithError(w, http.StatusInternalServerError, err.Error())

		return
	}

	respondWithJSON(w, http.StatusOK, response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, view.ErrorResponse{ErrorMessage: view.ErrorDetails{Code: code, Message: message}})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
