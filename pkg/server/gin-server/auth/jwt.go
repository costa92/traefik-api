package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"treafik-api/pkg/server"
)

const (
	UserKey = "id"
)

type Authorization struct {
	JwtConfig server.AuthJwt
}

func NewAuthorization(jwtConfig server.AuthJwt) *Authorization {
	return &Authorization{
		JwtConfig: jwtConfig,
	}
}

// GenerateJWT 生产
func (a *Authorization) GenerateJWT(username string) (string, error) {
	exportTime := time.Duration(a.JwtConfig.ExportTime) * time.Hour
	claims := UserInfo{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exportTime)), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                 // 生效时间
			Issuer:    "test",
		},
	}
	// 使用HS256签名算法
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString([]byte(a.JwtConfig.AppSecret))

	return s, err
}

// ParseJwt 验证
func (a *Authorization) ParseJwt(token string) (*UserInfo, error) {
	t, err := jwt.ParseWithClaims(token, &UserInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.JwtConfig.AppSecret), nil
	})
	if claims, ok := t.Claims.(*UserInfo); ok && t.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
