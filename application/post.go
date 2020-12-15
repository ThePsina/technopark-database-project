package application

import (
	"tech-db-project/domain/entity"
	"tech-db-project/domain/repository"
	"tech-db-project/infrasctructure/tools"
)

type PostApp struct {
	postRepo   repository.PostRepo
	threadRepo repository.ThreadRepo
}

func NewPostApp(postRepo repository.PostRepo, threadRepo repository.ThreadRepo) *PostApp {
	return &PostApp{postRepo, threadRepo}
}

func (postApp *PostApp) CreatePosts(p []*entity.Post, th *entity.Thread) error {
	var err error
	if err = postApp.threadRepo.GetBySlugOrId(th); err != nil {
		return tools.ThreadNotExist
	}
	if err = postApp.postRepo.InsertInto(p, th); err != nil {
		if err.Error() == "ERROR: Parent post was created in another thread (SQLSTATE 00404)" {
			return tools.ParentNotExist
		} else {
			return tools.UserNotExist
		}
	}

	for iter := range p {
		if err = postApp.threadRepo.InsertIntoForumUsers(p[iter].Forum, p[iter].Author); err != nil {
			return err
		}
	}

	return nil
}

func (postApp *PostApp) GetPost(p *entity.Post) error {
	err := postApp.postRepo.GetById(p)
	if err != nil {
		return tools.PostNotExist
	}

	return nil
}

func (postApp *PostApp) UpdatePost(p *entity.Post) error {
	message := p.Message

	err := postApp.postRepo.GetById(p)
	if err != nil {
		return tools.PostNotExist
	}

	if message != "" && message != p.Message {
		p.Message = message
		if err := postApp.postRepo.Update(p); err != nil {
			return tools.PostNotExist
		}
	}

	return nil
}
