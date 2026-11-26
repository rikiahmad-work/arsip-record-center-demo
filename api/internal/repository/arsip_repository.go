package repository

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
)

// arsipSortWhitelist mencegah SQL injection melalui parameter sort_by.
var arsipSortWhitelist = map[string]string{
	"id":                    "id",
	"no_item_arsip":         "no_item_arsip",
	"no_folder":             "no_folder",
	"kode_klasifikasi":      "kode_klasifikasi",
	"jenis_arsip":           "jenis_arsip",
	"kurun_waktu":           "kurun_waktu",
	"tingkat_perkembangan":  "tingkat_perkembangan",
	"lokasi_simpan":         "lokasi_simpan",
	"no_boks":               "no_boks",
	"no_rak":                "no_rak",
	"jangka_simpan_inaktif": "jangka_simpan_inaktif",
	"nasib_akhir":           "nasib_akhir",
	"skkad":                 "skkad",
	"kategori_arsip":        "kategori_arsip",
	"pencipta_arsip":        "pencipta_arsip",
	"status":                "status",
	"created_at":            "created_at",
	"updated_at":            "updated_at",
}

// ArsipRepository mendefinisikan interface untuk akses database arsip.
type ArsipRepository interface {
	FindAll(ctx context.Context, req dto.ArsipListRequest) ([]model.Arsip, *dto.PaginatedResponse, error)
	FindByID(ctx context.Context, id uint) (*model.Arsip, error)
	FindByNoItem(ctx context.Context, noItem string) (*model.Arsip, error)
	Create(ctx context.Context, arsip *model.Arsip) error
	Update(ctx context.Context, arsip *model.Arsip) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) (*model.Arsip, error)
	CountByKategori(ctx context.Context) (map[string]int64, error)
	CountBySKKAD(ctx context.Context) (map[string]int64, error)
	CountByStatus(ctx context.Context) (draft, final int64, err error)
	CountTotal(ctx context.Context) (int64, error)
	CountByUnor(ctx context.Context, unorID string) (int64, error)
	NextNoItemByUnor(ctx context.Context, unorID string, unitKerjaID *uint, penciptaArsip string) (string, error)
}

type arsipRepositoryImpl struct {
	db *gorm.DB
}

// NewArsipRepository membuat instance ArsipRepository.
func NewArsipRepository(db *gorm.DB) ArsipRepository {
	return &arsipRepositoryImpl{db: db}
}

func (r *arsipRepositoryImpl) FindAll(ctx context.Context, req dto.ArsipListRequest) ([]model.Arsip, *dto.PaginatedResponse, error) {
	var arsipList []model.Arsip
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Arsip{})

	// Filter search (`q` atau `search`)
	if s := req.ResolvedSearch(); s != "" {
		search := "%" + s + "%"
		query = query.Where(
			"no_item_arsip LIKE ? OR uraian_informasi_arsip LIKE ? OR kode_klasifikasi LIKE ? OR pencipta_arsip LIKE ? OR no_boks LIKE ? OR no_rak LIKE ? OR no_folder LIKE ?",
			search, search, search, search, search, search, search,
		)
	}

	// Filter kode klasifikasi (`kode` atau `kode_klasifikasi`)
	if kode := req.ResolvedKode(); kode != "" {
		query = query.Where("kode_klasifikasi = ?", kode)
	}

	// Filter jenis arsip
	if req.JenisArsip != "" {
		query = query.Where("UPPER(jenis_arsip) = ?", strings.ToUpper(req.JenisArsip))
	}

	// Filter kategori arsip (`kategori` atau `kategori_arsip`)
	if kat := req.ResolvedKategori(); kat != "" {
		query = query.Where("UPPER(kategori_arsip) = ?", strings.ToUpper(kat))
	}

	// Filter SKKAD
	if req.SKKAD != "" {
		query = query.Where("UPPER(skkad) = ?", strings.ToUpper(req.SKKAD))
	}

	// Filter tahun (`kurun_waktu`)
	if req.Tahun != "" {
		query = query.Where("kurun_waktu LIKE ?", "%"+req.Tahun+"%")
	}

	// Filter unit kerja
	if req.UnitKerjaID > 0 {
		query = query.Where("unit_kerja_id = ?", req.UnitKerjaID)
	} else if req.UnitKerja != "" {
		query = query.Where("pencipta_arsip LIKE ? OR unor_id = ?", "%"+req.UnitKerja+"%", req.UnitKerja)
	}

	// Filter lokasi simpan
	if req.LokasiSimpanID > 0 {
		query = query.Where("lokasi_simpan_id = ?", req.LokasiSimpanID)
	} else if req.LokasiSimpan != "" {
		query = query.Where("lokasi_simpan = ?", req.LokasiSimpan)
	}

	// Filter status
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// Hitung total sebelum pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Sort whitelist
	sortCol, ok := arsipSortWhitelist[strings.ToLower(req.SortBy)]
	if !ok {
		sortCol = "created_at"
	}
	sortOrder := strings.ToLower(req.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
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
	from := offset + 1

	// Query dengan preload relasi
	if err := query.
		Preload("PetugasInput").
		Preload("LokasiSimpanRef").
		Preload("UnitKerjaRef").
		Preload("PeminjamanList").
		Preload("PenyusutanList").
		Order(sortCol + " " + sortOrder).
		Offset(offset).
		Limit(perPage).
		Find(&arsipList).Error; err != nil {
		return nil, nil, err
	}

	to := offset + len(arsipList)
	if len(arsipList) == 0 {
		from = 0
		to = 0
	}

	meta := &dto.PaginatedResponse{
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     perPage,
		Total:       total,
		From:        from,
		To:          to,
	}

	return arsipList, meta, nil
}

func (r *arsipRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.Arsip, error) {
	var arsip model.Arsip
	err := r.db.WithContext(ctx).
		Preload("PetugasInput").
		Preload("LokasiSimpanRef").
		Preload("UnitKerjaRef").
		Preload("DokumenList").
		Preload("PeminjamanList").
		Preload("PenyusutanList").
		First(&arsip, id).Error
	if err != nil {
		return nil, err
	}
	return &arsip, nil
}

func (r *arsipRepositoryImpl) FindByNoItem(ctx context.Context, noItem string) (*model.Arsip, error) {
	var arsip model.Arsip
	err := r.db.WithContext(ctx).Where("no_item_arsip = ?", noItem).First(&arsip).Error
	if err != nil {
		return nil, err
	}
	return &arsip, nil
}

func (r *arsipRepositoryImpl) Create(ctx context.Context, arsip *model.Arsip) error {
	return r.db.WithContext(ctx).Create(arsip).Error
}

func (r *arsipRepositoryImpl) Update(ctx context.Context, arsip *model.Arsip) error {
	return r.db.WithContext(ctx).Save(arsip).Error
}

func (r *arsipRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Arsip{}, id).Error
}

func (r *arsipRepositoryImpl) Restore(ctx context.Context, id uint) (*model.Arsip, error) {
	var arsip model.Arsip
	err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&arsip).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Unscoped().Model(&arsip).Update("deleted_at", nil).Error
	if err != nil {
		return nil, err
	}

	return &arsip, nil
}

func (r *arsipRepositoryImpl) CountByKategori(ctx context.Context) (map[string]int64, error) {
	result := make(map[string]int64)

	type CountResult struct {
		KategoriArsip string
		Count         int64
	}

	var counts []CountResult
	err := r.db.WithContext(ctx).Model(&model.Arsip{}).
		Select("kategori_arsip, COUNT(*) as count").
		Group("kategori_arsip").
		Find(&counts).Error
	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		result[strings.ToUpper(c.KategoriArsip)] = c.Count
	}

	return result, nil
}

func (r *arsipRepositoryImpl) CountBySKKAD(ctx context.Context) (map[string]int64, error) {
	result := map[string]int64{
		string(model.SKKADBiasa):    0,
		string(model.SKKADTerbatas): 0,
		string(model.SKKADRahasia):  0,
	}

	type CountResult struct {
		SKKAD string
		Count int64
	}

	var counts []CountResult
	err := r.db.WithContext(ctx).Model(&model.Arsip{}).
		Select("skkad, COUNT(*) as count").
		Group("skkad").
		Find(&counts).Error
	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		result[c.SKKAD] = c.Count
	}

	return result, nil
}

func (r *arsipRepositoryImpl) CountByStatus(ctx context.Context) (draft, final int64, err error) {
	r.db.WithContext(ctx).Model(&model.Arsip{}).Where("status = ?", "draft").Count(&draft)
	r.db.WithContext(ctx).Model(&model.Arsip{}).Where("status = ?", "final").Count(&final)
	return draft, final, nil
}

func (r *arsipRepositoryImpl) CountTotal(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Arsip{}).Count(&count).Error
	return count, err
}

func (r *arsipRepositoryImpl) CountByUnor(ctx context.Context, unorID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Arsip{}).
		Where("unor_id = ?", unorID).
		Count(&count).Error
	return count, err
}

// NextNoItemByUnor menghasilkan nomor urut item arsip biasa ("1", "2", "3", ...)
// yang mengurut per Unit Kerja / Pencipta Arsip.
func (r *arsipRepositoryImpl) NextNoItemByUnor(ctx context.Context, unorID string, unitKerjaID *uint, penciptaArsip string) (string, error) {
	query := r.db.WithContext(ctx).Model(&model.Arsip{})
	switch {
	case unitKerjaID != nil && *unitKerjaID > 0 && penciptaArsip != "":
		query = query.Where("unor_id = ? OR unit_kerja_id = ? OR pencipta_arsip = ?", unorID, *unitKerjaID, penciptaArsip)
	case unitKerjaID != nil && *unitKerjaID > 0:
		query = query.Where("unor_id = ? OR unit_kerja_id = ?", unorID, *unitKerjaID)
	case penciptaArsip != "":
		query = query.Where("unor_id = ? OR pencipta_arsip = ?", unorID, penciptaArsip)
	default:
		query = query.Where("unor_id = ?", unorID)
	}

	var activeNoItems []string
	if err := query.Pluck("no_item_arsip", &activeNoItems).Error; err != nil {
		return "", err
	}

	usedActive := make(map[string]bool, len(activeNoItems))
	maxNum := 0
	for _, raw := range activeNoItems {
		trimmed := strings.TrimSpace(raw)
		usedActive[trimmed] = true
		if n, err := strconv.Atoi(trimmed); err == nil && n > maxNum {
			maxNum = n
		}
	}

	nextNum := maxNum + 1
	if len(activeNoItems)+1 > nextNum {
		nextNum = len(activeNoItems) + 1
	}
	for usedActive[strconv.Itoa(nextNum)] {
		nextNum++
	}
	nextStr := strconv.Itoa(nextNum)

	// Pastikan tidak ada baris soft-deleted yang menahan UNIQUE KEY (unor_id, no_item_arsip)
	var softDeleted []model.Arsip
	if err := r.db.WithContext(ctx).Unscoped().
		Where("unor_id = ? AND no_item_arsip = ? AND deleted_at IS NOT NULL", unorID, nextStr).
		Find(&softDeleted).Error; err == nil {
		for _, delRow := range softDeleted {
			_ = r.db.WithContext(ctx).Unscoped().
				Model(&model.Arsip{}).
				Where("id = ?", delRow.ID).
				Update("no_item_arsip", fmt.Sprintf("%s_del_%d", nextStr, delRow.ID)).Error
		}
	}

	return nextStr, nil
}
