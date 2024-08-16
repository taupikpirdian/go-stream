package psql_whitelist

import (
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/domain/repository"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/exceptions"
	"gorm.io/gorm"
)

type whiteListRepository struct {
	db *gorm.DB
}

type WhiteListRepositoryFactory struct {
	Db *gorm.DB
}

func (d *WhiteListRepositoryFactory) Create() (repository.WhitelistChatRepository, error) {
	if d.Db == nil {
		return nil, exceptions.ErrorRequired("Db")
	}
	return &whiteListRepository{
		db: d.Db,
	}, nil
}
