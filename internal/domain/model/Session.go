package model

import "time"

type Session struct {
	Token     Token
	ExpiresAt time.Time
}

type Token string

func (s Session) isExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
