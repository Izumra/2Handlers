package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Izumra/2Handlers/domain/dto/requests"
	"github.com/Izumra/2Handlers/domain/entity"
	"github.com/Izumra/2Handlers/internal/storage"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) Edit(ctx context.Context, id int, data requests.NewsData) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	noveltyDB := tx.QueryRow(
		ctx,
		`select * 
		from News 
		where Id = $1
		`,
		id,
	)

	var novelty entity.News
	if err := noveltyDB.Scan(
		&novelty.Id,
		&novelty.Title,
		&novelty.Content,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storage.ErrNoveltyNotFound
		}

		return err
	}

	categoriesDB, err := tx.Query(
		ctx,
		`select CategoryId
		from NewsCategories 
		where NewsId = $1`,
		id,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	} else if err == nil {
		for categoriesDB.Next() {
			var categoryID int

			categoriesDB.Scan(&categoryID)
			novelty.Categories = append(novelty.Categories, categoryID)
		}
	}
	categoriesDB.Close()

	added := ""
	prevCategories := make(map[int]bool)
	for _, v := range novelty.Categories {
		prevCategories[v] = true
	}
	for i, v := range data.Categories {
		if _, ok := prevCategories[v]; !ok {
			added += fmt.Sprintf("(%d,%d),", id, v)
		} else {
			delete(prevCategories, i)
		}
	}

	if len(prevCategories) != 0 {
		removed := fmt.Sprintf("delete from NewsCategories where NewsId=%d and CategoryId in (", id)

		for i := range prevCategories {
			removed += fmt.Sprintf("%d, ", i)
		}
		removed, _ = strings.CutSuffix(removed, ", ")
		removed += ")"

		_, err = tx.Exec(
			ctx,
			removed,
		)
		if err != nil {
			return err
		}
	}
	added, _ = strings.CutSuffix(added, ",")

	if added != "" {
		statement := fmt.Sprintf(`insert into NewsCategories(NewsId,CategoryId)values %s`, added)
		_, err = tx.Exec(
			ctx,
			statement,
		)
		if err != nil {
			return err
		}
	}

	updateNewsStatement := "update News set "
	if data.Title != novelty.Title && data.Title != "" {
		updateNewsStatement += fmt.Sprintf("Title='%s'", data.Title)
	}
	if data.Content != novelty.Content && data.Content != "" {
		if updateNewsStatement == "update News set " {
			updateNewsStatement += fmt.Sprintf("Content='%s'", data.Content)
		} else {
			updateNewsStatement += fmt.Sprintf(", Content='%s'", data.Content)
		}
	}

	if updateNewsStatement != "update News set " {
		updateNewsStatement += fmt.Sprintf("where Id=%d", id)

		_, err = tx.Exec(
			ctx,
			updateNewsStatement,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
func (s *Storage) List(ctx context.Context, offset, count int) ([]entity.News, error) {
	rows, err := s.db.Query(
		ctx,
		`select * 
	  from News 
	  offset $1
	  limit $2`,
		offset,
		count,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNewsNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var news []entity.News
	for rows.Next() {
		var novelty entity.News

		if err := rows.Scan(
			&novelty.Id,
			&novelty.Title,
			&novelty.Content,
		); err != nil {
			return nil, err
		}

		categoriesIDes, err := s.db.Query(
			ctx,
			`select CategoryId
		from NewsCategories 
		where NewsId=$1`,
			novelty.Id,
		)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		} else if err == nil {
			defer categoriesIDes.Close()

			for categoriesIDes.Next() {
				var categoryID int
				if err := categoriesIDes.Scan(&categoryID); err != nil {
					return nil, err
				}
				novelty.Categories = append(novelty.Categories, categoryID)
			}
		}

		news = append(news, novelty)
	}

	return news, nil
}
