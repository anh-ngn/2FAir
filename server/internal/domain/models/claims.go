package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `json:"user_id"`
	SessionID string             `json:"session_id"`
	jwt.StandardClaims
}
