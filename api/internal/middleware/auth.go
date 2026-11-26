package middleware

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// JWTClaims adalah struktur claims yang akan dimasukkan dalam token JWT.
type JWTClaims struct {
	JTI           string `json:"jti"`
	Sub           uint   `json:"sub"`
	Username      string `json:"username"`
	NIP           string `json:"nip"`
	Role          string `json:"role"`
	MfaVerified   bool   `json:"mfa_verified"`
	MfaVerifiedAt int64  `json:"mfa_verified_at,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware menyediakan middleware untuk validasi JWT token.
type AuthMiddleware struct {
	cfg    *config.Config
	db     *gorm.DB
	secret []byte
}

// NewAuthMiddleware membuat instance AuthMiddleware.
func NewAuthMiddleware(cfg *config.Config, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{
		cfg:    cfg,
		db:     db,
		secret: []byte(cfg.JWT.Secret),
	}
}

// RequireAuth adalah middleware untuk memvalidasi JWT token.
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c)
		if tokenString == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan")
			c.Abort()
			return
		}

		claims, err := m.validateToken(tokenString)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau sesi telah berakhir")
			c.Abort()
			return
		}

		// Cek apakah token sudah di-blacklist (jika db tersedia)
		if m.db != nil {
			blacklisted, err := m.isTokenBlacklisted(claims.JTI)
			if err != nil || blacklisted {
				helper.ErrorResponse(c, http.StatusUnauthorized, "Sesi telah berakhir, silakan login kembali")
				c.Abort()
				return
			}

			// Cek apakah user masih aktif di database
			var user model.User
			if err := m.db.First(&user, claims.Sub).Error; err != nil {
				helper.ErrorResponse(c, http.StatusUnauthorized, "Akun tidak ditemukan")
				c.Abort()
				return
			}
			if !user.IsActive {
				helper.ErrorResponse(c, http.StatusUnauthorized, "Akun Anda telah dinonaktifkan. Silakan hubungi administrator.")
				c.Abort()
				return
			}
		}

		// Validasi MFA jika diwajibkan oleh konfigurasi sistem
		if m.cfg.MFA.Required && !claims.MfaVerified {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Verifikasi MFA diperlukan.",
				"data":    gin.H{"mfa_required": true},
			})
			c.Abort()
			return
		}

		// Set user info ke context untuk digunakan handler selanjutnya
		c.Set("user_id", claims.Sub)
		c.Set("user_username", claims.Username)
		c.Set("user_nip", claims.NIP)
		c.Set("user_role", claims.Role)
		c.Set("token_jti", claims.JTI)
		if claims.ExpiresAt != nil {
			c.Set("token_exp", claims.ExpiresAt.Time.Unix())
		}

		c.Next()
	}
}

// GenerateToken membuat JWT token baru dengan JTI berformat UUID v4 kriptografis.
func (m *AuthMiddleware) GenerateToken(userID uint, username, nip, role string, rememberMe bool) (string, time.Time, error) {
	return m.GenerateTokenWithMFA(userID, username, nip, role, rememberMe, !m.cfg.MFA.Required)
}

// GenerateTokenWithMFA membuat JWT token baru dengan eksplisit flag status verifikasi MFA.
func (m *AuthMiddleware) GenerateTokenWithMFA(userID uint, username, nip, role string, rememberMe bool, mfaVerified bool) (string, time.Time, error) {
	expireHours := m.cfg.JWT.AccessExpireHours
	if rememberMe {
		expireHours = m.cfg.JWT.RefreshExpireHours
	}
	if expireHours <= 0 {
		expireHours = 7
	}

	expiresAt := time.Now().Add(time.Duration(expireHours) * time.Hour)
	jti := generateUUIDv4()

	var mfaVerifiedAt int64 = 0
	if mfaVerified {
		mfaVerifiedAt = time.Now().Unix()
	}

	claims := &JWTClaims{
		JTI:           jti,
		Sub:           userID,
		Username:      username,
		NIP:           nip,
		Role:          role,
		MfaVerified:   mfaVerified,
		MfaVerifiedAt: mfaVerifiedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// BlacklistToken menambahkan JTI token ke tabel blacklisted_tokens.
func (m *AuthMiddleware) BlacklistToken(jti string, expiresAt time.Time) error {
	if m.db == nil || jti == "" {
		return nil
	}
	return m.db.Exec(
		"INSERT IGNORE INTO blacklisted_tokens (token_jti, expires_at, created_at) VALUES (?, ?, NOW())",
		jti, expiresAt,
	).Error
}

// validateToken memverifikasi signature, expiry, dan secara eksplisit memeriksa algoritma HMAC (HS256)
// untuk mencegah serangan JWT algorithm confusion (alg: none / RS256->HS256).
func (m *AuthMiddleware) validateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok || t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func (m *AuthMiddleware) isTokenBlacklisted(jti string) (bool, error) {
	var count int64
	err := m.db.Table("blacklisted_tokens").
		Where("token_jti = ? AND expires_at > NOW()", jti).
		Count(&count).Error
	return count > 0, err
}

func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}

// generateUUIDv4 menghasilkan UUID v4 standar RFC 4122 menggunakan crypto/rand.
func generateUUIDv4() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
