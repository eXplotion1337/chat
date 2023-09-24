package usecase

import (
	"chat/internal/domain/repo"
	"chat/internal/domain/service"
	"chat/pkg/wrap"
	"context"
	"time"
)

type UserRegisterInput struct {
	Username string
	Password string
}

type UserRegisterUsecase struct {
	UserRepo repo.IUsersRepository
}

func NewUserRegisterUsecase(
	userRepo repo.IUsersRepository,
) *UserRegisterUsecase {
	return &UserRegisterUsecase{userRepo}
}

func (u *UserRegisterUsecase) ProcessRegister(i UserRegisterInput) error {
	ph, err := service.GetPasswordHash(i.Password)
	if err != nil {
		return wrap.Errorf(
			"failed to register user: password hash error: %w",
			err,
		)
	}

	q := repo.GetUserByNickQuery{
		Nick: i.Username,
	}

	/**
	@todo: сделать работу с контекстом красиво
	*/
	user, err := u.UserRepo.GetUserByNick(context.Background(), &q)

	if err != nil {
		return wrap.Errorf("Failed to check name for uniqueness: %w", err)
	}

	if user != nil {
		return wrap.Errorf("User with this name is already registered")
	}

	cq := &repo.CreateUserQuery{
		Nick:         i.Username,
		PasswordHash: ph,
		CreatedAt:    time.Now().UTC(),
	}

	err = u.UserRepo.CreateUser(context.Background(), cq)
	if err != nil {
		return wrap.Errorf("insert error", err)
	}

	return nil
}
