package auth

import (
	"testing"

	"treafik-api/pkg/server"
)

var AppSecret = "1111"

func TestGenerateJWT(t *testing.T) {
	m := NewAuthorization(server.AuthJwt{
		AppSecret:  AppSecret,
		ExportTime: 24,
	})
	res, err := m.GenerateJWT("costalong")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func Test_ParseJwt(t *testing.T) {
	m := NewAuthorization(server.AuthJwt{
		AppSecret:  AppSecret,
		ExportTime: 24,
	})
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MCwidXNlcm5hbWUiOiJjb3N0YWxvbmciLCJleHAiOjE2ODg5MTA2NjIsIm5iZiI6MTY4ODgyNDI2MiwiaWF0IjoxNjg4ODI0MjYyfQ.joXi4n1i-7_28Qlt83-TWnBqLQCIJu30Z25VcD6uirs"
	res, err := m.ParseJwt(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
	t.Log(res.GetExpirationTime())
}
