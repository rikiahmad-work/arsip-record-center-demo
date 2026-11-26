package model

import (
	"time"
)

// BlacklistedToken menyimpan JWT token yang sudah di-logout.
type BlacklistedToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TokenJTI  string    `gorm:"column:token_jti;type:varchar(255);uniqueIndex;not null" json:"token_jti"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName override default table name ke "blacklisted_tokens".
func (BlacklistedToken) TableName() string {
	return "blacklisted_tokens"
}
