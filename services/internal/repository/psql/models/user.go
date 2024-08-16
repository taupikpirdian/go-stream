package models

import "time"

type User struct {
	Id                   string     `gorm:"column:id"`
	Type                 int        `gorm:"column:type"`
	Name                 string     `gorm:"column:name"`
	Email                string     `gorm:"column:email"`
	Msisdn               string     `gorm:"column:msisdn"`
	EmailVerified        bool       `gorm:"column:email_verified"`
	MsisdnVerified       bool       `gorm:"column:msisdn_verified"`
	PhoneNumberTelkomsel bool       `gorm:"column:phone_number_telkomsel"`
	ProfileVerified      bool       `gorm:"column:profile_verified"`
	LinkAjaNumber        string     `gorm:"column:linkaja_number"`
	LinkAjaVerified      bool       `gorm:"column:link_aja_verified"`
	GopayNumber          string     `gorm:"column:gopay_number"`
	GopayVerified        bool       `gorm:"column:gopay_verified"`
	Roles                string     `gorm:"column:roles"`
	Status               int        `gorm:"column:status"`
	Photo                string     `gorm:"column:photo"`
	Gender               string     `gorm:"column:gender"`
	City                 string     `gorm:"column:city"`
	DateOfBirth          time.Time  `gorm:"column:dob"`
	DeviceOs             string     `gorm:"column:device_os"`
	DeviceType           string     `gorm:"column:device_type"`
	BusinessCategory     string     `gorm:"column:business_category"`
	WhatsappNumber       bool       `gorm:"column:whatsapp_number"`
	SubscribeNewsletter  bool       `gorm:"column:subscribe_newsletter"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
	DeletedAt            *time.Time `gorm:"column:deleted_at"`
	Onboarding           bool       `gorm:"column:onboarding"`
	PricePlan            string     `gorm:"column:price_plan"`
	YearOfBirth          int        `gorm:"column:yob"`
	Company              string     `gorm:"column:company"`
	JobRole              string     `gorm:"column:job_role"`
	PcaId                string     `gorm:"column:pca_id"`
	WalletId             string     `gorm:"column:wallet_id"`
	RatingTitle          string     `gorm:"column:rating_title"`
	Price                int        `gorm:"column:price"`
	SocialType           string     `gorm:"column:social_type"`
	SocialId             string     `gorm:"column:social_id"`
	SocialUserName       string     `gorm:"column:social_user_name"`
	AccountType          string     `gorm:"column:account_type"`
	SkipUpdateInfo       bool       `gorm:"column:skip_update_info"`
	UpdateInfoDate       time.Time  `gorm:"column:update_info_date"`
	EmailLastUpdate      time.Time  `gorm:"column:email_last_update"`
	MsisdnLastUpdate     time.Time  `gorm:"column:msisdn_last_update"`
	PoBalance            float64    `gorm:"column:po_balance"`
	IsAgreeRelaunch      bool       `gorm:"column:is_agree_relaunce_info"`
	IsEligibleForChatbot bool       `gorm:"column:is_eligible_for_chatbot"`
}

func (User) TableName() string {
	return "users"
}
