package service

import (
	"chat/internal/domain/model"
	"chat/internal/domain/repo"
	"chat/pkg/wrap"
	"crypto/rand"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type GetTokenInput struct {
	Password string
	User     *model.User
}

const CookieStringBytes = 32

type AuthService struct {
	sessionsRepo repo.ISessionsRepo
}

func NewAuthService(
	sessionsRepo repo.ISessionsRepo,
) *AuthService {
	return &AuthService{sessionsRepo}
}

func (s *AuthService) IsTokenValid(token string) bool {
	return s.sessionsRepo.IsTokenValid(token)
}

func (s *AuthService) GetToken(i GetTokenInput) (*model.Session, error) {
	session := &model.Session{}
	if !CheckPasswordHash(i.Password, i.User.PasswordHash) {
		return session, wrap.Errorf("wrong password")
	}

	t, err := GetAuthToken()
	if err != nil {
		return session, wrap.Errorf("failed to get auth token: %w", err)
	}

	session.Token = t
	session.ExpiresAt = time.Now().Add(time.Hour * 24)

	err = s.sessionsRepo.Add(session)
	if err != nil {
		return &model.Session{}, wrap.Errorf("failed to get auth token: %w", err)
	}

	return session, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetAuthToken() (model.Token, error) {
	return getCryptoString(CookieStringBytes)
}

func getCryptoString(bytes int) (model.Token, error) {
	b, err := getCryptoBytes(bytes)
	if err != nil {
		return "", err
	}

	return model.Token(base64.URLEncoding.EncodeToString(b)), nil
}

func getCryptoBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GetPasswordHash(p string) (string, error) {
	pw := []byte(p)
	hashedPw, err := bcrypt.GenerateFromPassword(pw, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPw), nil
}
