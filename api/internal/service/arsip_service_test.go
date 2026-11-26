package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// ─── Fake Repositories for Table-Driven Unit Testing ────────────

type fakeMasterRepo struct {
	kodes  map[string]*model.KodeKlasifikasi
	units  map[uint]*model.UnitKerja
	lokasi map[uint]*model.LokasiSimpan
}

func newFakeMasterRepo() *fakeMasterRepo {
	inaktif3 := "3 Tahun"
	inaktif4 := "4 Tahun"
	musnah := "Musnah"
	permanen := "Permanen"

	return &fakeMasterRepo{
		kodes: map[string]*model.KodeKlasifikasi{
			"DT.04.01": {
				ID:                  1,
				Kode:                "DT.04.01",
				Nama:                "Sub-Sub Masalah Data dan Informasi (DT.04.01)",
				JangkaSimpanInaktif: &inaktif4,
				NasibAkhir:          &musnah,
				SKKAD:               "Biasa",
				IsActive:            true,
			},
			"KA.03.06": {
				ID:                  2,
				Kode:                "KA.03.06",
				Nama:                "Sub-Sub Masalah Kearsipan (KA.03.06)",
				JangkaSimpanInaktif: &inaktif3,
				NasibAkhir:          &permanen,
				SKKAD:               "Rahasia",
				IsActive:            true,
			},
		},
		units: map[uint]*model.UnitKerja{
			1: {ID: 1, Kode: "BU.01", Nama: "Biro Umum & SDM", IsActive: true},
		},
		lokasi: map[uint]*model.LokasiSimpan{
			1: {ID: 1, Nama: "Record Center", IsActive: true},
		},
	}
}

func (f *fakeMasterRepo) FindAllKodeKlasifikasi(ctx context.Context, req dto.KodeKlasifikasiListRequest) ([]model.KodeKlasifikasi, *dto.PaginatedResponse, error) {
	return nil, nil, nil
}
func (f *fakeMasterRepo) FindKodeKlasifikasiByID(ctx context.Context, id uint) (*model.KodeKlasifikasi, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeMasterRepo) FindKodeKlasifikasiByKode(ctx context.Context, kode string) (*model.KodeKlasifikasi, error) {
	if item, ok := f.kodes[kode]; ok {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeMasterRepo) CreateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error {
	return nil
}
func (f *fakeMasterRepo) UpdateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error {
	return nil
}
func (f *fakeMasterRepo) DeleteKodeKlasifikasi(ctx context.Context, id uint) error { return nil }
func (f *fakeMasterRepo) CountArsipPerKode(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}
func (f *fakeMasterRepo) FindAllLokasiSimpan(ctx context.Context, req dto.LokasiSimpanListRequest) ([]model.LokasiSimpan, *dto.PaginatedResponse, error) {
	return nil, nil, nil
}
func (f *fakeMasterRepo) FindLokasiSimpanByID(ctx context.Context, id uint) (*model.LokasiSimpan, error) {
	if item, ok := f.lokasi[id]; ok {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeMasterRepo) CreateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error {
	return nil
}
func (f *fakeMasterRepo) UpdateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error {
	return nil
}
func (f *fakeMasterRepo) DeleteLokasiSimpan(ctx context.Context, id uint) error { return nil }
func (f *fakeMasterRepo) FindAllUnitKerja(ctx context.Context, req dto.UnitKerjaListRequest) ([]model.UnitKerja, *dto.PaginatedResponse, error) {
	return nil, nil, nil
}
func (f *fakeMasterRepo) FindUnitKerjaByID(ctx context.Context, id uint) (*model.UnitKerja, error) {
	if item, ok := f.units[id]; ok {
		return item, nil
	}
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeMasterRepo) FindUnitKerjaByKode(ctx context.Context, kode string) (*model.UnitKerja, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f *fakeMasterRepo) CreateUnitKerja(ctx context.Context, item *model.UnitKerja) error {
	return nil
}
func (f *fakeMasterRepo) UpdateUnitKerja(ctx context.Context, item *model.UnitKerja) error {
	return nil
}
func (f *fakeMasterRepo) DeleteUnitKerja(ctx context.Context, id uint) error { return nil }
func (f *fakeMasterRepo) CountArsipPerUnitKerja(ctx context.Context) (map[uint]int64, error) {
	return map[uint]int64{}, nil
}
func (f *fakeMasterRepo) GetAllKonfigurasi(ctx context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}
func (f *fakeMasterRepo) UpsertKonfigurasi(ctx context.Context, items map[string]string, updatedBy uint) error {
	return nil
}

type fakeArsipRepo struct {
	items []*model.Arsip
}

func (r *fakeArsipRepo) FindAll(ctx context.Context, req dto.ArsipListRequest) ([]model.Arsip, *dto.PaginatedResponse, error) {
	var res []model.Arsip
	for _, a := range r.items {
		if req.ResolvedKode() != "" && a.KodeKlasifikasi != req.ResolvedKode() {
			continue
		}
		if req.ResolvedKategori() != "" && !strings.EqualFold(a.KategoriArsip, req.ResolvedKategori()) {
			continue
		}
		if req.SKKAD != "" && !strings.EqualFold(a.SKKAD, req.SKKAD) {
			continue
		}
		res = append(res, *a)
	}
	return res, &dto.PaginatedResponse{CurrentPage: 1, LastPage: 1, PerPage: 10, Total: int64(len(res))}, nil
}
func (r *fakeArsipRepo) FindByID(ctx context.Context, id uint) (*model.Arsip, error) {
	for _, a := range r.items {
		if a.ID == id {
			return a, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeArsipRepo) FindByNoItem(ctx context.Context, noItem string) (*model.Arsip, error) {
	for _, a := range r.items {
		if a.NoItemArsip == noItem {
			return a, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *fakeArsipRepo) Create(ctx context.Context, arsip *model.Arsip) error {
	arsip.ID = uint(len(r.items) + 1)
	r.items = append(r.items, arsip)
	return nil
}
func (r *fakeArsipRepo) Update(ctx context.Context, arsip *model.Arsip) error { return nil }
func (r *fakeArsipRepo) Delete(ctx context.Context, id uint) error            { return nil }
func (r *fakeArsipRepo) Restore(ctx context.Context, id uint) (*model.Arsip, error) {
	return nil, nil
}
func (r *fakeArsipRepo) CountByKategori(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}
func (r *fakeArsipRepo) CountBySKKAD(ctx context.Context) (map[string]int64, error) {
	return map[string]int64{}, nil
}
func (r *fakeArsipRepo) CountByStatus(ctx context.Context) (int64, int64, error) { return 0, 0, nil }
func (r *fakeArsipRepo) CountTotal(ctx context.Context) (int64, error) {
	return int64(len(r.items)), nil
}
func (r *fakeArsipRepo) CountByUnor(ctx context.Context, unorID string) (int64, error) {
	var c int64
	for _, a := range r.items {
		if a.UnorID == unorID {
			c++
		}
	}
	return c, nil
}
func (r *fakeArsipRepo) NextNoItemByUnor(ctx context.Context, unorID string, unitKerjaID *uint, penciptaArsip string) (string, error) {
	maxNum := 0
	count := 0
	for _, a := range r.items {
		if a.UnorID == unorID || (penciptaArsip != "" && strings.EqualFold(a.PenciptaArsip, penciptaArsip)) {
			count++
			if n, err := strconv.Atoi(a.NoItemArsip); err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	next := maxNum + 1
	if count+1 > next {
		next = count + 1
	}
	return strconv.Itoa(next), nil
}

// ─── Tests for Phase 1 Acceptance Criteria ──────────────────────

func TestArsipService_AutoFillJRA_And_UnknownKode422(t *testing.T) {
	arsipRepo := &fakeArsipRepo{}
	masterRepo := newFakeMasterRepo()
	svc := NewArsipService(arsipRepo, masterRepo)

	folder := "F-01"
	unitID := uint(1)
	lokasiID := uint(1)

	// 1. Auto-fill JRA & SKKAD dari master kode_klasifikasi (Operator mencoba spoof SKKAD="Biasa" pada kode KA.03.06 yang Rahasia)
	resp, err := svc.Create(context.Background(), dto.StoreArsipRequest{
		NoFolder:            &folder,
		KodeKlasifikasi:     "KA.03.06",
		UraianInformasi:     "Dokumen Audit Kearsipan Rahasia",
		JenisArsip:          "Konvensional",
		KurunWaktu:          "2024",
		TingkatPerkembangan: "Asli",
		Jumlah:              "1 Berkas",
		LokasiSimpanID:      &lokasiID,
		NoBoks:              "B-01",
		NoRak:               "R-01",
		SKKAD:               "Biasa", // Harus diabaikan karena pemanggil adalah operator!
		KategoriArsip:       "Vital",
		UnitKerjaID:         &unitID,
	}, 2, string(model.RoleOperator))
	if err != nil {
		t.Fatalf("unexpected error on valid kode: %v", err)
	}

	if resp.JangkaSimpanInaktif != "3 Tahun" {
		t.Errorf("expected auto-filled JangkaSimpanInaktif='3 Tahun', got '%s'", resp.JangkaSimpanInaktif)
	}
	if resp.NasibAkhir == nil || *resp.NasibAkhir != "Permanen" {
		t.Errorf("expected auto-filled NasibAkhir='Permanen', got '%v'", resp.NasibAkhir)
	}
	if resp.SKKAD != "Rahasia" {
		t.Errorf("expected auto-filled SKKAD='Rahasia' from master (ignoring operator override), got '%s'", resp.SKKAD)
	}
	if resp.NoFolder == nil || *resp.NoFolder != "F-01" {
		t.Errorf("expected NoFolder='F-01', got '%v'", resp.NoFolder)
	}

	// 2. Administrator boleh override SKKAD
	respAdmin, err := svc.Create(context.Background(), dto.StoreArsipRequest{
		KodeKlasifikasi:     "DT.04.01",
		UraianInformasi:     "Nota Dinas Khusus Terbatas",
		JenisArsip:          "Elektronik",
		KurunWaktu:          "2025",
		TingkatPerkembangan: "Copy",
		Jumlah:              "1 File",
		SKKAD:               "Terbatas", // Override oleh administrator dari Biasa -> Terbatas
		KategoriArsip:       "Inaktif",
	}, 1, string(model.RoleAdministrator))
	if err != nil {
		t.Fatalf("unexpected error for admin create: %v", err)
	}
	if respAdmin.SKKAD != "Terbatas" {
		t.Errorf("expected admin SKKAD override 'Terbatas', got '%s'", respAdmin.SKKAD)
	}

	// 3. Kode tidak dikenal -> 422 Unprocessable Entity
	_, errUnknown := svc.Create(context.Background(), dto.StoreArsipRequest{
		KodeKlasifikasi:     "XX.99.99",
		UraianInformasi:     "Arsip dengan kode tidak valid",
		JenisArsip:          "Konvensional",
		KurunWaktu:          "2025",
		TingkatPerkembangan: "Asli",
		Jumlah:              "1 Berkas",
	}, 2, string(model.RoleOperator))
	if errUnknown == nil {
		t.Fatalf("expected error for unknown kode_klasifikasi, got nil")
	}
	appErr := helper.GetAppError(errUnknown)
	if appErr == nil || appErr.Code != 422 {
		t.Fatalf("expected HTTP 422 AppError for unknown kode, got: %v", errUnknown)
	}

	// 4. Filter list arsip berdasarkan kategori=Vital & skkad=Rahasia
	filtered, _, errFilter := svc.GetAll(context.Background(), dto.ArsipListRequest{
		Kategori: "Vital",
		SKKAD:    "Rahasia",
	})
	if errFilter != nil {
		t.Fatalf("unexpected error on filter: %v", errFilter)
	}
	if len(filtered) != 1 || filtered[0].KodeKlasifikasi != "KA.03.06" {
		t.Fatalf("expected 1 filtered item (KA.03.06), got %d", len(filtered))
	}
}

func TestRoleMiddleware_4RolesAuthorization(t *testing.T) {
	gin.SetMode(gin.TestMode)
	roleMw := middleware.NewRoleMiddleware()

	tests := []struct {
		name       string
		callerRole string
		allowed    []model.UserRole
		wantHTTP   int
	}{
		{
			name:       "viewer forbidden on POST /api/arsip",
			callerRole: "viewer",
			allowed:    []model.UserRole{model.RoleAdministrator, model.RoleOperator},
			wantHTTP:   http.StatusForbidden,
		},
		{
			name:       "approver forbidden on POST /api/arsip",
			callerRole: "approver",
			allowed:    []model.UserRole{model.RoleAdministrator, model.RoleOperator},
			wantHTTP:   http.StatusForbidden,
		},
		{
			name:       "operator allowed on POST /api/arsip",
			callerRole: "operator",
			allowed:    []model.UserRole{model.RoleAdministrator, model.RoleOperator},
			wantHTTP:   http.StatusOK,
		},
		{
			name:       "legacy super-admin mapped to administrator and allowed",
			callerRole: "super-admin",
			allowed:    []model.UserRole{model.RoleAdministrator},
			wantHTTP:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.POST("/api/arsip", func(ctx *gin.Context) {
				ctx.Set("user_role", tt.callerRole)
				ctx.Next()
			}, roleMw.RequireRole(tt.allowed...), func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			c.Request = httptest.NewRequest(http.MethodPost, "/api/arsip", nil)
			r.ServeHTTP(w, c.Request)

			if w.Code != tt.wantHTTP {
				t.Errorf("expected HTTP %d, got %d", tt.wantHTTP, w.Code)
			}
		})
	}
}
