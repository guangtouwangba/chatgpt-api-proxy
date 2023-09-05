package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/guangtouwangba/chatgpt-api-proxy/internal/db/model"
	"github.com/guangtouwangba/chatgpt-api-proxy/internal/db/repository"
)

func TestFindOne(t *testing.T) {
	// Setup
	repo := repository.NewGormOpenAIUsageRepository()
	openAIID := "testOpenAIID"
	identityID := "testIdentityID"
	usage := &model.OpenAIUsage{
		OpenAIID:   openAIID,
		IdentityID: identityID,
		Model:      "testModel",
		Usage:      100,
		Tokens:     50,
	}
	err := repo.Create(openAIID, identityID, usage.Model, usage.Usage, usage.Tokens)
	assert.Nil(t, err)

	// Test successful retrieval
	retrievedUsage, err := repo.FindOne(openAIID, identityID)
	assert.Nil(t, err)
	assert.Equal(t, usage, retrievedUsage)

	// Test error scenario
	_, err = repo.FindOne("nonExistentOpenAIID", "nonExistentIdentityID")
	assert.NotNil(t, err)
}
