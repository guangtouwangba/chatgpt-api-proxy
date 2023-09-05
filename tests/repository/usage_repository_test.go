package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/guangtouwangba/chatgpt-api-proxy/internal/db/model"
	"github.com/guangtouwangba/chatgpt-api-proxy/internal/db/repository"
)

var (
	mockDB   *gorm.DB
	mockSql  sqlmock.Sqlmock
	repo     *repository.GormOpenAIUsageRepository
)

func setup() {
	var err error
	mockDB, mockSql, err = sqlmock.New()
	if err != nil {
		panic("An error was not expected when opening a stub database connection")
	}

	repo = repository.NewGormOpenAIUsageRepository(mockDB)
}

func TestFindOneReturnsCorrectRecord(t *testing.T) {
	setup()

	mockSql.ExpectQuery("^SELECT (.+) FROM \"open_ai_usages\" WHERE (.+)$").
		WithArgs("validOpenAIID", "validIdentityID").
		WillReturnRows(sqlmock.NewRows([]string{"id", "openai_id", "identity_id", "usage", "model", "tokens"}).
			AddRow(1, "validOpenAIID", "validIdentityID", 100, "model", 50))

	usage, err := repo.FindOne("validOpenAIID", "validIdentityID")
	assert.NoError(t, err)
	assert.Equal(t, &model.OpenAIUsage{
		ID:         1,
		OpenAIID:   "validOpenAIID",
		IdentityID: "validIdentityID",
		Usage:      100,
		Model:      "model",
		Tokens:     50,
	}, usage)
}

func TestFindOneReturnsErrorForNonExistentRecord(t *testing.T) {
	setup()

	mockSql.ExpectQuery("^SELECT (.+) FROM \"open_ai_usages\" WHERE (.+)$").
		WithArgs("invalidOpenAIID", "invalidIdentityID").
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := repo.FindOne("invalidOpenAIID", "invalidIdentityID")
	assert.Error(t, err)
}

func TestFindOneReturnsErrorForDBConnectionProblem(t *testing.T) {
	setup()

	mockSql.ExpectQuery("^SELECT (.+) FROM \"open_ai_usages\" WHERE (.+)$").
		WithArgs("validOpenAIID", "validIdentityID").
		WillReturnError(gorm.ErrInvalidSQL)

	_, err := repo.FindOne("validOpenAIID", "validIdentityID")
	assert.Error(t, err)
}
