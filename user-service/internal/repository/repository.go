package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	customErros "social-network/user-service/internal/errors"
	"social-network/user-service/internal/logger"
	"time"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) RegisterUser(user *User) error {
	oldUser, _ := ur.GetUserByLogin(user.Login)
	if oldUser != nil {
		return &customErros.LoginAlreadyTakenError{}
	}

	user.UpdatedAt = time.Now()
	user.RegisteredAt = time.Now()

	_, err := ur.db.NewInsert().
		Model(user).
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("failed to insert user: %v", err))
		return err
	}

	return nil
}

func (ur *UserRepository) LoginUser(user *User) error {
	err := ur.db.NewSelect().
		Model(user).
		Where("login = ? and password = ?", user.Login, user.Password).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &customErros.NotFoundUserError{}
		}

		logger.Error(fmt.Sprintf("failed to login user: %v", err))
		return err
	}

	return nil
}

func (ur *UserRepository) GetUserByLogin(login string) (*User, error) {
	user := &User{}
	err := ur.db.NewSelect().
		Model(user).
		Where("login = ?", login).
		Scan(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &customErros.NotFoundUserError{}
		}
		logger.Error(fmt.Sprintf("failed to query user: %v", err))
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) UpdateUserProfile(login string, user *User) error {
	user.UpdatedAt = time.Now()
	_, err := ur.db.NewUpdate().
		Model(user).
		Where("login = ?", login).
		OmitZero().
		Exec(context.Background())
	if err != nil {
		logger.Error(fmt.Sprintf("failed to update user: %v", err))
		return err
	}

	return nil
}
