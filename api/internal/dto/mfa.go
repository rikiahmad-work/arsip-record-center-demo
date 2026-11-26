package dto

import "time"

// MfaVerifyRequest adalah request payload untuk verifikasi OTP TOTP atau Recovery Code.
type MfaVerifyRequest struct {
	Challenge string `json:"challenge" binding:"required,len=64,hexadecimal"`
	Code      string `json:"code" binding:"required,min=6,max=80"`
}

// MfaFactorRequest digunakan untuk aksi sensitif (regenerasi recovery code, reset mandiri, reset admin).
type MfaFactorRequest struct {
	Code string `json:"code" binding:"required,min=6,max=80"`
}

// MfaStatusResponse mengembalikan status aktivasi MFA pengguna.
type MfaStatusResponse struct {
	Enabled                bool       `json:"enabled"`
	EnabledAt              *time.Time `json:"enabled_at,omitempty"`
	RecoveryCodesRemaining int        `json:"recovery_codes_remaining"`
}

// MfaChallengeData adalah payload tantangan MFA yang dikembalikan saat login atau inisiasi reset.
type MfaChallengeData struct {
	MfaRequired bool   `json:"mfa_required"`
	Challenge   string `json:"challenge"`
	Mode        string `json:"mode"` // "setup" atau "verify"
	Secret      string `json:"secret,omitempty"`
	OtpauthURI  string `json:"otpauth_uri,omitempty"`
}

// MfaVerifyData adalah payload yang dikembalikan setelah verifikasi MFA berhasil.
type MfaVerifyData struct {
	AccessToken   string        `json:"access_token"`
	Token         string        `json:"token"`
	TokenType     string        `json:"token_type"`
	ExpiresIn     int64         `json:"expires_in"`
	ExpiresAt     time.Time     `json:"expires_at"`
	RecoveryCodes []string      `json:"recovery_codes,omitempty"`
	User          *UserResponse `json:"user"`
}
