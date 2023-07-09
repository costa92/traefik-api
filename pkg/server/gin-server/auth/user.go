package auth

import "github.com/golang-jwt/jwt/v5"

type UserInfo struct {
	ID       int64  `json:"id" gorm:"id"`
	Username string `json:"username" gorm:"username" json:"username,omitempty"`
	jwt.RegisteredClaims
}
