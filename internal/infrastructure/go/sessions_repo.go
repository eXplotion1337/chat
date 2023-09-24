package _go

import (
	"chat/internal/domain/model"
)

type SessionsRepository struct {
	storage map[string]model.Session
}

func NewSessionsRepository(
	storage map[string]model.Session,
) *SessionsRepository {
	return &SessionsRepository{storage}
}

func (s *SessionsRepository) Add(session *model.Session) error {
	key := string(session.Token)
	s.storage[key] = *session

	return nil
}

func (s *SessionsRepository) IsTokenValid(token string) bool {
	_, ok := s.storage[token]
	if !ok {
		return false
	}

	return true
}
