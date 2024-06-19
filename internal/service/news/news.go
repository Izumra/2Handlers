package news

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Izumra/2Handlers/domain/dto/requests"
	"github.com/Izumra/2Handlers/domain/entity"
	"github.com/Izumra/2Handlers/utils/JWT"
	"github.com/Izumra/2Handlers/utils/config"
)

var (
	ErrForbidden = errors.New("Доступ с данным токеном доступа запрещен")
)

type NewsRepository interface {
	Edit(ctx context.Context, id int, data requests.NewsData) error
	List(ctx context.Context, offset, count int) ([]entity.News, error)
}

type Service struct {
	log     *slog.Logger
	config  config.Token
	newsRep NewsRepository
}

func New(log *slog.Logger, config config.Token, nrep NewsRepository) *Service {
	return &Service{
		log,
		config,
		nrep,
	}
}

func (s *Service) EditNews(
	ctx context.Context,
	accessToken string,
	id int,
	data requests.NewsData,
) error {
	_, err := JWT.ValidateToken(accessToken, s.config.Secret)
	if err != nil {
		s.log.Info("Причина возникновения ошибки при проверке токена доступа", slog.Any("причина", err))
		return ErrForbidden
	}

	err = s.newsRep.Edit(ctx, id, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ListNews(
	ctx context.Context,
	accessToken string,
	offset int,
	count int,
) ([]entity.News, error) {
	_, err := JWT.ValidateToken(accessToken, s.config.Secret)
	if err != nil {
		s.log.Info("Причина возникновения ошибки при проверке токена доступа", slog.Any("причина", err))
		return nil, ErrForbidden
	}

	news, err := s.newsRep.List(ctx, offset, count)
	if err != nil {
		return nil, err
	}

	return news, nil
}
