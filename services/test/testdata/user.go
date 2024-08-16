package testdata_internal

import (
	"database/sql/driver"
	"reflect"
	"strings"
	"time"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/models"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/common"
	"github.com/bxcodec/faker/v3"
)

func NewUserModels() *models.User {
	return &models.User{
		Id:                   common.NewID().String(),
		Type:                 0,
		Name:                 faker.Name(),
		Email:                faker.Email(),
		Msisdn:               "6285846132417",
		EmailVerified:        true,
		MsisdnVerified:       true,
		PhoneNumberTelkomsel: false,
		ProfileVerified:      true,
		Roles:                "SURVEY",
		Status:               3,
		WhatsappNumber:       false,
		SubscribeNewsletter:  false,
		CreatedAt:            time.Time{},
		UpdatedAt:            time.Time{},
		DeletedAt:            &time.Time{},
		PricePlan:            "tsel",
		YearOfBirth:          1996,
		IsAgreeRelaunch:      false,
		IsEligibleForChatbot: true,
	}
}

func StructToColumns(data interface{}) (colums []string) {
	dataColums := reflect.TypeOf(data)
	for i := 0; i < dataColums.NumField(); i++ {
		field := dataColums.Field(i)
		columnName := field.Tag.Get("gorm")
		columnName = strings.TrimPrefix(columnName, "column:")
		colums = append(colums, columnName)
	}
	return colums
}

func StructToValue(data interface{}) (value []driver.Value) {
	dataValue := reflect.ValueOf(data)
	for i := 0; i < dataValue.Type().NumField(); i++ {
		dataV := dataValue.Field(i)
		value = append(value, dataV.Interface())
	}
	return value
}
