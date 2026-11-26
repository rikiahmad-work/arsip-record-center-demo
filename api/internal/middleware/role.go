package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// RoleMiddleware adalah middleware untuk membatasi akses berdasarkan role & dynamic permission.
// Dipanggil setelah AuthMiddleware.
type RoleMiddleware struct {
	db *gorm.DB
}

// NewRoleMiddleware membuat instance RoleMiddleware dengan opsional koneksi database untuk RBAC dinamis.
func NewRoleMiddleware(db ...*gorm.DB) *RoleMiddleware {
	var d *gorm.DB
	if len(db) > 0 {
		d = db[0]
	}
	return &RoleMiddleware{db: d}
}

// RequireRole mengembalikan middleware Gin yang memeriksa apakah user memiliki
// salah satu role yang diizinkan.
func (m *RoleMiddleware) RequireRole(allowedRoles ...model.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau sesi telah berakhir")
			c.Abort()
			return
		}

		userRole := model.NormalizeRole(role.(string))

		for _, allowed := range allowedRoles {
			if userRole == model.NormalizeRole(string(allowed)) {
				c.Next()
				return
			}
		}

		helper.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki akses ke sumber daya ini")
		c.Abort()
	}
}

// RequirePermission mengembalikan middleware Gin yang memeriksa apakah peran pengguna memiliki
// hak akses (permission) dinamis tertentu dari matriks RBAC.
func (m *RoleMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau sesi telah berakhir")
			c.Abort()
			return
		}

		userRole := model.NormalizeRole(roleVal.(string))

		// Administrator & Super-Admin selalu memiliki akses penuh (industry standard bypass)
		if userRole == model.RoleAdministrator || userRole == "super-admin" {
			c.Next()
			return
		}

		if m.db == nil {
			c.Next()
			return
		}

		normKey := strings.ReplaceAll(string(userRole), "_", "-")
		var rp model.RolePermission
		if err := m.db.Where("`key` = ? OR `key` = ?", normKey, string(userRole)).First(&rp).Error; err == nil && rp.Permissions != "" {
			var perms []string
			if jsonErr := json.Unmarshal([]byte(rp.Permissions), &perms); jsonErr == nil {
				for _, p := range perms {
					if p == permission {
						c.Next()
						return
					}
				}
			}
		}

		helper.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki hak akses untuk tindakan ini ("+permission+")")
		c.Abort()
	}
}
