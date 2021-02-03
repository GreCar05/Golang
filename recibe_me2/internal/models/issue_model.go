package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Id bson.ObjectId `json:"deliveryId" validate:"deliveryId"`

type IssueModel struct {
	ID          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	DeliveryID  string        `json:"delivery_id" bson:"delivery_id"`
	Description string        `json:"description" bson:"description"`
	IssueType   int64         `json:"issue_type" bson:"issue_type" `
}
