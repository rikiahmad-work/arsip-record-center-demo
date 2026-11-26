package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// OperasionalHandler menangani HTTP request untuk Peminjaman, Penyusutan, Dokumen Arsip, dan Audit Log.
type OperasionalHandler struct {
	opSvc service.OperasionalService
}

func NewOperasionalHandler(opSvc service.OperasionalService) *OperasionalHandler {
	return &OperasionalHandler{opSvc: opSvc}
}

func getCallerContext(c *gin.Context) (userID uint, role string, ip string) {
	if idVal, ok := c.Get("user_id"); ok {
		userID, _ = idVal.(uint)
	}
	if roleVal, ok := c.Get("user_role"); ok {
		role, _ = roleVal.(string)
	}
	ip = c.ClientIP()
	return userID, role, ip
}

// ─── 1. Peminjaman & Pengembalian (/api/peminjaman) ─────────────

func (h *OperasionalHandler) ListPeminjaman(c *gin.Context) {
	var req dto.PeminjamanListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.opSvc.ListPeminjaman(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data peminjaman")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}

func (h *OperasionalHandler) ShowPeminjaman(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.opSvc.GetPeminjamanByID(c.Request.Context(), uint(id))
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail peminjaman")
		return
	}
	helper.SuccessResponse(c, "", resp)
}

func (h *OperasionalHandler) StorePeminjaman(c *gin.Context) {
	var req dto.StorePeminjamanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	userID, _, ip := getCallerContext(c)
	resp, err := h.opSvc.CreatePeminjaman(c.Request.Context(), req, userID, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mencatat peminjaman")
		return
	}
	helper.CreatedResponse(c, "Peminjaman berhasil dicatat", resp)
}

func (h *OperasionalHandler) UpdatePeminjaman(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdatePeminjamanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	userID, _, ip := getCallerContext(c)
	resp, err := h.opSvc.UpdatePeminjaman(c.Request.Context(), uint(id), req, userID, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui peminjaman")
		return
	}
	helper.SuccessResponse(c, "Peminjaman berhasil diperbarui", resp)
}

func (h *OperasionalHandler) KembalikanPeminjaman(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.KembaliPeminjamanRequest
	_ = c.ShouldBindJSON(&req)

	userID, _, ip := getCallerContext(c)
	resp, err := h.opSvc.KembalikanPeminjaman(c.Request.Context(), uint(id), req, userID, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mencatat pengembalian arsip")
		return
	}
	helper.SuccessResponse(c, "Pengembalian arsip berhasil dicatat", resp)
}

func (h *OperasionalHandler) DestroyPeminjaman(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	userID, _, ip := getCallerContext(c)
	if err := h.opSvc.DeletePeminjaman(c.Request.Context(), uint(id), userID, ip); err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus data peminjaman")
		return
	}
	helper.SuccessResponse(c, "Peminjaman berhasil dihapus", nil)
}

// ─── 2. Penyusutan (/api/penyusutan) ────────────────────────────

func (h *OperasionalHandler) ListPenyusutan(c *gin.Context) {
	var req dto.PenyusutanListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.opSvc.ListPenyusutan(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data penyusutan")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}

func (h *OperasionalHandler) ShowPenyusutan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.opSvc.GetPenyusutanByID(c.Request.Context(), uint(id))
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil detail penyusutan")
		return
	}
	helper.SuccessResponse(c, "", resp)
}

func (h *OperasionalHandler) StorePenyusutan(c *gin.Context) {
	var req dto.StorePenyusutanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	userID, _, ip := getCallerContext(c)
	resp, err := h.opSvc.CreatePenyusutan(c.Request.Context(), req, userID, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengajukan penyusutan arsip")
		return
	}
	helper.CreatedResponse(c, "Pengajuan penyusutan berhasil dibuat", resp)
}

func (h *OperasionalHandler) ApprovePenyusutan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.DecisionPenyusutanRequest
	_ = c.ShouldBindJSON(&req)

	userID, role, ip := getCallerContext(c)
	resp, err := h.opSvc.ApprovePenyusutan(c.Request.Context(), uint(id), req, userID, role, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyetujui pengajuan penyusutan")
		return
	}
	helper.SuccessResponse(c, "Pengajuan penyusutan berhasil disetujui", resp)
}

func (h *OperasionalHandler) RejectPenyusutan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.DecisionPenyusutanRequest
	_ = c.ShouldBindJSON(&req)

	userID, role, ip := getCallerContext(c)
	resp, err := h.opSvc.RejectPenyusutan(c.Request.Context(), uint(id), req, userID, role, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menolak pengajuan penyusutan")
		return
	}
	helper.SuccessResponse(c, "Pengajuan penyusutan telah ditolak", resp)
}

func (h *OperasionalHandler) SelesaiPenyusutan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	userID, role, ip := getCallerContext(c)
	resp, err := h.opSvc.SelesaiPenyusutan(c.Request.Context(), uint(id), userID, role, ip)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyelesaikan pengajuan penyusutan")
		return
	}
	helper.SuccessResponse(c, "Penyusutan arsip ditandai selesai", resp)
}

// ─── 3. Dokumen Arsip (/api/arsip/:id/dokumen & /api/dokumen/:id) ─

func (h *OperasionalHandler) ListDokumen(c *gin.Context) {
	arsipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID arsip tidak valid")
		return
	}

	list, err := h.opSvc.ListDokumenByArsip(c.Request.Context(), uint(arsipID))
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar dokumen")
		return
	}
	helper.SuccessResponse(c, "", list)
}

func (h *OperasionalHandler) UploadDokumen(c *gin.Context) {
	arsipID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID arsip tidak valid")
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"file": "File dokumen wajib diunggah"})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca file dokumen")
		return
	}
	defer f.Close()

	var ketPtr *string
	if ket := c.PostForm("keterangan"); ket != "" {
		ketPtr = &ket
	}

	userID, _, ip := getCallerContext(c)
	resp, err := h.opSvc.UploadDokumen(
		c.Request.Context(),
		uint(arsipID),
		fileHeader.Filename,
		fileHeader.Size,
		f,
		ketPtr,
		userID,
		ip,
	)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengunggah dokumen arsip")
		return
	}
	helper.CreatedResponse(c, "Dokumen berhasil diunggah", resp)
}

func (h *OperasionalHandler) DownloadDokumen(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID dokumen tidak valid")
		return
	}

	userID, role, _ := getCallerContext(c)
	doc, err := h.opSvc.GetDokumenForDownload(c.Request.Context(), uint(id), role, userID)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil dokumen")
		return
	}

	data, err := os.ReadFile(doc.Path)
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "File fisik dokumen tidak ditemukan di storage")
		return
	}

	disposition := "attachment"
	if c.Query("disposition") == "inline" {
		disposition = "inline"
	}

	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=%q", disposition, doc.NamaFile))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, doc.MimeType, data)
}

func (h *OperasionalHandler) DestroyDokumen(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID dokumen tidak valid")
		return
	}

	userID, _, ip := getCallerContext(c)
	if err := h.opSvc.DeleteDokumen(c.Request.Context(), uint(id), userID, ip); err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus dokumen")
		return
	}
	helper.SuccessResponse(c, "Dokumen berhasil dihapus", nil)
}

// ─── 4. Audit Log (/api/audit-logs) ─────────────────────────────

func (h *OperasionalHandler) ListAuditLogs(c *gin.Context) {
	var req dto.AuditLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	list, meta, err := h.opSvc.ListAuditLog(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data audit log")
		return
	}
	helper.PaginatedResponse(c, list, meta)
}
