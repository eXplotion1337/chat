package registry

import (
	"chat/internal/app/usecase"
	"chat/internal/domain/repo"
	"chat/internal/domain/service"
	"github.com/uptrace/bun"
)

type Container struct {
	Database        *bun.DB
	UsersRepo       repo.IUsersRepository
	SessionsRepo    repo.ISessionsRepo
	AuthService     *service.AuthService
	LoginUsecase    *usecase.UserLoginUsecase
	RegisterUsecase *usecase.UserRegisterUsecase
}
