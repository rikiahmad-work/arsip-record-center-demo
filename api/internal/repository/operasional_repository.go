package repository

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
)

// OperasionalRepository mendefinisikan akses data untuk peminjaman, penyusutan, dokumen_arsip, dan audit_log.
type OperasionalRepository interface {
	// Peminjaman
	SyncOverduePeminjaman(ctx context.Context, now time.Time) error
	FindAllPeminjaman(ctx context.Context, req dto.PeminjamanListRequest) ([]model.Peminjaman, *dto.PaginatedResponse, error)
	FindPeminjamanByID(ctx context.Context, id uint) (*model.Peminjaman, error)
	CreatePeminjaman(ctx context.Context, item *model.Peminjaman) error
	UpdatePeminjaman(ctx context.Context, item *model.Peminjaman) error
	DeletePeminjaman(ctx context.Context, id uint) error
	CountPeminjaman(ctx context.Context) (int64, error)

	// Penyusutan
	FindAllPenyusutan(ctx context.Context, req dto.PenyusutanListRequest) ([]model.Penyusutan, *dto.PaginatedResponse, error)
	FindPenyusutanByID(ctx context.Context, id uint) (*model.Penyusutan, error)
	CreatePenyusutan(ctx context.Context, item *model.Penyusutan) error
	UpdatePenyusutan(ctx context.Context, item *model.Penyusutan) error
	CountPenyusutan(ctx context.Context) (int64, error)

	// Dokumen Arsip
	FindDokumenByArsipID(ctx context.Context, arsipID uint) ([]model.DokumenArsip, error)
	FindDokumenByID(ctx context.Context, id uint) (*model.DokumenArsip, error)
	CreateDokumen(ctx context.Context, item *model.DokumenArsip) error
	DeleteDokumen(ctx context.Context, id uint) error

	// Audit Log
	CreateAuditLog(ctx context.Context, item *model.AuditLog) error
	FindAllAuditLog(ctx context.Context, req dto.AuditLogListRequest) ([]model.AuditLog, *dto.PaginatedResponse, error)
}

type operasionalRepositoryImpl struct {
	db *gorm.DB
}

func NewOperasionalRepository(db *gorm.DB) OperasionalRepository {
	return &operasionalRepositoryImpl{db: db}
}

// ─── Peminjaman ─────────────────────────────────────────────────

func (r *operasionalRepositoryImpl) SyncOverduePeminjaman(ctx context.Context, now time.Time) error {
	today := now.Format("2006-01-02")
	return r.db.WithContext(ctx).Model(&model.Peminjaman{}).
		Where("status = ? AND tanggal_kembali_aktual IS NULL AND tanggal_kembali_rencana < ?", string(model.StatusPeminjamanDipinjam), today).
		Update("status", string(model.StatusPeminjamanTerlambat)).Error
}

func (r *operasionalRepositoryImpl) FindAllPeminjaman(ctx context.Context, req dto.PeminjamanListRequest) ([]model.Peminjaman, *dto.PaginatedResponse, error) {
	var list []model.Peminjaman
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Peminjaman{})

	if s := req.ResolvedSearch(); s != "" {
		like := "%" + s + "%"
		query = query.Where("no_tiket LIKE ? OR peminjam LIKE ? OR catatan LIKE ?", like, like, like)
	}
	if req.Status != "" {
		query = query.Where("UPPER(status) = ?", strings.ToUpper(req.Status))
	}
	if req.UnitKerjaID > 0 {
		query = query.Where("unit_kerja_id = ?", req.UnitKerjaID)
	}
	if req.ArsipID > 0 {
		query = query.Where("arsip_id = ?", req.ArsipID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 || perPage > 500 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.
		Preload("Arsip").
		Preload("UnitKerja").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&list).Error; err != nil {
		return nil, nil, err
	}

	from, to := offset+1, offset+len(list)
	if len(list) == 0 {
		from, to = 0, 0
	}

	return list, &dto.PaginatedResponse{
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}, nil
}

func (r *operasionalRepositoryImpl) FindPeminjamanByID(ctx context.Context, id uint) (*model.Peminjaman, error) {
	var item model.Peminjaman
	if err := r.db.WithContext(ctx).
		Preload("Arsip").
		Preload("UnitKerja").
		First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *operasionalRepositoryImpl) CreatePeminjaman(ctx context.Context, item *model.Peminjaman) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *operasionalRepositoryImpl) UpdatePeminjaman(ctx context.Context, item *model.Peminjaman) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *operasionalRepositoryImpl) DeletePeminjaman(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Peminjaman{}, id).Error
}

func (r *operasionalRepositoryImpl) CountPeminjaman(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Peminjaman{}).Count(&count).Error
	return count, err
}

// ─── Penyusutan ─────────────────────────────────────────────────

func (r *operasionalRepositoryImpl) FindAllPenyusutan(ctx context.Context, req dto.PenyusutanListRequest) ([]model.Penyusutan, *dto.PaginatedResponse, error) {
	var list []model.Penyusutan
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Penyusutan{})

	if s := req.ResolvedSearch(); s != "" {
		like := "%" + s + "%"
		query = query.Where("no_pengajuan LIKE ? OR alasan LIKE ?", like, like)
	}
	if req.JenisUsul != "" {
		query = query.Where("UPPER(jenis_usul) = ?", strings.ToUpper(req.JenisUsul))
	}
	if req.Status != "" {
		query = query.Where("UPPER(status) = ?", strings.ToUpper(req.Status))
	}
	if req.ArsipID > 0 {
		query = query.Where("arsip_id = ?", req.ArsipID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 || perPage > 500 {
		perPage = 10
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.
		Preload("Arsip").
		Preload("Pengaju").
		Preload("Approver").
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&list).Error; err != nil {
		return nil, nil, err
	}

	from, to := offset+1, offset+len(list)
	if len(list) == 0 {
		from, to = 0, 0
	}

	return list, &dto.PaginatedResponse{
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}, nil
}

func (r *operasionalRepositoryImpl) FindPenyusutanByID(ctx context.Context, id uint) (*model.Penyusutan, error) {
	var item model.Penyusutan
	if err := r.db.WithContext(ctx).
		Preload("Arsip").
		Preload("Pengaju").
		Preload("Approver").
		First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *operasionalRepositoryImpl) CreatePenyusutan(ctx context.Context, item *model.Penyusutan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *operasionalRepositoryImpl) UpdatePenyusutan(ctx context.Context, item *model.Penyusutan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *operasionalRepositoryImpl) CountPenyusutan(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Penyusutan{}).Count(&count).Error
	return count, err
}

// ─── Dokumen Arsip ──────────────────────────────────────────────

func (r *operasionalRepositoryImpl) FindDokumenByArsipID(ctx context.Context, arsipID uint) ([]model.DokumenArsip, error) {
	var list []model.DokumenArsip
	err := r.db.WithContext(ctx).Where("arsip_id = ?", arsipID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *operasionalRepositoryImpl) FindDokumenByID(ctx context.Context, id uint) (*model.DokumenArsip, error) {
	var item model.DokumenArsip
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *operasionalRepositoryImpl) CreateDokumen(ctx context.Context, item *model.DokumenArsip) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *operasionalRepositoryImpl) DeleteDokumen(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.DokumenArsip{}, id).Error
}

// ─── Audit Log ──────────────────────────────────────────────────

func (r *operasionalRepositoryImpl) CreateAuditLog(ctx context.Context, item *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *operasionalRepositoryImpl) FindAllAuditLog(ctx context.Context, req dto.AuditLogListRequest) ([]model.AuditLog, *dto.PaginatedResponse, error) {
	var list []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})
	if req.Aksi != "" {
		query = query.Where("aksi = ?", req.Aksi)
	}
	if req.Entitas != "" {
		query = query.Where("entitas = ?", req.Entitas)
	}
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 || perPage > 200 {
		perPage = 20
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, nil, err
	}

	from, to := offset+1, offset+len(list)
	if len(list) == 0 {
		from, to = 0, 0
	}

	return list, &dto.PaginatedResponse{
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}, nil
}
