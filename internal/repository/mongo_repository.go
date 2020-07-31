package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/flaambe/avito/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	Db *mongo.Database
}

type chatRepository struct {
	Db *mongo.Database
}

type messageRepository struct {
	Db *mongo.Database
}

func NewUserRepository(db *mongo.Database) UserRepository {
	return &userRepository{db}
}

func NewChatRepository(db *mongo.Database) ChatRepository {
	return &chatRepository{db}
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{db}
}

// User
func (u *userRepository) FindUserByID(id string) (model.User, error) {
	user := model.User{}
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.User{}, err
	}

	err = u.Db.Collection("users").FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) InsertUser(name string) (string, error) {
	result, err := u.Db.Collection("users").InsertOne(context.TODO(), bson.M{"username": name})
	if err != nil {
		return "-1", err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

// Chat
func (c *chatRepository) FindChatByID(id string) (model.Chat, error) {
	chat := model.Chat{}

	chatId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Chat{}, err
	}

	err = c.Db.Collection("chats").FindOne(context.TODO(), bson.M{"_id": chatId}).Decode(&chat)
	if err != nil {
		return model.Chat{}, err
	}

	return chat, nil
}

func (c *chatRepository) FindChats(user model.User) ([]model.Chat, error) {
	chats := []model.Chat{}

	cur, err := c.Db.Collection("chats").Find(context.TODO(), bson.M{"users": user})
	if err != nil {
		return []model.Chat{}, err
	}

	err = cur.All(context.TODO(), &chats)
	if err != nil {
		return []model.Chat{}, err
	}

	return chats, nil
}

func (c *chatRepository) InsertChat(name string, users []model.User) (string, error) {
	chat := model.Chat{
		Name:      name,
		Users:     users,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	result, err := c.Db.Collection("chats").InsertOne(context.TODO(), chat)
	if err != nil {
		return "-1", err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

// Message
func (m *messageRepository) InsertMessage(chat model.Chat, user model.User, text string) (string, error) {
	message := model.Message{
		Chat:      chat.ID,
		Author:    user.ID,
		Text:      text,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	result, err := m.Db.Collection("messages").InsertOne(context.TODO(), message)
	if err != nil {
		return "-1", err
	}

	oid, _ := result.InsertedID.(primitive.ObjectID)

	return oid.Hex(), nil
}

func (m *messageRepository) FindMessages(chat model.Chat) ([]model.Message, error) {
	messages := []model.Message{}

	cur, err := m.Db.Collection("messages").Find(context.TODO(), bson.M{"chat": chat.ID})
	if err != nil {
		return []model.Message{}, err
	}

	err = cur.All(context.TODO(), &messages)
	if err != nil {
		return []model.Message{}, err
	}

	return messages, nil
}
