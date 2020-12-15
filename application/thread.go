package application

import (
	"tech-db-project/domain/entity"
	"tech-db-project/domain/repository"
	"tech-db-project/infrasctructure/tools"
)

type ThreadApp struct {
	threadRepo repository.ThreadRepo
	forumRepo  repository.ForumRepo
}

func NewThreadApp(threadRepo repository.ThreadRepo, forumRepo repository.ForumRepo) *ThreadApp {
	return &ThreadApp{threadRepo, forumRepo}
}

func (threadApp *ThreadApp) CreateThread(th *entity.Thread) error {
	f := &entity.Forum{}
	f.Slug = th.Forum

	err := threadApp.forumRepo.GetBySlug(f)
	if err != nil {
		return tools.UserNotExist
	}

	th.Forum = f.Slug
	err = threadApp.threadRepo.InsertInto(th)
	if err != nil {
		err = threadApp.threadRepo.GetBySlugOrId(th)
		if err != nil {
			return tools.UserNotExist
		}

		return tools.ThreadExist
	}

	return nil
}

func (threadApp *ThreadApp) GetThreadInfo(th *entity.Thread) error {
	err := threadApp.threadRepo.GetBySlugOrId(th)
	if err != nil {
		return tools.ThreadNotExist
	}

	return nil
}

func (threadApp *ThreadApp) CreateVote(th *entity.Thread, v *entity.Vote) error {
	err := threadApp.threadRepo.GetBySlugOrId(th)
	if err != nil {
		return tools.ThreadNotExist
	}

	v.Thread = th.Id
	err = threadApp.threadRepo.InsertIntoVotes(th, v)
	if err != nil {
		return tools.UserNotExist
	}

	return nil
}

func (threadApp *ThreadApp) UpdateThread(th *entity.Thread) error {
	err := threadApp.threadRepo.Update(th)
	if err != nil {
		return tools.ThreadNotExist
	}

	return nil
}

func (threadApp *ThreadApp) GetThreadPosts(th *entity.Thread, desc, sort, limit, since string) ([]entity.Post, error) {
	err := threadApp.threadRepo.GetBySlugOrId(th)
	if err != nil {
		return nil, tools.ThreadNotExist
	}

	posts, err := threadApp.threadRepo.GetPosts(th, desc, sort, limit, since)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
