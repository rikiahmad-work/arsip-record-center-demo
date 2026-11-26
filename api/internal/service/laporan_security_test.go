package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/handler"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

func TestPhase3_DAI18Columns_And_LabelPDFGeneration(t *testing.T) {
	// 1. Pastikan struct DAIRowResponse memiliki tepat 18 kolom sesuai Format Daftar Arsip Inaktif (DAI)
	rt := reflect.TypeOf(dto.DAIRowResponse{})
	if rt.NumField() != 18 {
		t.Fatalf("expected DAIRowResponse to have exactly 18 columns, got %d", rt.NumField())
	}

	musnah := "Musnah"
	defNo := "DEF-001"
	rows := service.FormatDAIRows([]model.Arsip{
		{
			ID:                   1,
			NoItemArsip:          "1",
			KodeKlasifikasi:      "DT.04.01",
			UraianInformasiArsip: "Nota Dinas Pengelolaan Arsip",
			JenisArsip:           "Konvensional",
			KurunWaktu:           "2026",
			TingkatPerkembangan:  "Asli",
			Jumlah:               "1 Berkas",
			Keterangan:           "Baik",
			NomorDefinitif:       &defNo,
			LokasiSimpan:         "Record Center",
			NoBoks:               "B-01",
			NoRak:                "R-01",
			JangkaSimpanInaktif:  "4 Tahun",
			NasibAkhir:           &musnah,
			SKKAD:                "Biasa",
			KategoriArsip:        "Inaktif",
			PenciptaArsip:        "Biro Umum & SDM",
		},
	})
	if len(rows) != 1 || rows[0].No != 1 || rows[0].NasibAkhir != "Musnah" || rows[0].NomorDefinitif != "DEF-001" {
		t.Fatalf("unexpected DAI row formatting: %+v", rows)
	}

	// 2. Pastikan generator PDF Label menghasilkan dokumen PDF valid yang memuat Unit Pencipta, Tahun, Nomor Boks, dan Rentang Folder
	pdfBytes := helper.GenerateLabelPDF(helper.LabelEntry{
		JenisLabel: "LABEL BOKS ARSIP",
		UnitKerja:  "Biro Data dan Informasi",
		Kode:       "DT.04.01",
		Tahun:      "2026",
		NoRak:      "R-02",
		NoBoks:     "B-99",
		NoFolder:   "F-01 s/d F-10",
	})
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-1.4")) {
		t.Fatalf("expected generated label to start with %%PDF-1.4 header")
	}
	pdfStr := string(pdfBytes)
	for _, expectedText := range []string{"BIRO DATA DAN INFORMASI", "2026", "B-99", "F-01 s/d F-10"} {
		if !strings.Contains(pdfStr, expectedText) {
			t.Errorf("expected label PDF stream to contain %q", expectedText)
		}
	}
}

func TestPhase3_RateLimiter_Login6thRequestReturns429(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := middleware.NewRateLimiter(5, time.Minute)

	r := gin.New()
	r.POST("/api/auth/login", limiter.Middleware(), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 5 request pertama harus lolos (200 OK)
	for i := 1; i <= 5; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
		req.RemoteAddr = "192.168.1.100:12345"
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request #%d expected 200, got %d", i, w.Code)
		}
	}

	// Request ke-6 dalam 1 menit harus ditolak dengan 429 Too Many Requests
	w6 := httptest.NewRecorder()
	req6 := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	req6.RemoteAddr = "192.168.1.100:12345"
	r.ServeHTTP(w6, req6)
	if w6.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 6th login request in 1 minute to return 429, got %d", w6.Code)
	}
}

func TestPhase3_JWTHardening_And_SwaggerSpec(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:             "super-secret-256-bit-key-for-unit-testing",
			AccessExpireHours:  7,
			RefreshExpireHours: 30,
		},
	}
	authMw := middleware.NewAuthMiddleware(cfg, nil)

	tokenStr, _, err := authMw.GenerateToken(1, "administrator", "198501012010011001", "administrator", false)
	if err != nil || tokenStr == "" {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	// Verifikasi endpoint Swagger & OpenAPI JSON
	gin.SetMode(gin.TestMode)
	docsH := handler.NewDocsHandler()
	r := gin.New()
	r.GET("/swagger", docsH.ServeSwaggerUI)
	r.GET("/api/docs/openapi.json", docsH.ServeOpenAPISpec)

	wUI := httptest.NewRecorder()
	r.ServeHTTP(wUI, httptest.NewRequest(http.MethodGet, "/swagger", nil))
	if wUI.Code != http.StatusOK || !strings.Contains(wUI.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("expected /swagger to serve Swagger UI HTML")
	}

	wSpec := httptest.NewRecorder()
	r.ServeHTTP(wSpec, httptest.NewRequest(http.MethodGet, "/api/docs/openapi.json", nil))
	if wSpec.Code != http.StatusOK {
		t.Fatalf("expected /api/docs/openapi.json to return 200")
	}
	var spec map[string]interface{}
	if err := json.Unmarshal(wSpec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("expected valid OpenAPI JSON spec, got error: %v", err)
	}
	paths, _ := spec["paths"].(map[string]interface{})
	for _, requiredPath := range []string{"/api/auth/login", "/api/arsip", "/api/peminjaman", "/api/penyusutan", "/api/laporan/dai", "/api/label/box"} {
		if _, ok := paths[requiredPath]; !ok {
			t.Errorf("expected OpenAPI spec to contain path %s", requiredPath)
		}
	}
}
