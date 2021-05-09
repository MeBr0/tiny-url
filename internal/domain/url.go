package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type URL struct {
	Alias     string             `json:"alias" bson:"_id,omitempty"`
	Original  string             `json:"original" bson:"original"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	ExpiredAt time.Time          `json:"expiredAt" bson:"expiredAt"`
	Owner     primitive.ObjectID `json:"owner" bson:"owner"`
}

type URLCreate struct {
	Original string             `json:"original" bson:"original"`
	Owner    primitive.ObjectID `bson:"owner"`
}
