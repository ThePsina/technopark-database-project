package repository

import "tech-db-project/domain/entity"

type PostRepo interface {
	InsertInto(posts []*entity.Post, thread *entity.Thread) error
	GetById(post *entity.Post) error
	Update(post *entity.Post) error
	Prepare() error
}
