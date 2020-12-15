package repository

import "tech-db-project/domain/entity"

type UserRepo interface {
	InsertInto(user *entity.User) error

	GetByNickname(user *entity.User) error
	GetByNicknameOrEmail(user *entity.User) (entity.Users, error)
	GetStatus(status *entity.Status) error

	Update(user *entity.User) error

	DeleteAll() error

	Prepare() error
}
