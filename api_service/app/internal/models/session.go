package models

import (
	"encoding/json"
	"time"
)

type SessionResponse struct {
	Token  Token  `json:"token"`
	Role   string `json:"role"`
	UserId string `json:"userId"`
}

type Token struct {
	AccessToken string `json:"accessToken"`
	Exp         int64  `json:"exp"`
}

type SignInUserDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type LimitData struct {
	ClientIP string
	Count    int32
	Exp      time.Duration
}

func (i LimitData) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
