package repository

import (
	"chatgpt-api-proxy/internal/db/model"

	"gorm.io/gorm"
)

type GormOpenAIUsageRepository struct {
	db *gorm.DB
}

func (g *GormOpenAIUsageRepository) FindOne(openAIID string, identityID string, mod string) (*model.OpenAIUsage, error) {
	var usage model.OpenAIUsage
	if err := g.db.Where("openai_id = ? AND identity_id = ? AND model = ?", openAIID, identityID, mod).First(&usage).Error; err != nil {
		return nil, err
	}
	return &usage, nil
}

func (g *GormOpenAIUsageRepository) Create(openAIID string, identityID string, mod string, usage int64, tokens int64) error {
	return g.db.Create(&model.OpenAIUsage{
		OpenAIID:   openAIID,
		IdentityID: identityID,
		Model:      mod,
		Usage:      usage,
		Tokens:     tokens,
	}).Error
}

func (g *GormOpenAIUsageRepository) Update(openAIID string, identityID string, mod string, usage int64, tokens int64) error {
	return g.db.Model(&model.OpenAIUsage{}).Where("openai_id = ? AND identity_id = ? AND model = ?",
		openAIID, identityID, mod).Update("usage", usage).Update("tokens", tokens).Error
}

func NewGormOpenAIUsageRepository(db *gorm.DB) *GormOpenAIUsageRepository {
	return &GormOpenAIUsageRepository{
		db: db,
	}
}
