package repo

import (
	"chat/internal/domain/model"
	"context"
	"time"
)

type IUsersRepository interface {
	CreateUser(ctx context.Context, q *CreateUserQuery) error
	GetUserByNick(ctx context.Context, q *GetUserByNickQuery) (*model.User, error)
}

type CreateUserQuery struct {
	Nick         string
	PasswordHash string
	CreatedAt    time.Time
}

type GetUserByNickQuery struct {
	Nick string
}
