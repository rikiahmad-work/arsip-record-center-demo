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

// LaporanHandler menangani HTTP request untuk Dashboard, Laporan DAI, dan Cetak Label PDF.
type LaporanHandler struct {
	laporanSvc service.LaporanService
}

func NewLaporanHandler(laporanSvc service.LaporanService) *LaporanHandler {
	return &LaporanHandler{laporanSvc: laporanSvc}
}

// ─── 1. Dashboard (/api/dashboard/*) ────────────────────────────

func (h *LaporanHandler) DashboardStats(c *gin.Context) {
	stats, err := h.laporanSvc.GetDashboardStats(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil statistik dashboard")
		return
	}
	helper.SuccessResponse(c, "", stats)
}

func (h *LaporanHandler) DashboardChart(c *gin.Context) {
	tahun, _ := strconv.Atoi(c.Query("tahun"))
	chart, err := h.laporanSvc.GetDashboardChart(c.Request.Context(), tahun)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data grafik dashboard")
		return
	}
	helper.SuccessResponse(c, "", chart)
}

// ─── 2. Laporan & Rekap DAI (/api/laporan/*) ────────────────────

func (h *LaporanHandler) LaporanRekap(c *gin.Context) {
	var req dto.LaporanFilterRequest
	_ = c.ShouldBindQuery(&req)

	rekap, err := h.laporanSvc.GetLaporanRekap(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil rekapitulasi laporan")
		return
	}
	helper.SuccessResponse(c, "", rekap)
}

func (h *LaporanHandler) LaporanDAI(c *gin.Context) {
	var req dto.LaporanFilterRequest
	_ = c.ShouldBindQuery(&req)

	rows, err := h.laporanSvc.GetLaporanDAI(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil laporan Daftar Arsip Inaktif (DAI)")
		return
	}
	helper.SuccessResponse(c, "", rows)
}

func (h *LaporanHandler) LaporanExport(c *gin.Context) {
	var req dto.LaporanFilterRequest
	_ = c.ShouldBindQuery(&req)

	csvBytes, err := h.laporanSvc.ExportLaporanDAICSV(c.Request.Context(), req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengekspor laporan DAI")
		return
	}

	filename := fmt.Sprintf("rekap-dai-18-kolom-%s.csv", time.Now().Format("20060102-150405"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", csvBytes)
}

// ─── 3. Cetak Label PDF (/api/label/*) ──────────────────────────

func (h *LaporanHandler) LabelBox(c *gin.Context) {
	h.serveLabelPDF(c, "box", "label-boks.pdf")
}

func (h *LaporanHandler) LabelFolder(c *gin.Context) {
	h.serveLabelPDF(c, "folder", "label-folder.pdf")
}

func (h *LaporanHandler) LabelInaktif(c *gin.Context) {
	h.serveLabelPDF(c, "inaktif", "label-arsip-inaktif.pdf")
}

func (h *LaporanHandler) serveLabelPDF(c *gin.Context, jenis, filename string) {
	var req dto.LabelQueryRequest
	_ = c.ShouldBindQuery(&req)

	pdfBytes, err := h.laporanSvc.GenerateLabel(c.Request.Context(), jenis, req)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat PDF label")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	c.Header("X-Content-Type-Options", "nosniff")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
