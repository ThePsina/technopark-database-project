package application

import (
	"tech-db-project/domain/entity"
	"tech-db-project/domain/repository"
	"tech-db-project/infrasctructure/tools"
)

type UserApp struct {
	userRepo repository.UserRepo
}

func NewUserApp(userRepo repository.UserRepo) *UserApp {
	return &UserApp{userRepo}
}

func (userApp *UserApp) CreateUser(u *entity.User) (entity.Users, error) {
	err := userApp.userRepo.InsertInto(u)
	if err != nil {
		users, e := userApp.userRepo.GetByNicknameOrEmail(u)
		err = e
		tools.HandleError(err)
		return users, tools.UserExist
	}

	return entity.Users{}, nil
}

func (userApp *UserApp) GetUser(u *entity.User) error {
	err := userApp.userRepo.GetByNickname(u)
	if err != nil {
		return tools.UserNotExist
	}

	return nil
}

func (userApp *UserApp) UpdateUser(u *entity.User) error {
	uInfo := *u
	if err := userApp.userRepo.GetByNickname(&uInfo); err != nil {
		return tools.UserNotExist
	}

	if u.Email == "" {
		u.Email = uInfo.Email
	}

	if u.About == "" {
		u.About = uInfo.About
	}

	if u.Fullname == "" {
		u.Fullname = uInfo.Fullname
	}

	if err := userApp.userRepo.Update(u); err != nil {
		return tools.UserNotUpdated
	}

	return nil
}

func (userApp *UserApp) DeleteAll() error {
	err := userApp.userRepo.DeleteAll()
	if err != nil {
		return err
	}
	return nil
}

func (userApp *UserApp) GetStatus(s *entity.Status) error {
	err := userApp.userRepo.GetStatus(s)
	if err != nil {
		return err
	}

	return nil
}
