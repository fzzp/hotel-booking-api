package db

import (
	"context"

	"github.com/fzzp/hotel-booking-api/internal/models"
)

type SessionRepo interface {
	InsertOne(user *models.Session) (id uint, err error)
	GetOneByUq(map[string]string) (*models.Session, error)
}

var _ SessionRepo = (*sessionRepo)(nil)

type sessionRepo struct {
	DB Queryable
}

func NewSessionRepo(qb Queryable) *sessionRepo {
	return &sessionRepo{
		DB: qb,
	}
}

func (store *sessionRepo) InsertOne(ss *models.Session) (id uint, err error) {
	sql := `
	insert into sessions(
		user_id,
		token_id,
		refresh_token,
		client_ip,
		user_agent,
		expires_at
	)values(?,?,?,?,?,?);`
	return create(
		store.DB, sql,
		ss.UserID,
		ss.TokenID,
		ss.RefreshToken,
		ss.ClientIP,
		ss.UserAgent,
		ss.ExpiresAt,
	)
}

func (store *sessionRepo) GetOneByUq(fields map[string]string) (*models.Session, error) {
	sql := `
	select 
		id, user_id, token_id, refresh_token, client_ip, user_agent, expires_at, created_at
	from sessions 
	where 
	`
	query, params := unqQuerySQL(fields, "?")
	sql += query
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	var ss models.Session
	err := store.DB.GetContext(ctx, &ss, sql, params...)
	if err != nil {
		return nil, err
	}
	return &ss, nil
}
