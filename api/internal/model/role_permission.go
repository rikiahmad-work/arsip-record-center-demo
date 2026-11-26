package model

import "time"

// RolePermission merepresentasikan peran pengguna (role) beserta daftar hak akses (permissions JSON array)
// seperti pada implementasi Spatie Laravel-Permission di project referensi.
type RolePermission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Key         string    `gorm:"column:key;type:varchar(50);uniqueIndex;not null" json:"key"`
	Name        string    `gorm:"column:name;type:varchar(100);uniqueIndex;not null" json:"name"`
	IsSystem    bool      `gorm:"column:is_system;default:false;not null" json:"is_system"`
	Permissions string    `gorm:"column:permissions;type:text;not null" json:"permissions"` // JSON array of permission keys
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
