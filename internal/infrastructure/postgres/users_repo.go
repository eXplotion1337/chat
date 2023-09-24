package postgres

import (
	"chat/internal/domain/model"
	"chat/internal/domain/repo"
	"chat/pkg/wrap"
	"context"
	"github.com/uptrace/bun"
)

type UsersRepository struct {
	db *bun.DB
}

func NewUsersRepository(db *bun.DB) repo.IUsersRepository {
	return &UsersRepository{db}
}

func (r *UsersRepository) GetUserByNick(
	ctx context.Context,
	q *repo.GetUserByNickQuery,
) (*model.User, error) {
	users := []model.User{}
	err := r.db.NewSelect().Model(&users).Where(`u."nick"=?`, q.Nick).Scan(ctx)

	if err != nil {
		return nil, wrap.Errorf(
			"repository error: failed to get user by nick: %s",
			err.Error(),
		)
	}

	if len(users) == 0 {
		return nil, nil
	}

	return &users[0], nil
}

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	q *repo.CreateUserQuery,
) error {
	user := model.User{
		Nick:         q.Nick,
		PasswordHash: q.PasswordHash,
		CreatedAt:    q.CreatedAt,
	}

	_, err := r.db.NewInsert().Model(&user).Exec(ctx)

	if err != nil {
		return wrap.Errorf(
			"repository error: failed to create user: %w",
			err,
		)
	}

	return nil
}
