package storage

import "errors"

var (
	ErrUserNotFound    = errors.New("Пользователь не найден")
	ErrNewsNotFound    = errors.New("Новости не найдены")
	ErrNoveltyNotFound = errors.New("Новость не найдена")
)
