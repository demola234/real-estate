package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                      primitive.ObjectID `bson:"_id"`
	First_Name              *string            `json:"first_name" validate:"required,min=2,max=100"`
	Last_Name               *string            `json:"last_name" validate:"required,min=2,max=100"`
	Password                *string            `json:"password" validate:"required,min=6"`
	Email                   *string            `json:"email" validate:"required"`
	Otp_Verified            *bool              `json:"otp_verified"`
	Otp                     *string            `json:"otp"`
	Otp_Expiration          time.Time          `json:"otp_expiration"`
	Avatar                  *string            `json:"avatar"`
	Phone                   *string            `json:"phone" validate:"required"`
	Token                   *string            `json:"token"`
	Push_Notification_Token *string            `json:"push_notification_token"`
	Refresh_Token           *string            `json:"refresh_token"`
	Created_at              time.Time          `json:"created_at"`
	Updated_at              time.Time          `json:"updated_at"`
	User_id                 string             `json:"user_id"`
}
