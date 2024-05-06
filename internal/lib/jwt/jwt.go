package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/s0vunia/password-manager/internal/domain"
	"time"
)

func NewToken(user domain.User, app domain.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
