package models

import "gopkg.in/mgo.v2/bson"

type UserModel struct {
	ID                bson.ObjectId `json:"id" bson:"_id,omitempty" `
	FirstName         string        `json:"first_name" bson:"first_name"`
	LastName          string        `json:"last_name" bson:"last_name"`
	Phone             string        `json:"phone" bson:"phone"`
	Email             string        `json:"email" bson:"email"`
	Password          string        `json:"password" bson:"password"`
	Verified          bool          `json:"verified" bson:"verified"`
	VerificationCode  string        `json:"verification_code" bson:"verification_code"`
}

type UserAuthModel struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}

type UserLogResponseModel struct {
	User  UserModel `json:"user"`
	Token string    `json:"token"`
}
