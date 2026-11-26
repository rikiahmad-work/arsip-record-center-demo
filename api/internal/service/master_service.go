package service

import (
	"context"
	"errors"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// MasterService mendefinisikan business logic untuk master data (kode klasifikasi, lokasi simpan, unit kerja, konfigurasi).
type MasterService interface {
	// Kode Klasifikasi
	ListKodeKlasifikasi(ctx context.Context, req dto.KodeKlasifikasiListRequest) ([]dto.KodeKlasifikasiResponse, *dto.PaginatedResponse, error)
	CreateKodeKlasifikasi(ctx context.Context, req dto.StoreKodeKlasifikasiRequest) (*dto.KodeKlasifikasiResponse, error)
	UpdateKodeKlasifikasi(ctx context.Context, id uint, req dto.UpdateKodeKlasifikasiRequest) (*dto.KodeKlasifikasiResponse, error)
	DeleteKodeKlasifikasi(ctx context.Context, id uint) error

	// Lokasi Simpan
	ListLokasiSimpan(ctx context.Context, req dto.LokasiSimpanListRequest) ([]dto.LokasiSimpanResponse, *dto.PaginatedResponse, error)
	CreateLokasiSimpan(ctx context.Context, req dto.StoreLokasiSimpanRequest) (*dto.LokasiSimpanResponse, error)
	UpdateLokasiSimpan(ctx context.Context, id uint, req dto.UpdateLokasiSimpanRequest) (*dto.LokasiSimpanResponse, error)
	DeleteLokasiSimpan(ctx context.Context, id uint) error

	// Unit Kerja
	ListUnitKerja(ctx context.Context, req dto.UnitKerjaListRequest) ([]dto.UnitKerjaResponse, *dto.PaginatedResponse, error)
	CreateUnitKerja(ctx context.Context, req dto.StoreUnitKerjaRequest) (*dto.UnitKerjaResponse, error)
	UpdateUnitKerja(ctx context.Context, id uint, req dto.UpdateUnitKerjaRequest) (*dto.UnitKerjaResponse, error)
	DeleteUnitKerja(ctx context.Context, id uint) error

	// Konfigurasi Sistem
	GetKonfigurasi(ctx context.Context) (map[string]string, error)
	UpdateKonfigurasi(ctx context.Context, req dto.UpdateKonfigurasiRequest, updatedBy uint) (map[string]string, error)
}

type masterServiceImpl struct {
	repo repository.MasterRepository
}

func NewMasterService(repo repository.MasterRepository) MasterService {
	return &masterServiceImpl{repo: repo}
}

// ─── Kode Klasifikasi ───────────────────────────────────────────

func (s *masterServiceImpl) ListKodeKlasifikasi(ctx context.Context, req dto.KodeKlasifikasiListRequest) ([]dto.KodeKlasifikasiResponse, *dto.PaginatedResponse, error) {
	items, meta, err := s.repo.FindAllKodeKlasifikasi(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	counts, _ := s.repo.CountArsipPerKode(ctx)

	res := make([]dto.KodeKlasifikasiResponse, 0, len(items))
	for _, item := range items {
		res = append(res, dto.KodeKlasifikasiResponse{
			ID:                  item.ID,
			Kode:                item.Kode,
			Nama:                item.Nama,
			JangkaSimpanInaktif: item.JangkaSimpanInaktif,
			NasibAkhir:          item.NasibAkhir,
			SKKAD:               item.SKKAD,
			TotalArsip:          counts[item.Kode],
			IsActive:            item.IsActive,
			CreatedAt:           item.CreatedAt,
			UpdatedAt:           item.UpdatedAt,
		})
	}
	return res, meta, nil
}

func (s *masterServiceImpl) CreateKodeKlasifikasi(ctx context.Context, req dto.StoreKodeKlasifikasiRequest) (*dto.KodeKlasifikasiResponse, error) {
	existing, _ := s.repo.FindKodeKlasifikasiByKode(ctx, req.Kode)
	if existing != nil {
		return nil, helper.NewAppError(409, "Kode klasifikasi sudah terdaftar", nil)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item := &model.KodeKlasifikasi{
		Kode:                req.Kode,
		Nama:                req.Nama,
		JangkaSimpanInaktif: req.JangkaSimpanInaktif,
		NasibAkhir:          req.NasibAkhir,
		SKKAD:               req.SKKAD,
		IsActive:            isActive,
	}

	if err := s.repo.CreateKodeKlasifikasi(ctx, item); err != nil {
		return nil, err
	}

	return &dto.KodeKlasifikasiResponse{
		ID:                  item.ID,
		Kode:                item.Kode,
		Nama:                item.Nama,
		JangkaSimpanInaktif: item.JangkaSimpanInaktif,
		NasibAkhir:          item.NasibAkhir,
		SKKAD:               item.SKKAD,
		TotalArsip:          0,
		IsActive:            item.IsActive,
		CreatedAt:           item.CreatedAt,
		UpdatedAt:           item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) UpdateKodeKlasifikasi(ctx context.Context, id uint, req dto.UpdateKodeKlasifikasiRequest) (*dto.KodeKlasifikasiResponse, error) {
	item, err := s.repo.FindKodeKlasifikasiByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if req.Kode != item.Kode {
		existing, _ := s.repo.FindKodeKlasifikasiByKode(ctx, req.Kode)
		if existing != nil {
			return nil, helper.NewAppError(409, "Kode klasifikasi sudah terdaftar", nil)
		}
	}

	item.Kode = req.Kode
	item.Nama = req.Nama
	item.JangkaSimpanInaktif = req.JangkaSimpanInaktif
	item.NasibAkhir = req.NasibAkhir
	item.SKKAD = req.SKKAD
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateKodeKlasifikasi(ctx, item); err != nil {
		return nil, err
	}

	return &dto.KodeKlasifikasiResponse{
		ID:                  item.ID,
		Kode:                item.Kode,
		Nama:                item.Nama,
		JangkaSimpanInaktif: item.JangkaSimpanInaktif,
		NasibAkhir:          item.NasibAkhir,
		SKKAD:               item.SKKAD,
		IsActive:            item.IsActive,
		CreatedAt:           item.CreatedAt,
		UpdatedAt:           item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) DeleteKodeKlasifikasi(ctx context.Context, id uint) error {
	_, err := s.repo.FindKodeKlasifikasiByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return s.repo.DeleteKodeKlasifikasi(ctx, id)
}

// ─── Lokasi Simpan ──────────────────────────────────────────────

func (s *masterServiceImpl) ListLokasiSimpan(ctx context.Context, req dto.LokasiSimpanListRequest) ([]dto.LokasiSimpanResponse, *dto.PaginatedResponse, error) {
	items, meta, err := s.repo.FindAllLokasiSimpan(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	res := make([]dto.LokasiSimpanResponse, 0, len(items))
	for _, item := range items {
		res = append(res, dto.LokasiSimpanResponse{
			ID:        item.ID,
			Nama:      item.Nama,
			Gedung:    item.Gedung,
			Lantai:    item.Lantai,
			NoRak:     item.NoRak,
			NoBoks:    item.NoBoks,
			NoFolder:  item.NoFolder,
			Kapasitas: item.Kapasitas,
			Terpakai:  item.Terpakai,
			IsActive:  item.IsActive,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return res, meta, nil
}

func (s *masterServiceImpl) CreateLokasiSimpan(ctx context.Context, req dto.StoreLokasiSimpanRequest) (*dto.LokasiSimpanResponse, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	terpakai := 0
	if req.Terpakai != nil {
		terpakai = *req.Terpakai
	}

	item := &model.LokasiSimpan{
		Nama:      req.Nama,
		Gedung:    req.Gedung,
		Lantai:    req.Lantai,
		NoRak:     req.NoRak,
		NoBoks:    req.NoBoks,
		NoFolder:  req.NoFolder,
		Kapasitas: req.Kapasitas,
		Terpakai:  terpakai,
		IsActive:  isActive,
	}

	if err := s.repo.CreateLokasiSimpan(ctx, item); err != nil {
		return nil, err
	}

	return &dto.LokasiSimpanResponse{
		ID:        item.ID,
		Nama:      item.Nama,
		Gedung:    item.Gedung,
		Lantai:    item.Lantai,
		NoRak:     item.NoRak,
		NoBoks:    item.NoBoks,
		NoFolder:  item.NoFolder,
		Kapasitas: item.Kapasitas,
		Terpakai:  item.Terpakai,
		IsActive:  item.IsActive,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) UpdateLokasiSimpan(ctx context.Context, id uint, req dto.UpdateLokasiSimpanRequest) (*dto.LokasiSimpanResponse, error) {
	item, err := s.repo.FindLokasiSimpanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	item.Nama = req.Nama
	item.Gedung = req.Gedung
	item.Lantai = req.Lantai
	item.NoRak = req.NoRak
	item.NoBoks = req.NoBoks
	item.NoFolder = req.NoFolder
	item.Kapasitas = req.Kapasitas
	if req.Terpakai != nil {
		item.Terpakai = *req.Terpakai
	}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateLokasiSimpan(ctx, item); err != nil {
		return nil, err
	}

	return &dto.LokasiSimpanResponse{
		ID:        item.ID,
		Nama:      item.Nama,
		Gedung:    item.Gedung,
		Lantai:    item.Lantai,
		NoRak:     item.NoRak,
		NoBoks:    item.NoBoks,
		NoFolder:  item.NoFolder,
		Kapasitas: item.Kapasitas,
		Terpakai:  item.Terpakai,
		IsActive:  item.IsActive,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) DeleteLokasiSimpan(ctx context.Context, id uint) error {
	_, err := s.repo.FindLokasiSimpanByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return s.repo.DeleteLokasiSimpan(ctx, id)
}

// ─── Unit Kerja ─────────────────────────────────────────────────

func (s *masterServiceImpl) ListUnitKerja(ctx context.Context, req dto.UnitKerjaListRequest) ([]dto.UnitKerjaResponse, *dto.PaginatedResponse, error) {
	items, meta, err := s.repo.FindAllUnitKerja(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	counts, _ := s.repo.CountArsipPerUnitKerja(ctx)

	res := make([]dto.UnitKerjaResponse, 0, len(items))
	for _, item := range items {
		res = append(res, dto.UnitKerjaResponse{
			ID:         item.ID,
			Kode:       item.Kode,
			Nama:       item.Nama,
			Eselon:     item.Eselon,
			TotalArsip: counts[item.ID],
			IsActive:   item.IsActive,
			CreatedAt:  item.CreatedAt,
			UpdatedAt:  item.UpdatedAt,
		})
	}
	return res, meta, nil
}

func (s *masterServiceImpl) CreateUnitKerja(ctx context.Context, req dto.StoreUnitKerjaRequest) (*dto.UnitKerjaResponse, error) {
	existing, _ := s.repo.FindUnitKerjaByKode(ctx, req.Kode)
	if existing != nil {
		return nil, helper.NewAppError(409, "Kode unit kerja sudah terdaftar", nil)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item := &model.UnitKerja{
		Kode:     req.Kode,
		Nama:     req.Nama,
		Eselon:   req.Eselon,
		IsActive: isActive,
	}

	if err := s.repo.CreateUnitKerja(ctx, item); err != nil {
		return nil, err
	}

	return &dto.UnitKerjaResponse{
		ID:         item.ID,
		Kode:       item.Kode,
		Nama:       item.Nama,
		Eselon:     item.Eselon,
		TotalArsip: 0,
		IsActive:   item.IsActive,
		CreatedAt:  item.CreatedAt,
		UpdatedAt:  item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) UpdateUnitKerja(ctx context.Context, id uint, req dto.UpdateUnitKerjaRequest) (*dto.UnitKerjaResponse, error) {
	item, err := s.repo.FindUnitKerjaByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if req.Kode != item.Kode {
		existing, _ := s.repo.FindUnitKerjaByKode(ctx, req.Kode)
		if existing != nil {
			return nil, helper.NewAppError(409, "Kode unit kerja sudah terdaftar", nil)
		}
	}

	item.Kode = req.Kode
	item.Nama = req.Nama
	item.Eselon = req.Eselon
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}

	if err := s.repo.UpdateUnitKerja(ctx, item); err != nil {
		return nil, err
	}

	return &dto.UnitKerjaResponse{
		ID:        item.ID,
		Kode:      item.Kode,
		Nama:      item.Nama,
		Eselon:    item.Eselon,
		IsActive:  item.IsActive,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *masterServiceImpl) DeleteUnitKerja(ctx context.Context, id uint) error {
	_, err := s.repo.FindUnitKerjaByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return s.repo.DeleteUnitKerja(ctx, id)
}

// ─── Konfigurasi Sistem ─────────────────────────────────────────

func (s *masterServiceImpl) GetKonfigurasi(ctx context.Context) (map[string]string, error) {
	return s.repo.GetAllKonfigurasi(ctx)
}

func (s *masterServiceImpl) UpdateKonfigurasi(ctx context.Context, req dto.UpdateKonfigurasiRequest, updatedBy uint) (map[string]string, error) {
	updates := make(map[string]string)
	for k, v := range req.Items {
		updates[k] = v
	}
	if req.AppName != "" {
		updates["app_name"] = req.AppName
	}
	if req.DefaultJangkaSimpan != "" {
		updates["default_jangka_simpan"] = req.DefaultJangkaSimpan
	}
	if req.DefaultNasibAkhir != "" {
		updates["default_nasib_akhir"] = req.DefaultNasibAkhir
	}
	if req.NotifikasiJatuhTempoHari != "" {
		updates["notifikasi_jatuh_tempo_hari"] = req.NotifikasiJatuhTempoHari
	}
	if req.LogoAplikasi != "" {
		updates["logo_aplikasi"] = req.LogoAplikasi
	}

	if err := s.repo.UpsertKonfigurasi(ctx, updates, updatedBy); err != nil {
		return nil, err
	}
	return s.repo.GetAllKonfigurasi(ctx)
}
