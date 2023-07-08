package auth

import (
	"errors"
	"time"
)

const (
	PackageName = "internal/domains/auth/"

	RoleOwner  = "owner"
	RoleSeller = "seller"

	AccessKeyTTL  = time.Hour               // 1 hour
	RefreshKeyTTL = time.Hour * 24 * 30 * 2 // 2 months
)

var (
	ErrUsernameTaken       = errors.New("это имя пользователя уже занято")
	ErrUsernameNotFound    = errors.New("пользователь с таким именем не найден")
	ErrIdNotFound          = errors.New("пользователь с таким id не найден")
	ErrWrongPassword       = errors.New("неверный пароль")
	ErrInvalidRefreshToken = errors.New("инвалидный токен для обновления сессии")
	ErrInvalidAccessToken  = errors.New("инвалидный токен доступа")
	ErrDefault             = errors.New("что-то пошло не так")
)
