package app

import (
	"chat/internal/app/usecase"
	"chat/internal/domain/model"
	"chat/internal/domain/service"
	_go "chat/internal/infrastructure/go"
	"chat/internal/infrastructure/postgres"
	"chat/internal/registry"
)

func NewContainer(
	config *Config,
) (*registry.Container, error) {

	db, err := postgres.NewConnectDB(
		config.PostgresConfig.Dsn,
		config.PostgresConfig.ServiceName,
		config.PostgresConfig.DBPoolSize,
		config.PostgresConfig.ReadTimeout,
		config.PostgresConfig.WriteTimeout,
	)
	if err != nil {
		return nil, err
	}

	usersRepo := postgres.NewUsersRepository(db)
	sessionsRepo := _go.NewSessionsRepository(make(map[string]model.Session))
	authService := service.NewAuthService(sessionsRepo)
	loginUsecase := usecase.NewUserLoginUsecase(usersRepo, authService)
	registerUsecase := usecase.NewUserRegisterUsecase(usersRepo)

	container := registry.Container{
		Database:        db,
		UsersRepo:       usersRepo,
		SessionsRepo:    sessionsRepo,
		LoginUsecase:    loginUsecase,
		AuthService:     authService,
		RegisterUsecase: registerUsecase,
	}

	return &container, nil
}
