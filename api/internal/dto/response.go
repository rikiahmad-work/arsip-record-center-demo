package dto

import "time"

// LoginResponse dikembalikan saat login berhasil (bisa berupa token langsung atau MFA challenge).
type LoginResponse struct {
	MfaRequired bool          `json:"mfa_required,omitempty"`
	Challenge   string        `json:"challenge,omitempty"`
	Mode        string        `json:"mode,omitempty"`
	Secret      string        `json:"secret,omitempty"`
	OtpauthURI  string        `json:"otpauth_uri,omitempty"`
	User        *UserResponse `json:"user,omitempty"`
	Token       string        `json:"token,omitempty"`
	ExpiresAt   *time.Time    `json:"expires_at,omitempty"`
}

// UserResponse merepresentasikan data user yang dikembalikan ke client.
type UserResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Name        string     `json:"name"`
	NIP         string     `json:"nip"`
	Email       string     `json:"email"`
	Role        string     `json:"role"`
	RoleLabel   string     `json:"role_label"`
	Permissions []string   `json:"permissions"`
	UnitKerja   *string    `json:"unit_kerja"`
	Jabatan     *string    `json:"jabatan"`
	IsSSO       bool       `json:"is_sso"`
	IsActive    bool       `json:"is_active"`
	MfaEnabled  bool       `json:"mfa_enabled"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedBy   *uint      `json:"created_by"`
}

// PaginatedResponse digunakan untuk response dengan pagination.
type PaginatedResponse struct {
	Data        interface{} `json:"data"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	PerPage     int         `json:"per_page"`
	Total       int64       `json:"total"`
	From        int         `json:"from"`
	To          int         `json:"to"`
}

// UserStatsResponse digunakan untuk statistik user (4 role kanonik + alias kompatibilitas).
type UserStatsResponse struct {
	Total         int64 `json:"total"`
	Aktif         int64 `json:"aktif"`
	Nonaktif      int64 `json:"nonaktif"`
	Administrator int64 `json:"administrator"`
	Operator      int64 `json:"operator"`
	Viewer        int64 `json:"viewer"`
	Approver      int64 `json:"approver"`
	SuperAdmin    int64 `json:"super_admin"`
	AdminArsip    int64 `json:"admin_arsip"`
}
