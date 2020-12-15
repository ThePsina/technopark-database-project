package repository

import "tech-db-project/domain/entity"

type ThreadRepo interface {
	GetBySlug(thread *entity.Thread) error
	GetById(thread *entity.Thread) error
	GetBySlugOrId(thread *entity.Thread) error
	GetPosts(thread *entity.Thread, desc, sort, limit, since string) (entity.Posts, error)

	InsertInto(thread *entity.Thread) error
	InsertIntoVotes(thread *entity.Thread, vote *entity.Vote) error
	InsertIntoForumUsers(forum, nickname string) error

	Update(thread *entity.Thread) error

	Prepare() error
}
