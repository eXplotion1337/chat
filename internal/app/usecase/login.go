package usecase

import (
	"chat/internal/domain/model"
	"chat/internal/domain/repo"
	"chat/internal/domain/service"
	"chat/pkg/wrap"
	"context"
)

type UserLoginInput struct {
	Username string
	Password string
}

type UserLoginUsecase struct {
	userRepo    repo.IUsersRepository
	authService *service.AuthService
}

func NewUserLoginUsecase(
	userRepo repo.IUsersRepository,
	authService *service.AuthService,
) *UserLoginUsecase {
	return &UserLoginUsecase{userRepo, authService}
}

func (u *UserLoginUsecase) ProcessAuth(i *UserLoginInput) (*model.Session, error) {
	q := repo.GetUserByNickQuery{Nick: i.Username}
	user, err := u.userRepo.GetUserByNick(context.Background(), &q)
	if err != nil {
		return nil, wrap.Errorf("failed to process auth: %w", err)
	}

	serviceInput := service.GetTokenInput{
		User:     user,
		Password: i.Password,
	}

	session, err := u.authService.GetToken(serviceInput)
	if err != nil {
		return nil, wrap.Errorf("failed to process auth: %w", err)
	}

	return session, nil
}
