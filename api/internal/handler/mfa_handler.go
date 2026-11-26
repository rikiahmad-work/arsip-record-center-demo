package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// MfaHandler menangani HTTP request terkait Multi-Factor Authentication.
type MfaHandler struct {
	mfaService service.MfaService
	userRepo   repository.UserRepository
}

// NewMfaHandler membuat instance baru MfaHandler.
func NewMfaHandler(mfaService service.MfaService, userRepo repository.UserRepository) *MfaHandler {
	return &MfaHandler{
		mfaService: mfaService,
		userRepo:   userRepo,
	}
}

// Verify menangani POST /api/auth/mfa/verify
func (h *MfaHandler) Verify(c *gin.Context) {
	var req dto.MfaVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": "Challenge (64 karakter hex) dan kode verifikasi wajib diisi",
		})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	data, err := h.mfaService.Verify(c.Request.Context(), req, clientIP, userAgent)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses verifikasi MFA")
		return
	}

	helper.SuccessResponse(c, "MFA berhasil diverifikasi", data)
}

// Status menangani GET /api/auth/mfa/status
func (h *MfaHandler) Status(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid")
		return
	}
	userID := val.(uint)

	status, err := h.mfaService.GetStatus(c.Request.Context(), userID)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memeriksa status MFA")
		return
	}

	helper.SuccessResponse(c, "", status)
}

// ResetRecoveryCodes menangani POST /api/auth/mfa/recovery-codes
func (h *MfaHandler) ResetRecoveryCodes(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid")
		return
	}
	userID := val.(uint)

	var req dto.MfaFactorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"code": "Kode OTP atau recovery code aktif wajib diisi",
		})
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	codes, err := h.mfaService.ResetRecoveryCodes(c.Request.Context(), userID, req.Code, clientIP, userAgent)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui recovery code")
		return
	}

	helper.SuccessResponse(c, "Recovery code diperbarui", gin.H{
		"recovery_codes": codes,
	})
}

// Reset menangani POST /api/auth/mfa/reset (reset mandiri)
func (h *MfaHandler) Reset(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid")
		return
	}
	userID := val.(uint)

	var req dto.MfaFactorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"code": "Kode OTP atau recovery code aktif wajib diisi",
		})
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil || user == nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	challengeData, err := h.mfaService.ResetSelf(c.Request.Context(), user, req.Code, clientIP, userAgent)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mereset MFA")
		return
	}

	helper.SuccessResponse(c, "MFA direset. Selesaikan setup ulang.", challengeData)
}

// ResetUser menangani POST /api/users/:id/mfa/reset (reset oleh Administrator)
func (h *MfaHandler) ResetUser(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid")
		return
	}
	adminUserID := val.(uint)

	idParam := c.Param("id")
	targetID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID pengguna tidak valid")
		return
	}

	var req dto.MfaFactorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"code": "Kode OTP Administrator wajib diisi untuk verifikasi",
		})
		return
	}

	adminUser, err := h.userRepo.FindByID(c.Request.Context(), adminUserID)
	if err != nil || adminUser == nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Pengguna admin tidak ditemukan")
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.mfaService.ResetUserByAdmin(c.Request.Context(), adminUser, uint(targetID), req.Code, clientIP, userAgent); err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mereset MFA akun")
		return
	}

	helper.SuccessResponse(c, "MFA akun berhasil direset", nil)
}
