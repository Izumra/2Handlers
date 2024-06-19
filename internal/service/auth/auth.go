package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Izumra/2Handlers/domain/entity"
	"github.com/Izumra/2Handlers/utils/JWT"
	"github.com/Izumra/2Handlers/utils/config"
)

var (
	ErrUserNotFound = errors.New("Пользователь системы не зарегестрирован")
	ErrCreateToken  = errors.New("Произошла ошибка при создании токена доступа")
)

type UserRepository interface {
	ByUsername(ctx context.Context, username string) (entity.User, error)
}

type Service struct {
	log     *slog.Logger
	config  config.Token
	userRep UserRepository
}

func New(log *slog.Logger, urep UserRepository, config config.Token) *Service {
	return &Service{
		log,
		config,
		urep,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRep.ByUsername(ctx, username)
	if err != nil {
		return "", ErrUserNotFound
	}
	if password != user.Password {
		return "", fmt.Errorf("Пароли учетных записей не совпадают")
	}

	token, err := JWT.CreateToken(user.Id, s.config.Secret, s.config.TTL)
	if err != nil {
		s.log.Info("Причина возникновения ошибки при создании токена доступа", slog.Any("причина", err))
		return "", ErrCreateToken
	}

	return token, nil
}
