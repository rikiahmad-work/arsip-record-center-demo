package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// MasterHandler menangani HTTP request untuk master data (/api/master/*).
type MasterHandler struct {
	masterSvc service.MasterService
}

func NewMasterHandler(masterSvc service.MasterService) *MasterHandler {
	return &MasterHandler{masterSvc: masterSvc}
}

// ─── Kode Klasifikasi ───────────────────────────────────────────

func (h *MasterHandler) ListKodeKlasifikasi(c *gin.Context) {
	var req dto.KodeKlasifikasiListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.masterSvc.ListKodeKlasifikasi(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data kode klasifikasi")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}

func (h *MasterHandler) StoreKodeKlasifikasi(c *gin.Context) {
	var req dto.StoreKodeKlasifikasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.CreateKodeKlasifikasi(c.Request.Context(), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan kode klasifikasi")
		return
	}
	helper.CreatedResponse(c, "Kode klasifikasi berhasil ditambahkan", resp)
}

func (h *MasterHandler) UpdateKodeKlasifikasi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdateKodeKlasifikasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.UpdateKodeKlasifikasi(c.Request.Context(), uint(id), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui kode klasifikasi")
		return
	}
	helper.SuccessResponse(c, "Kode klasifikasi berhasil diperbarui", resp)
}

func (h *MasterHandler) DestroyKodeKlasifikasi(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	if err := h.masterSvc.DeleteKodeKlasifikasi(c.Request.Context(), uint(id)); err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus kode klasifikasi")
		return
	}
	helper.SuccessResponse(c, "Kode klasifikasi berhasil dihapus", nil)
}

// ─── Lokasi Simpan ──────────────────────────────────────────────

func (h *MasterHandler) ListLokasiSimpan(c *gin.Context) {
	var req dto.LokasiSimpanListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.masterSvc.ListLokasiSimpan(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data lokasi simpan")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}

func (h *MasterHandler) StoreLokasiSimpan(c *gin.Context) {
	var req dto.StoreLokasiSimpanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.CreateLokasiSimpan(c.Request.Context(), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan lokasi simpan")
		return
	}
	helper.CreatedResponse(c, "Lokasi simpan berhasil ditambahkan", resp)
}

func (h *MasterHandler) UpdateLokasiSimpan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdateLokasiSimpanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.UpdateLokasiSimpan(c.Request.Context(), uint(id), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui lokasi simpan")
		return
	}
	helper.SuccessResponse(c, "Lokasi simpan berhasil diperbarui", resp)
}

func (h *MasterHandler) DestroyLokasiSimpan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	if err := h.masterSvc.DeleteLokasiSimpan(c.Request.Context(), uint(id)); err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus lokasi simpan")
		return
	}
	helper.SuccessResponse(c, "Lokasi simpan berhasil dihapus", nil)
}

// ─── Unit Kerja ─────────────────────────────────────────────────

func (h *MasterHandler) ListUnitKerja(c *gin.Context) {
	var req dto.UnitKerjaListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.masterSvc.ListUnitKerja(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data unit kerja")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}

func (h *MasterHandler) StoreUnitKerja(c *gin.Context) {
	var req dto.StoreUnitKerjaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.CreateUnitKerja(c.Request.Context(), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan unit kerja")
		return
	}
	helper.CreatedResponse(c, "Unit kerja berhasil ditambahkan", resp)
}

func (h *MasterHandler) UpdateUnitKerja(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdateUnitKerjaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	resp, err := h.masterSvc.UpdateUnitKerja(c.Request.Context(), uint(id), req)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui unit kerja")
		return
	}
	helper.SuccessResponse(c, "Unit kerja berhasil diperbarui", resp)
}

func (h *MasterHandler) DestroyUnitKerja(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	if err := h.masterSvc.DeleteUnitKerja(c.Request.Context(), uint(id)); err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus unit kerja")
		return
	}
	helper.SuccessResponse(c, "Unit kerja berhasil dihapus", nil)
}

// ─── Konfigurasi Sistem ─────────────────────────────────────────

func (h *MasterHandler) GetKonfigurasi(c *gin.Context) {
	cfg, err := h.masterSvc.GetKonfigurasi(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil konfigurasi sistem")
		return
	}
	helper.SuccessResponse(c, "", cfg)
}

func (h *MasterHandler) UpdateKonfigurasi(c *gin.Context) {
	var req dto.UpdateKonfigurasiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(uint)

	cfg, err := h.masterSvc.UpdateKonfigurasi(c.Request.Context(), req, userID)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui konfigurasi sistem")
		return
	}
	helper.SuccessResponse(c, "Konfigurasi sistem berhasil diperbarui", cfg)
}
