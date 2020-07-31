package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Chat      primitive.ObjectID `bson:"chat"`
	Author    primitive.ObjectID `bson:"author"`
	Text      string             `bson:"text"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
