package model

import (
	"time"
)

// UserMfa menyimpan konfigurasi MFA TOTP dan hash recovery codes per pengguna.
type UserMfa struct {
	UserID          uint       `gorm:"primaryKey" json:"user_id"`
	Secret          string     `gorm:"type:text" json:"-"`
	PendingSecret   string     `gorm:"type:text" json:"-"`
	RecoveryCodes   *string    `gorm:"type:json" json:"-"` // JSON array string of SHA-256 hashes
	LastTotpCounter uint64     `gorm:"column:last_totp_counter" json:"last_totp_counter"`
	FailedAttempts  uint8      `gorm:"default:0" json:"failed_attempts"`
	LockedUntil     *time.Time `json:"locked_until"`
	EnabledAt       *time.Time `json:"enabled_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (UserMfa) TableName() string {
	return "user_mfas"
}

// MfaChallenge merepresentasikan sesi tantangan verifikasi MFA (5 menit) pra-login atau saat setup.
type MfaChallenge struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint       `gorm:"index;not null" json:"user_id"`
	ChallengeHash  string     `gorm:"type:char(64);uniqueIndex;not null" json:"-"`
	Purpose        string     `gorm:"type:varchar(24);not null" json:"purpose"` // "setup" atau "login"
	Attempts       uint8      `gorm:"default:0" json:"attempts"`
	IPAddress      string     `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgentHash  string     `gorm:"type:char(64)" json:"user_agent_hash"`
	ExpiresAt      time.Time  `json:"expires_at"`
	TokenExpiresAt time.Time  `json:"token_expires_at"`
	ConsumedAt     *time.Time `json:"consumed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (MfaChallenge) TableName() string {
	return "mfa_challenges"
}

// MfaAuditLog mencatat jejak audit operasi MFA.
type MfaAuditLog struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ActorUserID  *uint     `gorm:"index" json:"actor_user_id"`
	TargetUserID uint      `gorm:"index;not null" json:"target_user_id"`
	Action       string    `gorm:"type:varchar(40);not null" json:"action"`
	IPAddress    string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent    string    `gorm:"type:varchar(255)" json:"user_agent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (MfaAuditLog) TableName() string {
	return "mfa_audit_logs"
}
