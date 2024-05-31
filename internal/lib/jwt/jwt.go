package jwt

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/repositories/app"
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

// ProcessJWT Функция для извлечения app_id из JWT и проверки его валидности
func ProcessJWT(ctx context.Context, tokenString string, appRepo app.Repository) (error, *jwt.Token) {
	// Извлечение app_id из JWT без проверки подписи
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return fmt.Errorf("failed to parse JWT: %w", err), nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid JWT"), nil
	}

	appID, ok := claims["app_id"].(float64)
	if !ok {
		return fmt.Errorf("app_id not found in JWT"), nil
	}
	app, err := appRepo.App(ctx, int64(appID))
	if err != nil {
		return fmt.Errorf("failed to get app: %w", err), nil
	}

	// Проверка валидности JWT с использованием полученного секрета
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.Secret), nil
	})

	if err != nil {
		return fmt.Errorf("invalid JWT: %w", err), nil
	}
	return nil, token
}
