package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	// Unique id
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty" format:"hexidecimal string" example:"6095872d75ff40c9238bdb29"`
	// First name
	Name string `json:"name" bson:"name" example:"Sirius"`
	// Unique email
	Email string `json:"email" bson:"email" format:"email" example:"sirius@gmail.com"`
	// Secret password
	Password string `json:"password" bson:"password" example:"qweqweqwe"`
	// Time of registration
	RegisteredAt time.Time `json:"registeredAt" bson:"registeredAt" format:"yyyy-MM-ddThh:mm:ss.ZZZ" example:"2021-05-07T18:30:05.365Z"`
	// Last login time
	LastLogin time.Time `json:"lastLogin" bson:"lastLogin" format:"yyyy-MM-ddThh:mm:ss.ZZZ" example:"2021-05-07T18:30:05.365Z"`
} // @name User

type UserRegister struct {
	// First name
	Name string `json:"name" example:"Sirius"`
	// Unique email
	Email string `json:"email" format:"email" example:"sirius@gmail.com"`
	// Secret password
	Password string `json:"password" example:"qweqweqwe"`
} // @name UserRegister

type UserLogin struct {
	// Unique email
	Email string `json:"email" format:"email" example:"sirius@gmail.com"`
	// Secret password
	Password string `json:"password" example:"qweqweqwe"`
} // @name UserLogin

type Tokens struct {
	// Token used for accessing operations and/or resources
	AccessToken string `json:"accessToken" example:"access token"`
} // @name Tokens
