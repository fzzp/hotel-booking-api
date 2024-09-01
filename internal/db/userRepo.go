package db

import (
	"context"

	"github.com/fzzp/hotel-booking-api/internal/models"
)

type UserRepo interface {
	InsertOne(user *models.User) (id uint, err error)
	GetOneByUq(map[string]string) (*models.User, error)
	UpdateOne(uid uint, user *models.User) error
}

var _ UserRepo = (*userRepo)(nil)

type userRepo struct {
	DB Queryable
}

func NewUserRepo(qb Queryable) *userRepo {
	return &userRepo{
		DB: qb,
	}
}

// Insert 添加用户
func (store *userRepo) InsertOne(user *models.User) (id uint, err error) {
	sql := `
	insert into users(
			phone_number,
			password_hash,
			username,
			avatar)
	values(?,?,?,?);`
	return create(store.DB, sql, user.PhoneNumber, user.PasswordHash, user.Username, user.Avatar)
}

// GetOneByUq 根据唯一字段查询用户信息.
func (store *userRepo) GetOneByUq(fields map[string]string) (*models.User, error) {
	sql := `select 
		id, phone_number, username, password_hash, avatar
	from users
	where 
	`
	query, params := unqQuerySQL(fields, "?")
	sql += query + " and is_deleted = 1"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	var user models.User
	err := store.DB.GetContext(ctx, &user, sql, params...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (store *userRepo) UpdateOne(uid uint, user *models.User) error {
	sql := `update users
		set
			username=?,
			avatar=?,
			updated_at=now()
		where
		 id = ? and is_deleted = 1;
	`
	return update(store.DB, sql, user.Username, user.Avatar, uid)
}
