package psql_user_test

import (
	"context"
	"regexp"
	"testing"

	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/internal/repository/psql/psql_user"
	testdata_internal "cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-chat/services/test/testdata"
	"cicd-gitlab-ee.telkomsel.co.id/telkomsel/t-survey/crox-tsurvey/go-pkg/mocks"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_userRepository_GetUsersByID(t *testing.T) {
	type args struct {
		ctx    context.Context
		userId string
	}

	var (
		mdl     = testdata_internal.NewUserModels()
		columns = testdata_internal.StructToColumns(*mdl)
		values  = testdata_internal.StructToValue(*mdl)
	)

	query := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT 1`
	tests := []struct {
		name    string
		args    args
		want    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Find User Failed Test",
			args: args{
				ctx:    context.Background(),
				userId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
			},
			want: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(gorm.ErrInvalidData)
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userId: "6b31fdc7-b5ad-44a7-b216-d291fada3a21",
			},
			want: func(s sqlmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(sqlmock.NewRows(columns).AddRow(values...))
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mockSQL, err := mocks.MockDatabaseGorm(t)
			if err != nil {
				t.Fatalf("an error '%v' was not expected when opening a stub database connection", err)
			}
			if tt.want != nil {
				tt.want(mockSQL)
			}
			r := psql_user.UserRepositoryFactory{
				Db: mockDB,
			}
			a, _ := r.Create()
			_, err = a.GetUsersByID(tt.args.ctx, tt.args.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
