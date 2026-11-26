package model

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan entitas pengguna dalam sistem.
type User struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	NIP         string         `gorm:"column:nip;type:varchar(30);uniqueIndex;not null" json:"nip"`
	Email       string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password    string         `gorm:"type:varchar(255);not null" json:"-"`
	Role        string         `gorm:"type:enum('administrator','operator','viewer','approver');not null;default:'operator'" json:"role"`
	UnitKerja   *string        `gorm:"type:varchar(255)" json:"unit_kerja"`
	Jabatan     *string        `gorm:"type:varchar(255)" json:"jabatan"`
	IsSSO       bool           `gorm:"default:false;not null" json:"is_sso"`
	IsActive    bool           `gorm:"default:true;not null" json:"is_active"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedBy   *uint          `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName override default table name ke "users".
func (User) TableName() string {
	return "users"
}
