package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/s0vunia/password-manager/internal/domain"
	"github.com/s0vunia/password-manager/internal/lib/jwt"
	"github.com/s0vunia/password-manager/internal/lib/logger/sl"
	"github.com/s0vunia/password-manager/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type IOAuth interface {
	Login(
		ctx context.Context,
		login string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		login string,
		password string,
	) (userID uuid.UUID, err error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserSaver interface {
	Create(ctx context.Context, login string, passHash []byte) (uuid.UUID, error)
}

type UserProvider interface {
	Get(ctx context.Context, login string) (*domain.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (domain.App, error)
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// If user doesn't exist, returns error.
func (a *Auth) Login(
	ctx context.Context,
	login string,
	password string,
	appID int,
) (string, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", login),
	)

	log.Info("attempting to login user")

	user, err := a.usrProvider.Get(ctx, login)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("failed to get user", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		a.log.Info("failed to get app", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(*user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context.Context, login string, pass string) (uuid.UUID, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.Create(ctx, login, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))

		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
