package application

import (
	"tech-db-project/domain/entity"
	"tech-db-project/domain/repository"
	"tech-db-project/infrasctructure/tools"
)

type ForumApp struct {
	forumRepo repository.ForumRepo
	userRepo  repository.UserRepo
}

func NewForumApp(forumRepo repository.ForumRepo, userRepo repository.UserRepo) *ForumApp {
	return &ForumApp{forumRepo, userRepo}
}

func (forumApp *ForumApp) CreateForum(f *entity.Forum) error {
	u := &entity.User{}
	u.Nickname = f.User
	err := forumApp.userRepo.GetByNickname(u)
	if err != nil {
		return tools.UserNotExist
	}

	f.User = u.Nickname
	err = forumApp.forumRepo.GetBySlug(f)
	if err == nil {
		return tools.ForumExist
	}

	err = forumApp.forumRepo.InsertInto(f)
	if err != nil {
		return err
	}

	return nil
}

func (forumApp *ForumApp) GetForum(f *entity.Forum) error {
	err := forumApp.forumRepo.GetBySlug(f)
	if err != nil {
		return tools.ForumNotExist
	}

	return nil
}

func (forumApp *ForumApp) GetForumThreads(f *entity.Forum, desc, limit, since string) (entity.Threads, error) {
	err := forumApp.forumRepo.GetBySlug(f)
	if err != nil {
		return nil, tools.ForumNotExist
	}

	ths, err := forumApp.forumRepo.GetThreads(f, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return ths, nil
}

func (forumApp *ForumApp) GetForumUsers(f *entity.Forum, desc, limit, since string) (entity.Users, error) {
	err := forumApp.forumRepo.GetBySlug(f)
	if err != nil {
		return nil, tools.ForumNotExist
	}
	usr, err := forumApp.forumRepo.GetUsers(f, desc, limit, since)
	if err != nil {
		return nil, err
	}

	return usr, nil
}
