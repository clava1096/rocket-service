package model

import "github.com/google/uuid"

type SessionDB struct {
	SessionID uuid.UUID
	UserID    uuid.UUID
}

func (s *SessionDB) Key(prefix string) string {
	return prefix + s.SessionID.String()
}

func (s *SessionDB) Value() string {
	return s.UserID.String()
}
