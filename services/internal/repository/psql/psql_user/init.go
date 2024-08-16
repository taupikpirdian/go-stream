package psql_user

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepositoryFactory struct {
	Db *gorm.DB
}

func (d *UserRepositoryFactory) Create() (repository.ChatUsersRepository, error) {
	if d.Db == nil {
		return nil, exceptions.ErrorRequired("Db")
	}
	return &userRepository{
		db: d.Db,
	}, nil
}
