package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Chat struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Users     []User             `bson:"users"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
