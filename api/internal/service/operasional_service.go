package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// allowedDocumentMIME membatasi jenis file dokumen digitalisasi sesuai standar keamanan industri (PDF, JPG, PNG).
var allowedDocumentMIME = map[string]string{
	"application/pdf": ".pdf",
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
}

// OperasionalService mendefinisikan business logic untuk Peminjaman, Penyusutan, Dokumen Arsip, dan Audit Log.
type OperasionalService interface {
	// Peminjaman
	ListPeminjaman(ctx context.Context, req dto.PeminjamanListRequest) ([]dto.PeminjamanResponse, *dto.PaginatedResponse, error)
	GetPeminjamanByID(ctx context.Context, id uint) (*dto.PeminjamanResponse, error)
	CreatePeminjaman(ctx context.Context, req dto.StorePeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error)
	UpdatePeminjaman(ctx context.Context, id uint, req dto.UpdatePeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error)
	KembalikanPeminjaman(ctx context.Context, id uint, req dto.KembaliPeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error)
	DeletePeminjaman(ctx context.Context, id uint, userID uint, ipAddr string) error

	// Penyusutan
	ListPenyusutan(ctx context.Context, req dto.PenyusutanListRequest) ([]dto.PenyusutanResponse, *dto.PaginatedResponse, error)
	GetPenyusutanByID(ctx context.Context, id uint) (*dto.PenyusutanResponse, error)
	CreatePenyusutan(ctx context.Context, req dto.StorePenyusutanRequest, userID uint, ipAddr string) (*dto.PenyusutanResponse, error)
	ApprovePenyusutan(ctx context.Context, id uint, req dto.DecisionPenyusutanRequest, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error)
	RejectPenyusutan(ctx context.Context, id uint, req dto.DecisionPenyusutanRequest, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error)
	SelesaiPenyusutan(ctx context.Context, id uint, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error)

	// Dokumen Arsip
	ListDokumenByArsip(ctx context.Context, arsipID uint) ([]dto.DokumenArsipResponse, error)
	UploadDokumen(ctx context.Context, arsipID uint, originalFilename string, fileSize int64, fileReader io.Reader, keterangan *string, userID uint, ipAddr string) (*dto.DokumenArsipResponse, error)
	GetDokumenForDownload(ctx context.Context, id uint, callerRole string, callerID uint) (*model.DokumenArsip, error)
	DeleteDokumen(ctx context.Context, id uint, userID uint, ipAddr string) error

	// Audit Log
	RecordAudit(ctx context.Context, userID *uint, aksi, entitas string, entitasID *uint, dataLama, dataBaru interface{}, ipAddr string)
	ListAuditLog(ctx context.Context, req dto.AuditLogListRequest) ([]model.AuditLog, *dto.PaginatedResponse, error)
}

type operasionalServiceImpl struct {
	opRepo    repository.OperasionalRepository
	arsipRepo repository.ArsipRepository
	storage   config.StorageConfig
}

func NewOperasionalService(
	opRepo repository.OperasionalRepository,
	arsipRepo repository.ArsipRepository,
	storageCfg config.StorageConfig,
) OperasionalService {
	if storageCfg.Path == "" {
		storageCfg.Path = "./storage/dokumen"
	}
	if storageCfg.MaxUploadMB <= 0 {
		storageCfg.MaxUploadMB = 10
	}
	return &operasionalServiceImpl{
		opRepo:    opRepo,
		arsipRepo: arsipRepo,
		storage:   storageCfg,
	}
}

// ─── Helper & Audit Log ─────────────────────────────────────────

func (s *operasionalServiceImpl) RecordAudit(ctx context.Context, userID *uint, aksi, entitas string, entitasID *uint, dataLama, dataBaru interface{}, ipAddr string) {
	if s.opRepo == nil {
		return
	}
	var lamaStr, baruStr *string
	if dataLama != nil {
		if b, err := json.Marshal(dataLama); err == nil {
			str := string(b)
			lamaStr = &str
		}
	}
	if dataBaru != nil {
		if b, err := json.Marshal(dataBaru); err == nil {
			str := string(b)
			baruStr = &str
		}
	}
	var ipPtr *string
	if ipAddr != "" {
		ipPtr = &ipAddr
	}
	_ = s.opRepo.CreateAuditLog(ctx, &model.AuditLog{
		UserID:    userID,
		Aksi:      aksi,
		Entitas:   entitas,
		EntitasID: entitasID,
		DataLama:  lamaStr,
		DataBaru:  baruStr,
		IPAddress: ipPtr,
	})
}

func (s *operasionalServiceImpl) ListAuditLog(ctx context.Context, req dto.AuditLogListRequest) ([]model.AuditLog, *dto.PaginatedResponse, error) {
	return s.opRepo.FindAllAuditLog(ctx, req)
}

func parseDateISO(val, fieldName string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", strings.TrimSpace(val))
	if err != nil {
		return time.Time{}, helper.NewAppError(422, fmt.Sprintf("Format %s harus YYYY-MM-DD", fieldName), nil)
	}
	return t, nil
}

// computePeminjamanStatus menghitung status peminjaman secara dinamis (termasuk status Terlambat bila melewati tenggat).
func computePeminjamanStatus(tglRencana time.Time, tglAktual *time.Time, currentStatus string, now time.Time) string {
	if tglAktual != nil || strings.EqualFold(currentStatus, string(model.StatusPeminjamanDikembalikan)) {
		return string(model.StatusPeminjamanDikembalikan)
	}
	todayStr := now.Format("2006-01-02")
	rencanaStr := tglRencana.Format("2006-01-02")
	if rencanaStr < todayStr {
		return string(model.StatusPeminjamanTerlambat)
	}
	return string(model.StatusPeminjamanDipinjam)
}

// ─── 1. Peminjaman & Pengembalian ───────────────────────────────

func (s *operasionalServiceImpl) ListPeminjaman(ctx context.Context, req dto.PeminjamanListRequest) ([]dto.PeminjamanResponse, *dto.PaginatedResponse, error) {
	now := time.Now()
	_ = s.opRepo.SyncOverduePeminjaman(ctx, now)

	items, meta, err := s.opRepo.FindAllPeminjaman(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.PeminjamanResponse, 0, len(items))
	for _, item := range items {
		res = append(res, s.toPeminjamanResponse(item, now))
	}
	return res, meta, nil
}

func (s *operasionalServiceImpl) GetPeminjamanByID(ctx context.Context, id uint) (*dto.PeminjamanResponse, error) {
	now := time.Now()
	_ = s.opRepo.SyncOverduePeminjaman(ctx, now)

	item, err := s.opRepo.FindPeminjamanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}
	resp := s.toPeminjamanResponse(*item, now)
	return &resp, nil
}

func (s *operasionalServiceImpl) CreatePeminjaman(ctx context.Context, req dto.StorePeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error) {
	// Pastikan arsip ada
	if _, err := s.arsipRepo.FindByID(ctx, req.ArsipID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewAppError(422, "Arsip yang dipinjam tidak ditemukan", nil)
		}
		return nil, err
	}

	tglPinjam, err := parseDateISO(req.TanggalPinjam, "tanggal_pinjam")
	if err != nil {
		return nil, err
	}
	tglRencana, err := parseDateISO(req.TanggalKembaliRencana, "tanggal_kembali_rencana")
	if err != nil {
		return nil, err
	}
	if tglRencana.Before(tglPinjam) {
		return nil, helper.NewAppError(422, "Tanggal kembali rencana tidak boleh sebelum tanggal pinjam", nil)
	}

	now := time.Now()
	noTiket := strings.TrimSpace(req.NoTiket)
	if noTiket == "" {
		count, _ := s.opRepo.CountPeminjaman(ctx)
		noTiket = fmt.Sprintf("PIN/%s/%03d", now.Format("2006/01"), count+1)
	}

	status := computePeminjamanStatus(tglRencana, nil, string(model.StatusPeminjamanDipinjam), now)

	item := &model.Peminjaman{
		NoTiket:               noTiket,
		ArsipID:               req.ArsipID,
		Peminjam:              strings.TrimSpace(req.Peminjam),
		UnitKerjaID:           req.UnitKerjaID,
		TanggalPinjam:         tglPinjam,
		TanggalKembaliRencana: tglRencana,
		Status:                status,
		Catatan:               req.Catatan,
		CreatedBy:             userID,
	}

	if err := s.opRepo.CreatePeminjaman(ctx, item); err != nil {
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "create", "peminjaman", &item.ID, nil, item, ipAddr)

	created, err := s.opRepo.FindPeminjamanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPeminjamanResponse(*item, now)
		return &resp, nil
	}
	resp := s.toPeminjamanResponse(*created, now)
	return &resp, nil
}

func (s *operasionalServiceImpl) UpdatePeminjaman(ctx context.Context, id uint, req dto.UpdatePeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error) {
	item, err := s.opRepo.FindPeminjamanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	tglPinjam, err := parseDateISO(req.TanggalPinjam, "tanggal_pinjam")
	if err != nil {
		return nil, err
	}
	tglRencana, err := parseDateISO(req.TanggalKembaliRencana, "tanggal_kembali_rencana")
	if err != nil {
		return nil, err
	}
	if tglRencana.Before(tglPinjam) {
		return nil, helper.NewAppError(422, "Tanggal kembali rencana tidak boleh sebelum tanggal pinjam", nil)
	}

	oldCopy := *item
	now := time.Now()

	item.ArsipID = req.ArsipID
	item.Peminjam = strings.TrimSpace(req.Peminjam)
	item.UnitKerjaID = req.UnitKerjaID
	item.TanggalPinjam = tglPinjam
	item.TanggalKembaliRencana = tglRencana
	item.Catatan = req.Catatan
	item.Status = computePeminjamanStatus(tglRencana, item.TanggalKembaliAktual, item.Status, now)

	if err := s.opRepo.UpdatePeminjaman(ctx, item); err != nil {
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "update", "peminjaman", &item.ID, oldCopy, item, ipAddr)

	updated, err := s.opRepo.FindPeminjamanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPeminjamanResponse(*item, now)
		return &resp, nil
	}
	resp := s.toPeminjamanResponse(*updated, now)
	return &resp, nil
}

func (s *operasionalServiceImpl) KembalikanPeminjaman(ctx context.Context, id uint, req dto.KembaliPeminjamanRequest, userID uint, ipAddr string) (*dto.PeminjamanResponse, error) {
	item, err := s.opRepo.FindPeminjamanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	oldCopy := *item
	now := time.Now()
	tglKembali := now
	if strings.TrimSpace(req.TanggalKembaliAktual) != "" {
		parsed, errParse := parseDateISO(req.TanggalKembaliAktual, "tanggal_kembali_aktual")
		if errParse != nil {
			return nil, errParse
		}
		tglKembali = parsed
	}

	item.TanggalKembaliAktual = &tglKembali
	item.Status = string(model.StatusPeminjamanDikembalikan)
	if req.Catatan != nil {
		item.Catatan = req.Catatan
	}

	if err := s.opRepo.UpdatePeminjaman(ctx, item); err != nil {
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "update", "peminjaman", &item.ID, oldCopy, item, ipAddr)

	updated, err := s.opRepo.FindPeminjamanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPeminjamanResponse(*item, now)
		return &resp, nil
	}
	resp := s.toPeminjamanResponse(*updated, now)
	return &resp, nil
}

func (s *operasionalServiceImpl) DeletePeminjaman(ctx context.Context, id uint, userID uint, ipAddr string) error {
	item, err := s.opRepo.FindPeminjamanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	if err := s.opRepo.DeletePeminjaman(ctx, id); err != nil {
		return err
	}
	s.RecordAudit(ctx, &userID, "delete", "peminjaman", &id, item, nil, ipAddr)
	return nil
}

func (s *operasionalServiceImpl) toPeminjamanResponse(item model.Peminjaman, now time.Time) dto.PeminjamanResponse {
	var tglAktual *string
	if item.TanggalKembaliAktual != nil {
		str := item.TanggalKembaliAktual.Format("2006-01-02")
		tglAktual = &str
	}

	unitNama := item.Arsip.PenciptaArsip
	if item.UnitKerja != nil && item.UnitKerja.Nama != "" {
		unitNama = item.UnitKerja.Nama
	}

	status := computePeminjamanStatus(item.TanggalKembaliRencana, item.TanggalKembaliAktual, item.Status, now)

	return dto.PeminjamanResponse{
		ID:                    item.ID,
		NoTiket:               item.NoTiket,
		ArsipID:               item.ArsipID,
		NoItemArsip:           item.Arsip.NoItemArsip,
		KodeKlasifikasi:       item.Arsip.KodeKlasifikasi,
		UraianArsip:           item.Arsip.UraianInformasiArsip,
		Peminjam:              item.Peminjam,
		UnitKerjaID:           item.UnitKerjaID,
		UnitKerjaNama:         unitNama,
		TanggalPinjam:         item.TanggalPinjam.Format("2006-01-02"),
		TanggalKembaliRencana: item.TanggalKembaliRencana.Format("2006-01-02"),
		TanggalKembaliAktual:  tglAktual,
		Status:                status,
		Catatan:               item.Catatan,
		CreatedBy:             item.CreatedBy,
		CreatedAt:             item.CreatedAt,
		UpdatedAt:             item.UpdatedAt,
	}
}

// ─── 2. Penyusutan (Usul Musnah & Usul Serah) ───────────────────

func normalizeJenisUsul(raw string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "MUSNAH", "USUL MUSNAH":
		return string(model.JenisUsulMusnah), nil
	case "SERAH", "USUL SERAH":
		return string(model.JenisUsulSerah), nil
	default:
		return "", helper.NewAppError(422, "Jenis usul harus 'Musnah' atau 'Serah'", nil)
	}
}

func (s *operasionalServiceImpl) ListPenyusutan(ctx context.Context, req dto.PenyusutanListRequest) ([]dto.PenyusutanResponse, *dto.PaginatedResponse, error) {
	items, meta, err := s.opRepo.FindAllPenyusutan(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.PenyusutanResponse, 0, len(items))
	for _, item := range items {
		res = append(res, s.toPenyusutanResponse(item))
	}
	return res, meta, nil
}

func (s *operasionalServiceImpl) GetPenyusutanByID(ctx context.Context, id uint) (*dto.PenyusutanResponse, error) {
	item, err := s.opRepo.FindPenyusutanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}
	resp := s.toPenyusutanResponse(*item)
	return &resp, nil
}

func (s *operasionalServiceImpl) CreatePenyusutan(ctx context.Context, req dto.StorePenyusutanRequest, userID uint, ipAddr string) (*dto.PenyusutanResponse, error) {
	jenisUsul, err := normalizeJenisUsul(req.JenisUsul)
	if err != nil {
		return nil, err
	}

	arsip, err := s.arsipRepo.FindByID(ctx, req.ArsipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewAppError(422, "Arsip tidak ditemukan", nil)
		}
		return nil, err
	}

	tglRetensi, err := parseDateISO(req.TanggalRetensi, "tanggal_retensi")
	if err != nil {
		return nil, err
	}

	noPengajuan := strings.TrimSpace(req.NoPengajuan)
	if noPengajuan == "" {
		count, _ := s.opRepo.CountPenyusutan(ctx)
		prefix := "USM"
		if jenisUsul == string(model.JenisUsulSerah) {
			prefix = "USS"
		}
		noPengajuan = fmt.Sprintf("%s/%s/%03d", prefix, time.Now().Format("2006/01"), count+1)
	}

	item := &model.Penyusutan{
		NoPengajuan:    noPengajuan,
		JenisUsul:      jenisUsul,
		ArsipID:        req.ArsipID,
		TanggalRetensi: tglRetensi,
		Alasan:         req.Alasan,
		Status:         string(model.StatusPenyusutanDiajukan),
		DiajukanOleh:   userID,
	}

	if err := s.opRepo.CreatePenyusutan(ctx, item); err != nil {
		return nil, err
	}

	// Sinkronisasi kategori_arsip menjadi 'Usul Musnah' atau 'Usul Serah' sesuai Roadmap Task 2.2
	if jenisUsul == string(model.JenisUsulMusnah) {
		arsip.KategoriArsip = string(model.KategoriUsulMusnahVal)
	} else {
		arsip.KategoriArsip = string(model.KategoriUsulSerahVal)
	}
	_ = s.arsipRepo.Update(ctx, arsip)

	s.RecordAudit(ctx, &userID, "create", "penyusutan", &item.ID, nil, item, ipAddr)

	created, err := s.opRepo.FindPenyusutanByID(ctx, item.ID)
	if err != nil {
		item.Arsip = *arsip
		resp := s.toPenyusutanResponse(*item)
		return &resp, nil
	}
	resp := s.toPenyusutanResponse(*created)
	return &resp, nil
}

func (s *operasionalServiceImpl) ApprovePenyusutan(ctx context.Context, id uint, req dto.DecisionPenyusutanRequest, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error) {
	if !model.UserRole(callerRole).CanApprove() {
		return nil, helper.ErrForbidden
	}

	item, err := s.opRepo.FindPenyusutanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if item.Status != string(model.StatusPenyusutanDiajukan) {
		return nil, helper.NewAppError(422, fmt.Sprintf("Pengajuan dengan status '%s' tidak dapat disetujui", item.Status), nil)
	}

	oldCopy := *item
	item.Status = string(model.StatusPenyusutanDisetujui)
	item.DisetujuiOleh = &userID
	item.CatatanApprover = req.CatatanApprover

	if err := s.opRepo.UpdatePenyusutan(ctx, item); err != nil {
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "approve", "penyusutan", &item.ID, oldCopy, item, ipAddr)

	updated, err := s.opRepo.FindPenyusutanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPenyusutanResponse(*item)
		return &resp, nil
	}
	resp := s.toPenyusutanResponse(*updated)
	return &resp, nil
}

func (s *operasionalServiceImpl) RejectPenyusutan(ctx context.Context, id uint, req dto.DecisionPenyusutanRequest, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error) {
	if !model.UserRole(callerRole).CanApprove() {
		return nil, helper.ErrForbidden
	}

	item, err := s.opRepo.FindPenyusutanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if item.Status != string(model.StatusPenyusutanDiajukan) {
		return nil, helper.NewAppError(422, fmt.Sprintf("Pengajuan dengan status '%s' tidak dapat ditolak", item.Status), nil)
	}

	oldCopy := *item
	item.Status = string(model.StatusPenyusutanDitolak)
	item.DisetujuiOleh = &userID
	item.CatatanApprover = req.CatatanApprover

	if err := s.opRepo.UpdatePenyusutan(ctx, item); err != nil {
		return nil, err
	}

	// Kembalikan kategori_arsip ke Inaktif bila pengajuan ditolak
	if arsip, errArsip := s.arsipRepo.FindByID(ctx, item.ArsipID); errArsip == nil && arsip != nil {
		arsip.KategoriArsip = string(model.KategoriInaktifVal)
		_ = s.arsipRepo.Update(ctx, arsip)
	}

	s.RecordAudit(ctx, &userID, "reject", "penyusutan", &item.ID, oldCopy, item, ipAddr)

	updated, err := s.opRepo.FindPenyusutanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPenyusutanResponse(*item)
		return &resp, nil
	}
	resp := s.toPenyusutanResponse(*updated)
	return &resp, nil
}

func (s *operasionalServiceImpl) SelesaiPenyusutan(ctx context.Context, id uint, userID uint, callerRole, ipAddr string) (*dto.PenyusutanResponse, error) {
	if model.NormalizeRole(callerRole) != model.RoleAdministrator {
		return nil, helper.ErrForbidden
	}

	item, err := s.opRepo.FindPenyusutanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if item.Status != string(model.StatusPenyusutanDisetujui) {
		return nil, helper.NewAppError(422, "Hanya pengajuan yang sudah Disetujui yang dapat ditandai Selesai", nil)
	}

	oldCopy := *item
	item.Status = string(model.StatusPenyusutanSelesai)

	if err := s.opRepo.UpdatePenyusutan(ctx, item); err != nil {
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "update", "penyusutan", &item.ID, oldCopy, item, ipAddr)

	updated, err := s.opRepo.FindPenyusutanByID(ctx, item.ID)
	if err != nil {
		resp := s.toPenyusutanResponse(*item)
		return &resp, nil
	}
	resp := s.toPenyusutanResponse(*updated)
	return &resp, nil
}

func (s *operasionalServiceImpl) toPenyusutanResponse(item model.Penyusutan) dto.PenyusutanResponse {
	pengajuNama := ""
	if item.Pengaju.ID != 0 {
		pengajuNama = item.Pengaju.Name
	}
	approverNama := ""
	if item.Approver != nil && item.Approver.ID != 0 {
		approverNama = item.Approver.Name
	}
	nasib := ""
	if item.Arsip.NasibAkhir != nil {
		nasib = *item.Arsip.NasibAkhir
	} else if item.Arsip.JangkaSimpanNasibAkhir != nil {
		nasib = *item.Arsip.JangkaSimpanNasibAkhir
	}

	return dto.PenyusutanResponse{
		ID:              item.ID,
		NoPengajuan:     item.NoPengajuan,
		JenisUsul:       item.JenisUsul,
		ArsipID:         item.ArsipID,
		NoItemArsip:     item.Arsip.NoItemArsip,
		KodeKlasifikasi: item.Arsip.KodeKlasifikasi,
		UraianArsip:     item.Arsip.UraianInformasiArsip,
		NasibAkhir:      nasib,
		TanggalRetensi:  item.TanggalRetensi.Format("2006-01-02"),
		Alasan:          item.Alasan,
		Status:          item.Status,
		DiajukanOleh:    item.DiajukanOleh,
		DiajukanNama:    pengajuNama,
		DisetujuiOleh:   item.DisetujuiOleh,
		DisetujuiNama:   approverNama,
		CatatanApprover: item.CatatanApprover,
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}
}

// ─── 3. Dokumen Arsip (Digitalisasi) ────────────────────────────

func generateUUIDFilename(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant
	return fmt.Sprintf("%x-%x-%x-%x-%x%s", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16], ext), nil
}

func (s *operasionalServiceImpl) ListDokumenByArsip(ctx context.Context, arsipID uint) ([]dto.DokumenArsipResponse, error) {
	if _, err := s.arsipRepo.FindByID(ctx, arsipID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	items, err := s.opRepo.FindDokumenByArsipID(ctx, arsipID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.DokumenArsipResponse, 0, len(items))
	for _, d := range items {
		res = append(res, s.toDokumenResponse(d))
	}
	return res, nil
}

func (s *operasionalServiceImpl) UploadDokumen(
	ctx context.Context,
	arsipID uint,
	originalFilename string,
	fileSize int64,
	fileReader io.Reader,
	keterangan *string,
	userID uint,
	ipAddr string,
) (*dto.DokumenArsipResponse, error) {
	// 1. Cek arsip exists
	if _, err := s.arsipRepo.FindByID(ctx, arsipID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	// 2. Validasi ukuran file (MAX_UPLOAD_MB)
	maxBytes := int64(s.storage.MaxUploadMB) * 1024 * 1024
	if fileSize <= 0 {
		return nil, helper.NewAppError(422, "File dokumen tidak boleh kosong", nil)
	}
	if fileSize > maxBytes {
		return nil, helper.NewAppError(422, fmt.Sprintf("Ukuran file melebihi batas maksimal %d MB", s.storage.MaxUploadMB), nil)
	}

	// 3. Baca maksimal maxBytes + 1 dan deteksi MIME berdasarkan magic bytes (512 byte pertama)
	limitedReader := io.LimitReader(fileReader, maxBytes+1)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, helper.NewAppError(400, "Gagal membaca konten file", err)
	}
	if int64(len(content)) > maxBytes {
		return nil, helper.NewAppError(422, fmt.Sprintf("Ukuran file melebihi batas maksimal %d MB", s.storage.MaxUploadMB), nil)
	}

	detectedMIME := http.DetectContentType(content)
	// Normalisasi parameter tambahan jika ada (mis. "; charset=...")
	if idx := strings.Index(detectedMIME, ";"); idx != -1 {
		detectedMIME = strings.TrimSpace(detectedMIME[:idx])
	}

	ext, allowed := allowedDocumentMIME[detectedMIME]
	if !allowed {
		return nil, helper.NewAppError(422, fmt.Sprintf("Tipe file '%s' tidak diizinkan (hanya PDF, JPG, PNG)", detectedMIME), nil)
	}

	// 4. Simpan ke STORAGE_PATH dengan nama file UUID v4
	if err := os.MkdirAll(s.storage.Path, 0o750); err != nil {
		return nil, helper.NewAppError(500, "Gagal menyiapkan direktori penyimpanan dokumen", err)
	}

	uuidFile, err := generateUUIDFilename(ext)
	if err != nil {
		uuidFile = fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), hex.EncodeToString([]byte(originalFilename))[:8], ext)
	}
	fullPath := filepath.Join(s.storage.Path, uuidFile)

	if err := os.WriteFile(fullPath, content, 0o640); err != nil {
		return nil, helper.NewAppError(500, "Gagal menyimpan file dokumen", err)
	}

	safeName := filepath.Base(strings.TrimSpace(originalFilename))
	if safeName == "" || safeName == "." {
		safeName = uuidFile
	}

	doc := &model.DokumenArsip{
		ArsipID:     arsipID,
		NamaFile:    safeName,
		Path:        fullPath,
		MimeType:    detectedMIME,
		UkuranBytes: int64(len(content)),
		Keterangan:  keterangan,
		UploadedBy:  userID,
	}

	if err := s.opRepo.CreateDokumen(ctx, doc); err != nil {
		_ = os.Remove(fullPath)
		return nil, err
	}

	s.RecordAudit(ctx, &userID, "create", "dokumen_arsip", &doc.ID, nil, doc, ipAddr)

	resp := s.toDokumenResponse(*doc)
	return &resp, nil
}

func (s *operasionalServiceImpl) GetDokumenForDownload(ctx context.Context, id uint, callerRole string, callerID uint) (*model.DokumenArsip, error) {
	doc, err := s.opRepo.FindDokumenByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	// Cek otorisasi SKKAD: jika arsip berstatus Rahasia, viewer tidak boleh mengunduh
	if arsip, errArsip := s.arsipRepo.FindByID(ctx, doc.ArsipID); errArsip == nil && arsip != nil {
		if strings.EqualFold(arsip.SKKAD, string(model.SKKADRahasia)) {
			normRole := model.NormalizeRole(callerRole)
			if normRole == model.RoleViewer {
				return nil, helper.ErrForbidden
			}
		}
	}

	return doc, nil
}

func (s *operasionalServiceImpl) DeleteDokumen(ctx context.Context, id uint, userID uint, ipAddr string) error {
	doc, err := s.opRepo.FindDokumenByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}

	if err := s.opRepo.DeleteDokumen(ctx, id); err != nil {
		return err
	}
	if doc.Path != "" {
		_ = os.Remove(doc.Path)
	}

	s.RecordAudit(ctx, &userID, "delete", "dokumen_arsip", &id, doc, nil, ipAddr)
	return nil
}

func (s *operasionalServiceImpl) toDokumenResponse(d model.DokumenArsip) dto.DokumenArsipResponse {
	return dto.DokumenArsipResponse{
		ID:          d.ID,
		ArsipID:     d.ArsipID,
		NamaFile:    d.NamaFile,
		MimeType:    d.MimeType,
		UkuranBytes: d.UkuranBytes,
		Keterangan:  d.Keterangan,
		DownloadURL: fmt.Sprintf("/api/dokumen/%d/download", d.ID),
		UploadedBy:  d.UploadedBy,
		CreatedAt:   d.CreatedAt,
	}
}
