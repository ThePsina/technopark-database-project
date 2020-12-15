package repository

import "tech-db-project/domain/entity"

type ForumRepo interface {
	InsertInto(forum *entity.Forum) error
	GetBySlug(forum *entity.Forum) error
	GetThreads(forum *entity.Forum, desc, limit, since string) (entity.Threads, error)
	GetUsers(forum *entity.Forum, desc, limit, since string) (entity.Users, error)
	Prepare() error
}
