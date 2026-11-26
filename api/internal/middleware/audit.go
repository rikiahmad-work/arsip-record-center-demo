package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
)

// AuditMiddleware mencatat secara otomatis request mutasi (POST, PUT, DELETE) yang berhasil (HTTP 2xx)
// ke tabel audit_log untuk entitas yang belum dicatat manual di service layer (seperti auth/login, users, arsip, master).
func AuditMiddleware(opRepo repository.OperasionalRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if opRepo == nil {
			c.Next()
			return
		}

		method := c.Request.Method
		if method != http.MethodPost && method != http.MethodPut && method != http.MethodDelete {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		// Skip entitas yang sudah memiliki audit detail (data_lama & data_baru) langsung di OperasionalService
		if strings.HasPrefix(path, "/api/peminjaman") ||
			strings.HasPrefix(path, "/api/penyusutan") ||
			strings.HasPrefix(path, "/api/dokumen") ||
			strings.Contains(path, "/dokumen") {
			c.Next()
			return
		}

		var bodyPreview string
		if c.Request.Body != nil && strings.Contains(c.GetHeader("Content-Type"), "application/json") && !strings.HasSuffix(path, "/login") {
			raw, _ := io.ReadAll(io.LimitReader(c.Request.Body, 4096))
			c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))
			if len(raw) > 0 {
				bodyPreview = string(raw)
			}
		}

		c.Next()

		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			return
		}

		aksi := "create"
		switch {
		case strings.HasSuffix(path, "/login"):
			aksi = "login"
		case strings.HasSuffix(path, "/logout"):
			aksi = "logout"
		case method == http.MethodPut:
			aksi = "update"
		case method == http.MethodDelete:
			aksi = "delete"
		}

		entitas := "system"
		switch {
		case strings.HasPrefix(path, "/api/auth"):
			entitas = "auth"
		case strings.HasPrefix(path, "/api/users"):
			entitas = "user"
		case strings.HasPrefix(path, "/api/arsip"):
			entitas = "arsip"
		case strings.HasPrefix(path, "/api/master"):
			entitas = "master"
		}

		var userIDPtr *uint
		if idVal, ok := c.Get("user_id"); ok {
			if uid, okCast := idVal.(uint); okCast && uid > 0 {
				userIDPtr = &uid
			}
		}

		var dataBaruPtr *string
		if bodyPreview != "" {
			dataBaruPtr = &bodyPreview
		}

		ip := c.ClientIP()
		_ = opRepo.CreateAuditLog(c.Request.Context(), &model.AuditLog{
			UserID:    userIDPtr,
			Aksi:      aksi,
			Entitas:   entitas,
			DataBaru:  dataBaruPtr,
			IPAddress: &ip,
		})
	}
}
