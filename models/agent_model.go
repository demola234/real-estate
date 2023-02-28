package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agents struct {
	ID                      primitive.ObjectID `bson:"_id"`
	Created_at              time.Time          `json:"created_at"`
	Updated_at              time.Time          `json:"updated_at"`
	User_id                 string             `json:"user_id"`
}
