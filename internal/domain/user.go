package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name         string               `json:"name" bson:"name"`
	Email        string               `json:"email" bson:"email"`
	Password     string               `json:"password" bson:"password"`
	RegisteredAt time.Time            `json:"registeredAt" bson:"registeredAt"`
	LastLogin    time.Time            `json:"lastLogin" bson:"lastLogin"`
}

type UserRegister struct {
	Name         string               `json:"name"`
	Email        string               `json:"email"`
	Password     string               `json:"password"`
}
