package repo

import "chat/internal/domain/model"

type ISessionsRepo interface {
	Add(session *model.Session) error
	IsTokenValid(token string) bool
}
