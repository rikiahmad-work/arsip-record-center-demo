package service

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// ─── Fake OperasionalRepository for Phase 2 Unit Testing ────────

type fakeOperasionalRepo struct {
	peminjaman []*model.Peminjaman
	penyusutan []*model.Penyusutan
	dokumen    []*model.DokumenArsip
	audits     []*model.AuditLog
}

func (f *fakeOperasionalRepo) SyncOverduePeminjaman(ctx context.Context, now time.Time) error {
	today := now.Format("2006-01-02")
	for _, p := range f.peminjaman {
		if p.Status == string(model.StatusPeminjamanDipinjam) && p.TanggalKembaliAktual == nil && p.TanggalKembaliRencana.Format("2006-01-02") < today {
			p.Status = string(model.StatusPeminjamanTerlambat)
		}
	}
	return nil
}

func (f *fakeOperasionalRepo) FindAllPeminjaman(ctx context.Context, req dto.PeminjamanListRequest) ([]model.Peminjaman, *dto.PaginatedResponse, error) {
	var res []model.Peminjaman
	for _, p := range f.peminjaman {
		res = append(res, *p)
	}
	return res, &dto.PaginatedResponse{CurrentPage: 1, LastPage: 1, PerPage: 10, Total: int64(len(res))}, nil
}

func (f *fakeOperasionalRepo) FindPeminjamanByID(ctx context.Context, id uint) (*model.Peminjaman, error) {
	for _, p := range f.peminjaman {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeOperasionalRepo) CreatePeminjaman(ctx context.Context, item *model.Peminjaman) error {
	item.ID = uint(len(f.peminjaman) + 1)
	f.peminjaman = append(f.peminjaman, item)
	return nil
}

func (f *fakeOperasionalRepo) UpdatePeminjaman(ctx context.Context, item *model.Peminjaman) error {
	return nil
}

func (f *fakeOperasionalRepo) DeletePeminjaman(ctx context.Context, id uint) error {
	return nil
}

func (f *fakeOperasionalRepo) CountPeminjaman(ctx context.Context) (int64, error) {
	return int64(len(f.peminjaman)), nil
}

func (f *fakeOperasionalRepo) FindAllPenyusutan(ctx context.Context, req dto.PenyusutanListRequest) ([]model.Penyusutan, *dto.PaginatedResponse, error) {
	var res []model.Penyusutan
	for _, p := range f.penyusutan {
		res = append(res, *p)
	}
	return res, &dto.PaginatedResponse{CurrentPage: 1, LastPage: 1, PerPage: 10, Total: int64(len(res))}, nil
}

func (f *fakeOperasionalRepo) FindPenyusutanByID(ctx context.Context, id uint) (*model.Penyusutan, error) {
	for _, p := range f.penyusutan {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeOperasionalRepo) CreatePenyusutan(ctx context.Context, item *model.Penyusutan) error {
	item.ID = uint(len(f.penyusutan) + 1)
	f.penyusutan = append(f.penyusutan, item)
	return nil
}

func (f *fakeOperasionalRepo) UpdatePenyusutan(ctx context.Context, item *model.Penyusutan) error {
	return nil
}

func (f *fakeOperasionalRepo) CountPenyusutan(ctx context.Context) (int64, error) {
	return int64(len(f.penyusutan)), nil
}

func (f *fakeOperasionalRepo) FindDokumenByArsipID(ctx context.Context, arsipID uint) ([]model.DokumenArsip, error) {
	var res []model.DokumenArsip
	for _, d := range f.dokumen {
		if d.ArsipID == arsipID {
			res = append(res, *d)
		}
	}
	return res, nil
}

func (f *fakeOperasionalRepo) FindDokumenByID(ctx context.Context, id uint) (*model.DokumenArsip, error) {
	for _, d := range f.dokumen {
		if d.ID == id {
			return d, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeOperasionalRepo) CreateDokumen(ctx context.Context, item *model.DokumenArsip) error {
	item.ID = uint(len(f.dokumen) + 1)
	f.dokumen = append(f.dokumen, item)
	return nil
}

func (f *fakeOperasionalRepo) DeleteDokumen(ctx context.Context, id uint) error {
	return nil
}

func (f *fakeOperasionalRepo) CreateAuditLog(ctx context.Context, item *model.AuditLog) error {
	item.ID = uint(len(f.audits) + 1)
	f.audits = append(f.audits, item)
	return nil
}

func (f *fakeOperasionalRepo) FindAllAuditLog(ctx context.Context, req dto.AuditLogListRequest) ([]model.AuditLog, *dto.PaginatedResponse, error) {
	var res []model.AuditLog
	for _, a := range f.audits {
		res = append(res, *a)
	}
	return res, &dto.PaginatedResponse{CurrentPage: 1, LastPage: 1, PerPage: 20, Total: int64(len(res))}, nil
}

// ─── Phase 2 Acceptance Criteria Tests ──────────────────────────

func TestPeminjaman_OverdueAndReturnFlow(t *testing.T) {
	opRepo := &fakeOperasionalRepo{}
	arsipRepo := &fakeArsipRepo{
		items: []*model.Arsip{
			{ID: 1, NoItemArsip: "ARS/2026/DT.04.01/001", KodeKlasifikasi: "DT.04.01", UraianInformasiArsip: "Arsip Uji", KategoriArsip: "Inaktif", SKKAD: "Biasa"},
		},
	}
	svc := NewOperasionalService(opRepo, arsipRepo, config.StorageConfig{Path: t.TempDir(), MaxUploadMB: 1})

	// 1. Buat peminjaman dengan tanggal_kembali_rencana yang sudah lewat (2020-01-10) -> otomatis berstatus Terlambat
	resp, err := svc.CreatePeminjaman(context.Background(), dto.StorePeminjamanRequest{
		ArsipID:               1,
		Peminjam:              "Budi Santoso",
		TanggalPinjam:         "2020-01-01",
		TanggalKembaliRencana: "2020-01-10",
	}, 2, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error creating peminjaman: %v", err)
	}
	if resp.Status != string(model.StatusPeminjamanTerlambat) {
		t.Errorf("expected overdue status 'Terlambat', got '%s'", resp.Status)
	}

	// 2. Kembalikan arsip -> status berubah menjadi Dikembalikan
	returned, err := svc.KembalikanPeminjaman(context.Background(), resp.ID, dto.KembaliPeminjamanRequest{}, 2, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error returning peminjaman: %v", err)
	}
	if returned.Status != string(model.StatusPeminjamanDikembalikan) {
		t.Errorf("expected status 'Dikembalikan', got '%s'", returned.Status)
	}
	if returned.TanggalKembaliAktual == nil {
		t.Errorf("expected tanggal_kembali_aktual to be populated")
	}

	// 3. Pastikan terekam di audit_log
	if len(opRepo.audits) < 2 {
		t.Errorf("expected at least 2 audit logs, got %d", len(opRepo.audits))
	}
}

func TestPenyusutan_CategorySync_And_ApproverAuthorization(t *testing.T) {
	opRepo := &fakeOperasionalRepo{}
	arsipItem := &model.Arsip{
		ID:                   1,
		NoItemArsip:          "ARS/2026/DT.04.01/001",
		KodeKlasifikasi:      "DT.04.01",
		UraianInformasiArsip: "Arsip Inaktif Siap Susut",
		KategoriArsip:        "Inaktif",
		SKKAD:                "Biasa",
	}
	arsipRepo := &fakeArsipRepo{items: []*model.Arsip{arsipItem}}
	svc := NewOperasionalService(opRepo, arsipRepo, config.StorageConfig{Path: t.TempDir(), MaxUploadMB: 1})

	// 1. Operator mengajukan Usul Musnah -> kategori_arsip otomatis menjadi 'Usul Musnah'
	created, err := svc.CreatePenyusutan(context.Background(), dto.StorePenyusutanRequest{
		JenisUsul:      "Musnah",
		ArsipID:        1,
		TanggalRetensi: "2026-12-31",
	}, 2, "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error creating penyusutan: %v", err)
	}
	if arsipItem.KategoriArsip != string(model.KategoriUsulMusnahVal) {
		t.Errorf("expected arsip.kategori_arsip='Usul Musnah', got '%s'", arsipItem.KategoriArsip)
	}

	// 2. Operator mencoba menyetujui -> ditolak dengan 403 Forbidden
	_, errForbidden := svc.ApprovePenyusutan(context.Background(), created.ID, dto.DecisionPenyusutanRequest{}, 2, string(model.RoleOperator), "127.0.0.1")
	if errForbidden == nil {
		t.Fatalf("expected 403 error when operator tries to approve, got nil")
	}
	if appErr := helper.GetAppError(errForbidden); appErr == nil || appErr.Code != 403 {
		t.Fatalf("expected HTTP 403 AppError, got: %v", errForbidden)
	}

	// 3. Approver menyetujui -> status menjadi Disetujui
	note := "Disetujui sesuai JRA"
	approved, err := svc.ApprovePenyusutan(context.Background(), created.ID, dto.DecisionPenyusutanRequest{
		CatatanApprover: &note,
	}, 4, string(model.RoleApprover), "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error when approver approves: %v", err)
	}
	if approved.Status != string(model.StatusPenyusutanDisetujui) {
		t.Errorf("expected status 'Disetujui', got '%s'", approved.Status)
	}

	// 4. Administrator menandai Selesai -> status menjadi Selesai
	finished, err := svc.SelesaiPenyusutan(context.Background(), created.ID, 1, string(model.RoleAdministrator), "127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error when admin marks selesai: %v", err)
	}
	if finished.Status != string(model.StatusPenyusutanSelesai) {
		t.Errorf("expected status 'Selesai', got '%s'", finished.Status)
	}
}

func TestDokumenUpload_MIMEAndSizeValidation(t *testing.T) {
	opRepo := &fakeOperasionalRepo{}
	arsipRepo := &fakeArsipRepo{
		items: []*model.Arsip{
			{ID: 1, NoItemArsip: "1", KodeKlasifikasi: "DT.04.01", SKKAD: "Rahasia"},
		},
	}
	svc := NewOperasionalService(opRepo, arsipRepo, config.StorageConfig{Path: t.TempDir(), MaxUploadMB: 1})

	// 1. Upload file PDF valid (%PDF-1.7 magic header) -> sukses
	pdfData := append([]byte("%PDF-1.7\n"), bytes.Repeat([]byte("A"), 128)...)
	docResp, err := svc.UploadDokumen(context.Background(), 1, "dokumen_uji.pdf", int64(len(pdfData)), bytes.NewReader(pdfData), nil, 2, "127.0.0.1")
	if err != nil {
		t.Fatalf("expected valid PDF upload to succeed, got: %v", err)
	}
	if docResp.MimeType != "application/pdf" {
		t.Errorf("expected mime_type='application/pdf', got '%s'", docResp.MimeType)
	}

	// 2. Upload file dengan MIME tidak diizinkan (text/plain) -> 422
	txtData := []byte("Halo ini file teks biasa, bukan PDF atau gambar")
	_, errMIME := svc.UploadDokumen(context.Background(), 1, "script.sh", int64(len(txtData)), bytes.NewReader(txtData), nil, 2, "127.0.0.1")
	if errMIME == nil {
		t.Fatalf("expected 422 error for invalid MIME type, got nil")
	}
	if appErr := helper.GetAppError(errMIME); appErr == nil || appErr.Code != 422 {
		t.Fatalf("expected HTTP 422 for invalid MIME, got: %v", errMIME)
	}

	// 3. Upload file melebihi batas ukuran (> 1 MB) -> 422
	largeData := append([]byte("%PDF-1.7\n"), bytes.Repeat([]byte("X"), 2*1024*1024)...)
	_, errSize := svc.UploadDokumen(context.Background(), 1, "besar.pdf", int64(len(largeData)), bytes.NewReader(largeData), nil, 2, "127.0.0.1")
	if errSize == nil {
		t.Fatalf("expected 422 error for oversized file, got nil")
	}
	if appErr := helper.GetAppError(errSize); appErr == nil || appErr.Code != 422 {
		t.Fatalf("expected HTTP 422 for oversized file, got: %v", errSize)
	}

	// 4. Viewer mencoba mengunduh dokumen dari arsip ber-SKKAD Rahasia -> 403 Forbidden
	_, errViewerDownload := svc.GetDokumenForDownload(context.Background(), docResp.ID, string(model.RoleViewer), 3)
	if errViewerDownload == nil {
		t.Fatalf("expected 403 when viewer downloads Rahasia document")
	}
}
