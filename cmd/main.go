package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/users/add", chatHandler.AddUser)
	serveMux.HandleFunc("/chats/add", chatHandler.AddChat)
	serveMux.HandleFunc("/chats/get", chatHandler.GetChats)
	serveMux.HandleFunc("/messages/add", chatHandler.AddMessage)
	serveMux.HandleFunc("/messages/get", chatHandler.GetMessages)

	srv := &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      serveMux,
	}

	go func() {
		panic(srv.ListenAndServe())
	}()

	// Create channel for shutdown signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGINT)

	//Recieve shutdown signals.
	<-stop

	// Disconnect database client
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		log.Fatal(err)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down server %s", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}
