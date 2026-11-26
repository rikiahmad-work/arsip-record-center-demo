package repository

import (
	"context"
	"math"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MasterRepository mendefinisikan akses data untuk kode_klasifikasi, lokasi_simpan, unit_kerja, dan konfigurasi_sistem.
type MasterRepository interface {
	// Kode Klasifikasi
	FindAllKodeKlasifikasi(ctx context.Context, req dto.KodeKlasifikasiListRequest) ([]model.KodeKlasifikasi, *dto.PaginatedResponse, error)
	FindKodeKlasifikasiByID(ctx context.Context, id uint) (*model.KodeKlasifikasi, error)
	FindKodeKlasifikasiByKode(ctx context.Context, kode string) (*model.KodeKlasifikasi, error)
	CreateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error
	UpdateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error
	DeleteKodeKlasifikasi(ctx context.Context, id uint) error
	CountArsipPerKode(ctx context.Context) (map[string]int64, error)

	// Lokasi Simpan
	FindAllLokasiSimpan(ctx context.Context, req dto.LokasiSimpanListRequest) ([]model.LokasiSimpan, *dto.PaginatedResponse, error)
	FindLokasiSimpanByID(ctx context.Context, id uint) (*model.LokasiSimpan, error)
	CreateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error
	UpdateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error
	DeleteLokasiSimpan(ctx context.Context, id uint) error

	// Unit Kerja
	FindAllUnitKerja(ctx context.Context, req dto.UnitKerjaListRequest) ([]model.UnitKerja, *dto.PaginatedResponse, error)
	FindUnitKerjaByID(ctx context.Context, id uint) (*model.UnitKerja, error)
	FindUnitKerjaByKode(ctx context.Context, kode string) (*model.UnitKerja, error)
	CreateUnitKerja(ctx context.Context, item *model.UnitKerja) error
	UpdateUnitKerja(ctx context.Context, item *model.UnitKerja) error
	DeleteUnitKerja(ctx context.Context, id uint) error
	CountArsipPerUnitKerja(ctx context.Context) (map[uint]int64, error)

	// Konfigurasi Sistem
	GetAllKonfigurasi(ctx context.Context) (map[string]string, error)
	UpsertKonfigurasi(ctx context.Context, items map[string]string, updatedBy uint) error
}

type masterRepositoryImpl struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepositoryImpl{db: db}
}

// ─── Kode Klasifikasi ───────────────────────────────────────────

func (r *masterRepositoryImpl) FindAllKodeKlasifikasi(ctx context.Context, req dto.KodeKlasifikasiListRequest) ([]model.KodeKlasifikasi, *dto.PaginatedResponse, error) {
	var list []model.KodeKlasifikasi
	var total int64

	query := r.db.WithContext(ctx).Model(&model.KodeKlasifikasi{})

	search := req.Search
	if search == "" {
		search = req.Q
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("kode LIKE ? OR nama LIKE ?", like, like)
	}
	if req.SKKAD != "" {
		query = query.Where("skkad = ?", req.SKKAD)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 500 // default kembalikan seluruh master kode klasifikasi untuk dropdown/master
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.Order("id ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
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

func (r *masterRepositoryImpl) FindKodeKlasifikasiByID(ctx context.Context, id uint) (*model.KodeKlasifikasi, error) {
	var item model.KodeKlasifikasi
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *masterRepositoryImpl) FindKodeKlasifikasiByKode(ctx context.Context, kode string) (*model.KodeKlasifikasi, error) {
	var item model.KodeKlasifikasi
	if err := r.db.WithContext(ctx).Where("kode = ?", kode).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *masterRepositoryImpl) CreateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *masterRepositoryImpl) UpdateKodeKlasifikasi(ctx context.Context, item *model.KodeKlasifikasi) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *masterRepositoryImpl) DeleteKodeKlasifikasi(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.KodeKlasifikasi{}, id).Error
}

func (r *masterRepositoryImpl) CountArsipPerKode(ctx context.Context) (map[string]int64, error) {
	result := make(map[string]int64)
	type row struct {
		Kode  string
		Total int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).Model(&model.Arsip{}).
		Select("kode_klasifikasi as kode, COUNT(*) as total").
		Group("kode_klasifikasi").
		Scan(&rows).Error; err != nil {
		return result, err
	}
	for _, item := range rows {
		result[item.Kode] = item.Total
	}
	return result, nil
}

// ─── Lokasi Simpan ──────────────────────────────────────────────

func (r *masterRepositoryImpl) FindAllLokasiSimpan(ctx context.Context, req dto.LokasiSimpanListRequest) ([]model.LokasiSimpan, *dto.PaginatedResponse, error) {
	var list []model.LokasiSimpan
	var total int64

	query := r.db.WithContext(ctx).Model(&model.LokasiSimpan{})

	search := req.Search
	if search == "" {
		search = req.Q
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("nama LIKE ? OR gedung LIKE ?", like, like)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 100
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.Order("id ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
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

func (r *masterRepositoryImpl) FindLokasiSimpanByID(ctx context.Context, id uint) (*model.LokasiSimpan, error) {
	var item model.LokasiSimpan
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *masterRepositoryImpl) CreateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *masterRepositoryImpl) UpdateLokasiSimpan(ctx context.Context, item *model.LokasiSimpan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *masterRepositoryImpl) DeleteLokasiSimpan(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.LokasiSimpan{}, id).Error
}

// ─── Unit Kerja ─────────────────────────────────────────────────

func (r *masterRepositoryImpl) FindAllUnitKerja(ctx context.Context, req dto.UnitKerjaListRequest) ([]model.UnitKerja, *dto.PaginatedResponse, error) {
	var list []model.UnitKerja
	var total int64

	query := r.db.WithContext(ctx).Model(&model.UnitKerja{})

	search := req.Search
	if search == "" {
		search = req.Q
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("kode LIKE ? OR nama LIKE ?", like, like)
	}
	if req.IsActive != nil {
		query = query.Where("is_active = ?", *req.IsActive)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	page := req.Page
	if page < 1 {
		page = 1
	}
	perPage := req.PerPage
	if perPage < 1 {
		perPage = 100
	}

	offset := (page - 1) * perPage
	lastPage := int(math.Ceil(float64(total) / float64(perPage)))
	if lastPage < 1 {
		lastPage = 1
	}

	if err := query.Order("id ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
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

func (r *masterRepositoryImpl) FindUnitKerjaByID(ctx context.Context, id uint) (*model.UnitKerja, error) {
	var item model.UnitKerja
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *masterRepositoryImpl) FindUnitKerjaByKode(ctx context.Context, kode string) (*model.UnitKerja, error) {
	var item model.UnitKerja
	if err := r.db.WithContext(ctx).Where("kode = ?", kode).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *masterRepositoryImpl) CreateUnitKerja(ctx context.Context, item *model.UnitKerja) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *masterRepositoryImpl) UpdateUnitKerja(ctx context.Context, item *model.UnitKerja) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *masterRepositoryImpl) DeleteUnitKerja(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.UnitKerja{}, id).Error
}

func (r *masterRepositoryImpl) CountArsipPerUnitKerja(ctx context.Context) (map[uint]int64, error) {
	result := make(map[uint]int64)
	type row struct {
		UnitKerjaID uint
		Total       int64
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Table("unit_kerja u").
		Select("u.id AS unit_kerja_id, COUNT(a.id) AS total").
		Joins("LEFT JOIN arsip a ON a.deleted_at IS NULL AND (a.unit_kerja_id = u.id OR a.unor_id = u.kode OR a.pencipta_arsip = u.nama)").
		Where("u.deleted_at IS NULL").
		Group("u.id").
		Scan(&rows).Error; err != nil {
		return result, err
	}
	for _, item := range rows {
		result[item.UnitKerjaID] = item.Total
	}
	return result, nil
}

// ─── Konfigurasi Sistem ─────────────────────────────────────────

func (r *masterRepositoryImpl) GetAllKonfigurasi(ctx context.Context) (map[string]string, error) {
	var list []model.KonfigurasiSistem
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	res := make(map[string]string, len(list))
	for _, item := range list {
		if item.Value != nil {
			res[item.Key] = *item.Value
		} else {
			res[item.Key] = ""
		}
	}
	return res, nil
}

func (r *masterRepositoryImpl) UpsertKonfigurasi(ctx context.Context, items map[string]string, updatedBy uint) error {
	for k, v := range items {
		val := v
		rec := model.KonfigurasiSistem{
			Key:       k,
			Value:     &val,
			UpdatedBy: &updatedBy,
		}
		if err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "updated_by", "updated_at"}),
		}).Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}
