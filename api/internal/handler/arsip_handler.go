package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// ArsipHandler menangani HTTP request terkait arsip.
type ArsipHandler struct {
	arsipService service.ArsipService
}

// NewArsipHandler membuat instance ArsipHandler.
func NewArsipHandler(arsipService service.ArsipService) *ArsipHandler {
	return &ArsipHandler{arsipService: arsipService}
}

// Index menangani GET /api/arsip
func (h *ArsipHandler) Index(c *gin.Context) {
	var req dto.ArsipListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	arsipList, meta, err := h.arsipService.GetAll(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data arsip")
		return
	}

	helper.PaginatedResponse(c, arsipList, meta)
}

// Show menangani GET /api/arsip/:id
func (h *ArsipHandler) Show(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.arsipService.GetByID(c.Request.Context(), uint(id))
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

// Store menangani POST /api/arsip
func (h *ArsipHandler) Store(c *gin.Context) {
	var req dto.StoreArsipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	petugasInputID, _ := c.Get("user_id")
	userID, ok := petugasInputID.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid")
		return
	}
	roleVal, _ := c.Get("user_role")
	callerRole, _ := roleVal.(string)

	resp, err := h.arsipService.Create(c.Request.Context(), req, userID, callerRole)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan arsip")
		return
	}

	helper.CreatedResponse(c, "Arsip berhasil ditambahkan", resp)
}

// Update menangani PUT /api/arsip/:id
func (h *ArsipHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	var req dto.UpdateArsipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{
			"validation": err.Error(),
		})
		return
	}

	roleVal, _ := c.Get("user_role")
	callerRole, _ := roleVal.(string)

	resp, err := h.arsipService.Update(c.Request.Context(), uint(id), req, callerRole)
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengupdate arsip")
		return
	}

	helper.SuccessResponse(c, "Arsip berhasil diupdate", resp)
}

// Destroy menangani DELETE /api/arsip/:id
func (h *ArsipHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	err = h.arsipService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus arsip")
		return
	}

	helper.SuccessResponse(c, "Arsip berhasil dihapus", nil)
}

// Restore menangani PUT /api/arsip/:id/restore
func (h *ArsipHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID tidak valid")
		return
	}

	resp, err := h.arsipService.Restore(c.Request.Context(), uint(id))
	if err != nil {
		appErr := helper.GetAppError(err)
		if appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal memulihkan arsip")
		return
	}

	helper.SuccessResponse(c, "Arsip berhasil dipulihkan", resp)
}

// Stats menangani GET /api/arsip/stats
func (h *ArsipHandler) Stats(c *gin.Context) {
	stats, err := h.arsipService.GetStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil statistik arsip")
		return
	}

	helper.SuccessResponse(c, "", stats)
}

// DaftarUNOR menangani GET /api/arsip/unor
func (h *ArsipHandler) DaftarUNOR(c *gin.Context) {
	daftar := h.arsipService.GetDaftarUNOR()
	helper.SuccessResponse(c, "", daftar)
}

// Export menangani GET /api/arsip/export (mendukung format=xlsx untuk Excel 18 kolom DAI asli maupun format=csv)
func (h *ArsipHandler) Export(c *gin.Context) {
	var req dto.ArsipListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	format := c.Query("format")
	if format == "xlsx" || format == "excel" {
		data, err := h.arsipService.ExportExcel(c.Request.Context(), req)
		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengekspor data arsip ke Excel")
			return
		}
		filename := fmt.Sprintf("DAFTAR-ARSIP-INAKTIF-%s.xlsx", time.Now().Format("20060102-150405"))
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		c.Header("Content-Length", strconv.Itoa(len(data)))
		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
		return
	}

	data, err := h.arsipService.ExportCSV(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengekspor data arsip")
		return
	}

	filename := fmt.Sprintf("daftar-arsip-inaktif-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", data)
}

// Import menangani POST /api/arsip/import (multipart file CSV)
func (h *ArsipHandler) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"file": "File import wajib diunggah"})
		return
	}

	f, err := file.Open()
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Gagal membaca file import")
		return
	}
	defer f.Close()

	petugasInputID, _ := c.Get("user_id")
	userID, _ := petugasInputID.(uint)
	roleVal, _ := c.Get("user_role")
	callerRole, _ := roleVal.(string)

	res, err := h.arsipService.ImportCSV(c.Request.Context(), f, userID, callerRole)
	if err != nil {
		if appErr := helper.GetAppError(err); appErr != nil {
			helper.AppErrorResponse(c, appErr)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengimpor data arsip")
		return
	}

	helper.SuccessResponse(c, fmt.Sprintf("Berhasil mengimpor %d dari %d baris arsip", res.ImportedRows, res.TotalRows), res)
}
