package models

type WhitelistChat struct {
	Id    string `gorm:"column:id"`
	Email string `gorm:"column:email"`
}

func (WhitelistChat) TableName() string {
	return "feature_whitelist_chats"
}
