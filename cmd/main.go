package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/flaambe/avito/internal/handler"
	"github.com/flaambe/avito/internal/repository"
	"github.com/flaambe/avito/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database connected")

	db := client.Database(os.Getenv("DB_NAME"))

	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatService := service.NewChatService(userRepo, chatRepo, messageRepo)
	chatHandler := handler.NewChatHandler(chatService)

	http.HandleFunc("/users/add", chatHandler.AddUser)
	http.HandleFunc("/chats/add", chatHandler.AddChat)
	http.HandleFunc("/chats/get", chatHandler.GetChats)
	http.HandleFunc("/messages/add", chatHandler.AddMessage)
	http.HandleFunc("/messages/get", chatHandler.GetMessages)

	if err := http.ListenAndServe(os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
