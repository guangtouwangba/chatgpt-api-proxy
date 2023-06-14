package model

type OpenAIUsage struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	OpenAIID   string `gorm:"column:openai_id" json:"openai_id"`
	IdentityID string `gorm:"column:identity_id" json:"identity_id"`
	Usage      int64  `gorm:"column:usage" json:"usage"`
	Model      string `gorm:"column:model" json:"model"`
	Tokens     int64  `gorm:"column:tokens" json:"tokens"`
}

type OpenAIUsageRepository interface {
	FindOne(openAIID string, identityID string, model string) (*OpenAIUsage, error)
	Create(openAIID string, identityID string, model string, usage int64, tokens int64) error
	Update(openAIID string, identityID string, model string, usage int64, tokens int64) error
}
