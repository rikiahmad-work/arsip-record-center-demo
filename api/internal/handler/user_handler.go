package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// UserHandler menangani HTTP request terkait manajemen user.
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler membuat instance UserHandler.
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Index menangani GET /api/users
func (h *UserHandler) Index(c *gin.Context) {
	var req dto.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	users, meta, err := h.userService.GetAll(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data user")
		return
	}

	helper.PaginatedResponse(c, users, meta)
}

// Show menangani GET /api/users/:id
func (h *UserHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.userService.GetByID(c.Request.Context(), uint(id))
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

// Store menangani POST /api/users
func (h *UserHandler) Store(c *gin.Context) {
	var req dto.StoreUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	createdBy, _ := c.Get("user_id")
	creatorID, ok := createdBy.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid")
		return
	}

	resp, err := h.userService.Create(c.Request.Context(), req, creatorID)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan user")
		return
	}

	helper.CreatedResponse(c, "User berhasil ditambahkan", resp)
}

// Update menangani PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	currentUserID, _ := c.Get("user_id")
	currID, _ := currentUserID.(uint)

	resp, err := h.userService.Update(c.Request.Context(), uint(id), req, currID)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate user")
		return
	}

	helper.SuccessResponse(c, "User berhasil diupdate", resp)
}

// Destroy menangani DELETE /api/users/:id
func (h *UserHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	currentUserID, _ := c.Get("user_id")
	userID, ok := currentUserID.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid")
		return
	}

	err = h.userService.Delete(c.Request.Context(), uint(id), userID)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus user")
		return
	}

	helper.SuccessResponse(c, "User berhasil dihapus", nil)
}

// Restore menangani PUT /api/users/:id/restore
func (h *UserHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.userService.Restore(c.Request.Context(), uint(id))
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memulihkan user")
		return
	}

	helper.SuccessResponse(c, "User berhasil dipulihkan", resp)
}

// Stats menangani GET /api/users/stats
func (h *UserHandler) Stats(c *gin.Context) {
	stats, err := h.userService.GetStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil statistik user")
		return
	}

	helper.SuccessResponse(c, "", stats)
}
