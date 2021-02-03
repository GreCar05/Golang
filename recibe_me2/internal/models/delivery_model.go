package models

import "gopkg.in/mgo.v2/bson"

const (
	DELIVERY_STATE_INVALID = iota
	DELIVERY_STATE_SENT
	DELIVERY_STATE_IN_TRANSIT
	DELIVERY_STATE_RECEIVED
)

type DeliveryState int64

type DeliveryModel struct {
	ID               bson.ObjectId          `json:"id" bson:"_id,omitempty" `
	Rated            bool                   `json:"rated" bson:"rated"`
	Rating           int64                  `json:"rating" bson:"rating"`
	Courier          CourierModel           `json:"courier" bson:"courier"`
	Description      string                 `json:"description" bson:"description"`
	DeliveryState    DeliveryState          `json:"delivery_state" bson:"delivery_state"`
	DeliveryHistory  []DeliveryHistoryModel `json:"delivery_history" bson:"delivery_history"`
	EstimatedDueDate int64                  `json:"estimated_due_date" bson:"estimated_due_date"`
}

type DeliveryHistoryModel struct {
	Timestamp     int64         `json:"timestamp" bson:"timestamp"`
	Description   string        `json:"description" bson:"description"`
	DeliveryState DeliveryState `json:"delivery_state" bson:"delivery_state"`
}

type CourierModel struct {
	Name    string `json:"name" bson:"name"`
	IconURL string `json:"icon_url" bson:"icon_url"`
}
