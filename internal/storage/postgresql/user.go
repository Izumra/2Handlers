package postgresql

import (
	"context"
	"errors"

	"github.com/Izumra/2Handlers/domain/entity"
	"github.com/Izumra/2Handlers/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) ByUsername(ctx context.Context, username string) (entity.User, error) {
	row := s.db.QueryRow(
		ctx,
		"select Username, Password from Users where Username=$1",
		username,
	)

	var user entity.User
	if err := row.Scan(&user.Username, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, storage.ErrUserNotFound
		}

		return entity.User{}, err
	}

	return user, nil
}
