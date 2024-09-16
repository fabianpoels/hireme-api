package models

import (
	"hireme-api/config"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Participant struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SessionId string             `json:"sessionId" bson:"sessionId"`
	Email     string             `json:"email" bson:"email"`
	Page      string             `json:"page" bson:"page"`
	Score     int                `json:"score" bson:"score"`
	Pages     Page               `json:"pages" bson:"pages"`
	Finished  bool               `json:"finished" bson:"finished"`
	CreatedAt time.Time          `json:"-" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedAt,omitempty"`
}

func GetParticipantCollection(client mongo.Client) *mongo.Collection {
	return client.Database(config.GetConfig().GetString("database")).Collection("participants")
}