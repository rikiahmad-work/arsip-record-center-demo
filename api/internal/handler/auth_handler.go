package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// AuthHandler menangani HTTP request terkait autentikasi.
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler membuat instance AuthHandler.
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login menangani POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	req.ClientIP = c.ClientIP()
	req.UserAgent = c.GetHeader("User-Agent")

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server")
		return
	}

	msg := "Login berhasil"
	if resp != nil && resp.MfaRequired {
		msg = "Verifikasi MFA diperlukan"
	}
	helper.SuccessResponse(c, msg, resp)
}

// Logout menangani POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tokenJTI, _ := c.Get("token_jti")
	jti, ok := tokenJTI.(string)
	if !ok {
		helper.ErrorResponse(c, http.StatusBadRequest, "Token tidak valid")
		return
	}

	err := h.authService.Logout(c.Request.Context(), jti)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal logout")
		return
	}

	helper.SuccessResponse(c, "Logout berhasil", nil)
}

// Me menangani GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid")
		return
	}

	resp, err := h.authService.Me(c.Request.Context(), id)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan server")
		return
	}

	helper.SuccessResponse(c, "", resp)
}

// ChangePassword menangani PUT /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, ok := userID.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(c.Request.Context(), id, req)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengubah kata sandi")
		return
	}

	helper.SuccessResponse(c, "Kata sandi berhasil diperbarui", nil)
}
