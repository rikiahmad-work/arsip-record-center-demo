package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// ─── Dummy UNOR (simulasi SIMPEG untuk backward-compatibility) ──
var dummyUNOR = map[string]string{
	"a1b2c3d4-e5f6-7890-abcd-ef1234567890": "Biro Umum & SDM",
	"b2c3d4e5-f6a7-8901-bcde-f12345678901": "Biro Data dan Statistik",
	"c3d4e5f6-a7b8-9012-cdef-123456789012": "Biro Hukum",
	"d4e5f6a7-b8c9-0123-defa-234567890123": "Biro Perencanaan",
	"e5f6a7b8-c9d0-1234-efab-345678901234": "Sekretariat Jenderal",
	"f6a7b8c9-d0e1-2345-fabc-456789012345": "Inspektorat",
	"a1a2a3a4-b5c6-7890-abcd-ef1234567890": "Biro Data dan Informasi",
}

// ArsipService mendefinisikan interface untuk business logic arsip.
type ArsipService interface {
	GetAll(ctx context.Context, req dto.ArsipListRequest) ([]dto.ArsipResponse, *dto.PaginatedResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.ArsipResponse, error)
	Create(ctx context.Context, req dto.StoreArsipRequest, petugasInputID uint, callerRole string) (*dto.ArsipResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateArsipRequest, callerRole string) (*dto.ArsipResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) (*dto.ArsipResponse, error)
	GetStats(ctx context.Context) (*dto.ArsipStatsResponse, error)
	GetDaftarUNOR() map[string]string
	ExportCSV(ctx context.Context, req dto.ArsipListRequest) ([]byte, error)
	ImportCSV(ctx context.Context, reader io.Reader, petugasInputID uint, callerRole string) (*dto.ArsipImportResultResponse, error)
}

type arsipServiceImpl struct {
	arsipRepo  repository.ArsipRepository
	masterRepo repository.MasterRepository
}

// NewArsipService membuat instance ArsipService dengan injeksi ArsipRepository dan MasterRepository.
func NewArsipService(arsipRepo repository.ArsipRepository, masterRepos ...repository.MasterRepository) ArsipService {
	var mRepo repository.MasterRepository
	if len(masterRepos) > 0 {
		mRepo = masterRepos[0]
	}
	return &arsipServiceImpl{
		arsipRepo:  arsipRepo,
		masterRepo: mRepo,
	}
}

func (s *arsipServiceImpl) GetDaftarUNOR() map[string]string {
	result := make(map[string]string, len(dummyUNOR))
	for k, v := range dummyUNOR {
		result[k] = v
	}
	return result
}

func resolveUnor(unorID string) (string, bool) {
	label, ok := dummyUNOR[unorID]
	return label, ok
}

// normalizeJenisArsip menormalisasi input jenis_arsip ke nilai kanonik ('Konvensional' / 'Elektronik').
func normalizeJenisArsip(raw string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "KONVENSIONAL":
		return string(model.JenisArsipKonvensional), nil
	case "ELEKTRONIK", "DIGITAL":
		return string(model.JenisArsipElektronik), nil
	default:
		return "", helper.NewAppError(422, fmt.Sprintf("Jenis arsip tidak valid: %s", raw), nil)
	}
}

// normalizeTingkatPerkembangan menormalisasi input tingkat_perkembangan ke nilai kanonik.
func normalizeTingkatPerkembangan(raw string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "ASLI":
		return string(model.TingkatPerkembanganAsli), nil
	case "COPY":
		return string(model.TingkatPerkembanganCopy), nil
	case "SALINAN":
		return string(model.TingkatPerkembanganSalinan), nil
	case "TEMBUSAN":
		return string(model.TingkatPerkembanganTembusan), nil
	default:
		return "", helper.NewAppError(422, fmt.Sprintf("Tingkat perkembangan tidak valid: %s", raw), nil)
	}
}

// normalizeKategoriArsip menormalisasi kategori_arsip ke nilai kanonik Record Center.
func normalizeKategoriArsip(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return string(model.KategoriInaktifVal), nil
	}
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "INAKTIF", "AKTIF":
		return string(model.KategoriInaktifVal), nil
	case "VITAL":
		return string(model.KategoriVitalVal), nil
	case "TERJAGA":
		return string(model.KategoriTerjagaVal), nil
	case "USUL MUSNAH":
		return string(model.KategoriUsulMusnahVal), nil
	case "USUL SERAH":
		return string(model.KategoriUsulSerahVal), nil
	default:
		return "", helper.NewAppError(422, fmt.Sprintf("Kategori arsip tidak valid: %s", raw), nil)
	}
}

// resolveJRAFromMaster menerapkan aturan Server-Side Authority:
// mengambil jangka_simpan_inaktif, nasib_akhir, dan skkad dari tabel master kode_klasifikasi.
// Mengembalikan 422 jika kode klasifikasi tidak ditemukan.
func (s *arsipServiceImpl) resolveJRAFromMaster(ctx context.Context, kode, requestedSKKAD, callerRole string) (jangkaSimpan string, nasibAkhir *string, skkad string, err error) {
	if s.masterRepo == nil {
		return "3 Tahun", strPtrVal(string(model.NasibAkhirMusnah)), string(model.SKKADBiasa), nil
	}

	master, err := s.masterRepo.FindKodeKlasifikasiByKode(ctx, strings.TrimSpace(kode))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, "", helper.NewAppError(422, fmt.Sprintf("Kode klasifikasi '%s' tidak ditemukan di master data", kode), nil)
		}
		return "", nil, "", err
	}

	if master.JangkaSimpanInaktif != nil && *master.JangkaSimpanInaktif != "" {
		jangkaSimpan = *master.JangkaSimpanInaktif
	} else {
		jangkaSimpan = "3 Tahun"
	}

	if master.NasibAkhir != nil && *master.NasibAkhir != "" {
		nasibVal := *master.NasibAkhir
		nasibAkhir = &nasibVal
	} else {
		nasibAkhir = strPtrVal(string(model.NasibAkhirMusnah))
	}

	skkad = master.SKKAD
	if skkad == "" {
		skkad = string(model.SKKADBiasa)
	}

	// Hanya Administrator yang boleh melakukan override SKKAD dengan nilai kanonik yang valid
	if model.NormalizeRole(callerRole) == model.RoleAdministrator && requestedSKKAD != "" {
		if model.SKKAD(requestedSKKAD).IsValid() {
			skkad = requestedSKKAD
		}
	}

	return jangkaSimpan, nasibAkhir, skkad, nil
}

// resolveUnitKerjaAndLokasi memetakan unit_kerja_id / unor_id dan lokasi_simpan_id / lokasi_simpan.
func (s *arsipServiceImpl) resolveUnitKerjaAndLokasi(ctx context.Context, unitKerjaID *uint, unorID, pencipta string, lokasiID *uint, lokasiNama string) (resolvedUnitID *uint, resolvedUnorID, resolvedPencipta string, resolvedLokasiID *uint, resolvedLokasiNama string, err error) {
	resolvedUnitID = unitKerjaID
	resolvedUnorID = unorID
	resolvedPencipta = strings.TrimSpace(pencipta)

	if s.masterRepo != nil && unitKerjaID != nil && *unitKerjaID > 0 {
		uk, errUK := s.masterRepo.FindUnitKerjaByID(ctx, *unitKerjaID)
		if errUK != nil {
			if errors.Is(errUK, gorm.ErrRecordNotFound) {
				return nil, "", "", nil, "", helper.NewAppError(422, "Unit kerja tidak ditemukan", nil)
			}
			return nil, "", "", nil, "", errUK
		}
		resolvedPencipta = uk.Nama
		if resolvedUnorID == "" {
			resolvedUnorID = uk.Kode
		}
	} else if resolvedUnorID != "" {
		if label, ok := resolveUnor(resolvedUnorID); ok {
			resolvedPencipta = label
		} else if resolvedPencipta == "" {
			resolvedPencipta = resolvedUnorID
		}
	} else {
		resolvedUnorID = "BU.01"
		if resolvedPencipta == "" {
			resolvedPencipta = "Biro Umum & SDM"
		}
	}

	resolvedLokasiID = lokasiID
	resolvedLokasiNama = strings.TrimSpace(lokasiNama)
	if s.masterRepo != nil && lokasiID != nil && *lokasiID > 0 {
		lok, errLok := s.masterRepo.FindLokasiSimpanByID(ctx, *lokasiID)
		if errLok != nil {
			if errors.Is(errLok, gorm.ErrRecordNotFound) {
				return nil, "", "", nil, "", helper.NewAppError(422, "Lokasi simpan tidak ditemukan", nil)
			}
			return nil, "", "", nil, "", errLok
		}
		resolvedLokasiNama = lok.Nama
	} else if resolvedLokasiNama == "" {
		resolvedLokasiNama = "Record Center"
	}

	return resolvedUnitID, resolvedUnorID, resolvedPencipta, resolvedLokasiID, resolvedLokasiNama, nil
}

func (s *arsipServiceImpl) GetAll(ctx context.Context, req dto.ArsipListRequest) ([]dto.ArsipResponse, *dto.PaginatedResponse, error) {
	arsipList, meta, err := s.arsipRepo.FindAll(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	responses := make([]dto.ArsipResponse, 0, len(arsipList))
	for _, arsip := range arsipList {
		responses = append(responses, s.toResponse(arsip))
	}

	return responses, meta, nil
}

func (s *arsipServiceImpl) GetByID(ctx context.Context, id uint) (*dto.ArsipResponse, error) {
	arsip, err := s.arsipRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	resp := s.toResponse(*arsip)
	return &resp, nil
}

func (s *arsipServiceImpl) Create(ctx context.Context, req dto.StoreArsipRequest, petugasInputID uint, callerRole string) (*dto.ArsipResponse, error) {
	uraian := strings.TrimSpace(req.ResolvedUraian())
	if uraian == "" {
		return nil, helper.NewAppError(422, "Uraian informasi arsip wajib diisi", nil)
	}

	jenis, err := normalizeJenisArsip(req.JenisArsip)
	if err != nil {
		return nil, err
	}
	tingkat, err := normalizeTingkatPerkembangan(req.TingkatPerkembangan)
	if err != nil {
		return nil, err
	}
	kategori, err := normalizeKategoriArsip(req.KategoriArsip)
	if err != nil {
		return nil, err
	}

	// Auto-fill JRA & SKKAD dari master kode_klasifikasi
	jangkaSimpan, nasibAkhir, skkad, err := s.resolveJRAFromMaster(ctx, req.KodeKlasifikasi, req.SKKAD, callerRole)
	if err != nil {
		return nil, err
	}

	unitID, unorID, pencipta, lokasiID, lokasiNama, err := s.resolveUnitKerjaAndLokasi(
		ctx, req.UnitKerjaID, req.UnorID, req.PenciptaArsip, req.LokasiSimpanID, req.LokasiSimpan,
	)
	if err != nil {
		return nil, err
	}

	// Auto-generate no_item_arsip jika tidak diisi (format: ARS/YYYY/KODE/NNN)
	noItem := strings.TrimSpace(req.NoItemArsip)
	if noItem == "" {
		count, err := s.arsipRepo.CountTotal(ctx)
		if err != nil {
			return nil, err
		}
		noItem = fmt.Sprintf("ARS/%d/%s/%03d", time.Now().Year(), req.KodeKlasifikasi, count+1)
	} else {
		existing, _ := s.arsipRepo.FindByNoItem(ctx, noItem)
		if existing != nil && existing.UnorID == unorID {
			return nil, helper.NewAppError(409, "No item arsip sudah terdaftar", nil)
		}
	}

	status := req.Status
	if status == "" {
		status = string(model.StatusArsipFinal)
	}

	var nomorDefinitif *string
	if req.NomorDefinitif != nil && *req.NomorDefinitif != "" {
		nomorDefinitif = req.NomorDefinitif
	}

	var legacyNasib *string
	if nasibAkhir != nil {
		upper := strings.ToUpper(*nasibAkhir)
		legacyNasib = &upper
	}

	arsip := &model.Arsip{
		NoItemArsip:            noItem,
		NoFolder:               req.NoFolder,
		KodeKlasifikasi:        strings.TrimSpace(req.KodeKlasifikasi),
		UraianInformasiArsip:   uraian,
		JenisArsip:             jenis,
		KurunWaktu:             req.KurunWaktu,
		TingkatPerkembangan:    tingkat,
		Jumlah:                 req.Jumlah,
		Keterangan:             req.Keterangan,
		NomorDefinitif:         nomorDefinitif,
		LokasiSimpanID:         lokasiID,
		LokasiSimpan:           lokasiNama,
		NoBoks:                 req.NoBoks,
		NoRak:                  req.NoRak,
		JangkaSimpanInaktif:    jangkaSimpan,
		JangkaSimpanNasibAkhir: legacyNasib,
		NasibAkhir:             nasibAkhir,
		SKKAD:                  skkad,
		KategoriArsip:          kategori,
		UnitKerjaID:            unitID,
		UnorID:                 unorID,
		PenciptaArsip:          pencipta,
		PetugasInputID:         petugasInputID,
		Status:                 status,
	}

	if err := s.arsipRepo.Create(ctx, arsip); err != nil {
		return nil, err
	}

	created, err := s.arsipRepo.FindByID(ctx, arsip.ID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(*created)
	return &resp, nil
}

func (s *arsipServiceImpl) Update(ctx context.Context, id uint, req dto.UpdateArsipRequest, callerRole string) (*dto.ArsipResponse, error) {
	arsip, err := s.arsipRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	uraian := strings.TrimSpace(req.ResolvedUraian())
	if uraian == "" {
		return nil, helper.NewAppError(422, "Uraian informasi arsip wajib diisi", nil)
	}

	jenis, err := normalizeJenisArsip(req.JenisArsip)
	if err != nil {
		return nil, err
	}
	tingkat, err := normalizeTingkatPerkembangan(req.TingkatPerkembangan)
	if err != nil {
		return nil, err
	}
	kategori, err := normalizeKategoriArsip(req.KategoriArsip)
	if err != nil {
		return nil, err
	}

	// Auto-fill JRA & SKKAD dari master kode_klasifikasi
	jangkaSimpan, nasibAkhir, skkad, err := s.resolveJRAFromMaster(ctx, req.KodeKlasifikasi, req.SKKAD, callerRole)
	if err != nil {
		return nil, err
	}

	unitID, unorID, pencipta, lokasiID, lokasiNama, err := s.resolveUnitKerjaAndLokasi(
		ctx, req.UnitKerjaID, req.UnorID, req.PenciptaArsip, req.LokasiSimpanID, req.LokasiSimpan,
	)
	if err != nil {
		return nil, err
	}

	var nomorDefinitif *string
	if req.NomorDefinitif != nil && *req.NomorDefinitif != "" {
		nomorDefinitif = req.NomorDefinitif
	}

	var legacyNasib *string
	if nasibAkhir != nil {
		upper := strings.ToUpper(*nasibAkhir)
		legacyNasib = &upper
	}

	if req.NoItemArsip != "" {
		arsip.NoItemArsip = req.NoItemArsip
	}
	if req.NoFolder != nil {
		arsip.NoFolder = req.NoFolder
	}
	arsip.KodeKlasifikasi = strings.TrimSpace(req.KodeKlasifikasi)
	arsip.UraianInformasiArsip = uraian
	arsip.JenisArsip = jenis
	arsip.KurunWaktu = req.KurunWaktu
	arsip.TingkatPerkembangan = tingkat
	arsip.Jumlah = req.Jumlah
	arsip.Keterangan = req.Keterangan
	arsip.NomorDefinitif = nomorDefinitif
	arsip.LokasiSimpanID = lokasiID
	arsip.LokasiSimpan = lokasiNama
	arsip.NoBoks = req.NoBoks
	arsip.NoRak = req.NoRak
	arsip.JangkaSimpanInaktif = jangkaSimpan
	arsip.JangkaSimpanNasibAkhir = legacyNasib
	arsip.NasibAkhir = nasibAkhir
	arsip.SKKAD = skkad
	arsip.KategoriArsip = kategori
	arsip.UnitKerjaID = unitID
	if unorID != "" {
		arsip.UnorID = unorID
	}
	if pencipta != "" {
		arsip.PenciptaArsip = pencipta
	}
	if req.Status != "" {
		arsip.Status = req.Status
	}

	if err := s.arsipRepo.Update(ctx, arsip); err != nil {
		return nil, err
	}

	updated, err := s.arsipRepo.FindByID(ctx, arsip.ID)
	if err != nil {
		return nil, err
	}

	resp := s.toResponse(*updated)
	return &resp, nil
}

func (s *arsipServiceImpl) Delete(ctx context.Context, id uint) error {
	_, err := s.arsipRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}
	return s.arsipRepo.Delete(ctx, id)
}

func (s *arsipServiceImpl) Restore(ctx context.Context, id uint) (*dto.ArsipResponse, error) {
	arsip, err := s.arsipRepo.Restore(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	resp := s.toResponse(*arsip)
	return &resp, nil
}

func (s *arsipServiceImpl) GetStats(ctx context.Context) (*dto.ArsipStatsResponse, error) {
	total, err := s.arsipRepo.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	draft, final, err := s.arsipRepo.CountByStatus(ctx)
	if err != nil {
		return nil, err
	}

	byKategori, err := s.arsipRepo.CountByKategori(ctx)
	if err != nil {
		return nil, err
	}

	bySKKAD, err := s.arsipRepo.CountBySKKAD(ctx)
	if err != nil {
		return nil, err
	}

	inaktif := byKategori["INAKTIF"] + byKategori["AKTIF"]
	vital := byKategori["VITAL"]
	terjaga := byKategori["TERJAGA"]
	usulMusnah := byKategori["USUL MUSNAH"]
	usulSerah := byKategori["USUL SERAH"]

	return &dto.ArsipStatsResponse{
		Total:           total,
		ArsipInaktif:    inaktif,
		ArsipVital:      vital,
		ArsipTerjaga:    terjaga,
		ArsipUsulMusnah: usulMusnah,
		ArsipUsulSerah:  usulSerah,
		ArsipAktif:      vital,
		ArsipDraft:      draft,
		ArsipFinal:      final,
		BySKKAD:         bySKKAD,
	}, nil
}

// ExportCSV mengekspor daftar arsip sesuai filter menjadi format CSV 19 kolom DAI (+ No. Folder).
func (s *arsipServiceImpl) ExportCSV(ctx context.Context, req dto.ArsipListRequest) ([]byte, error) {
	req.Page = 1
	req.PerPage = 500
	items, _, err := s.arsipRepo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{
		"NO", "NO_ITEM_ARSIP", "NO_FOLDER", "KODE_KLASIFIKASI", "URAIAN_INFORMASI_ARSIP",
		"JENIS_ARSIP", "KURUN_WAKTU", "TINGKAT_PERKEMBANGAN", "JUMLAH", "KETERANGAN",
		"NOMOR_DEFINITIF", "LOKASI_SIMPAN", "NO_BOKS", "NO_RAK",
		"JANGKA_SIMPAN_INAKTIF", "NASIB_AKHIR", "SKKAD", "KATEGORI_ARSIP", "PENCIPTA_ARSIP",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for i, a := range items {
		noFolder := ""
		if a.NoFolder != nil {
			noFolder = *a.NoFolder
		}
		noDef := ""
		if a.NomorDefinitif != nil {
			noDef = *a.NomorDefinitif
		}
		nasib := ""
		if a.NasibAkhir != nil {
			nasib = *a.NasibAkhir
		} else if a.JangkaSimpanNasibAkhir != nil {
			nasib = *a.JangkaSimpanNasibAkhir
		}

		row := []string{
			fmt.Sprintf("%d", i+1),
			a.NoItemArsip,
			noFolder,
			a.KodeKlasifikasi,
			a.UraianInformasiArsip,
			a.JenisArsip,
			a.KurunWaktu,
			a.TingkatPerkembangan,
			a.Jumlah,
			a.Keterangan,
			noDef,
			a.LokasiSimpan,
			a.NoBoks,
			a.NoRak,
			a.JangkaSimpanInaktif,
			nasib,
			a.SKKAD,
			a.KategoriArsip,
			a.PenciptaArsip,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// ImportCSV mengimpor data arsip dari file CSV dengan auto-fill JRA & SKKAD dari master kode_klasifikasi.
func (s *arsipServiceImpl) ImportCSV(ctx context.Context, reader io.Reader, petugasInputID uint, callerRole string) (*dto.ArsipImportResultResponse, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, helper.NewAppError(422, "Format file CSV tidak valid", nil)
	}
	if len(records) < 2 {
		return nil, helper.NewAppError(422, "File import kosong atau hanya berisi header", nil)
	}

	result := &dto.ArsipImportResultResponse{
		TotalRows: len(records) - 1,
	}

	for idx, row := range records[1:] {
		lineNum := idx + 2
		if len(row) < 5 {
			result.FailedRows++
			result.Errors = append(result.Errors, fmt.Sprintf("Baris %d: jumlah kolom kurang dari minimal", lineNum))
			continue
		}

		getCol := func(i int) string {
			if i < len(row) {
				return strings.TrimSpace(row[i])
			}
			return ""
		}

		noFolder := getCol(2)
		var noFolderPtr *string
		if noFolder != "" {
			noFolderPtr = &noFolder
		}
		noDef := getCol(10)
		var noDefPtr *string
		if noDef != "" {
			noDefPtr = &noDef
		}

		jenis := getCol(5)
		if jenis == "" {
			jenis = string(model.JenisArsipKonvensional)
		}
		kurun := getCol(6)
		if kurun == "" {
			kurun = fmt.Sprintf("%d", time.Now().Year())
		}
		tingkat := getCol(7)
		if tingkat == "" {
			tingkat = string(model.TingkatPerkembanganAsli)
		}
		jumlah := getCol(8)
		if jumlah == "" {
			jumlah = "1 Berkas"
		}

		storeReq := dto.StoreArsipRequest{
			NoItemArsip:         getCol(1),
			NoFolder:            noFolderPtr,
			KodeKlasifikasi:     getCol(3),
			UraianInformasi:     getCol(4),
			JenisArsip:          jenis,
			KurunWaktu:          kurun,
			TingkatPerkembangan: tingkat,
			Jumlah:              jumlah,
			Keterangan:          getCol(9),
			NomorDefinitif:      noDefPtr,
			LokasiSimpan:        getCol(11),
			NoBoks:              getCol(12),
			NoRak:               getCol(13),
			SKKAD:               getCol(16),
			KategoriArsip:       getCol(17),
			PenciptaArsip:       getCol(18),
			Status:              string(model.StatusArsipFinal),
		}

		if _, err := s.Create(ctx, storeReq, petugasInputID, callerRole); err != nil {
			result.FailedRows++
			result.Errors = append(result.Errors, fmt.Sprintf("Baris %d (%s): %v", lineNum, storeReq.KodeKlasifikasi, err.Error()))
			continue
		}
		result.ImportedRows++
	}

	return result, nil
}

func (s *arsipServiceImpl) toResponse(arsip model.Arsip) dto.ArsipResponse {
	petugasNama := ""
	if arsip.PetugasInput.ID != 0 {
		petugasNama = arsip.PetugasInput.Name
	}

	lokasiNama := arsip.LokasiSimpan
	if arsip.LokasiSimpanRef != nil && arsip.LokasiSimpanRef.Nama != "" {
		lokasiNama = arsip.LokasiSimpanRef.Nama
	}

	unitNama := arsip.PenciptaArsip
	if arsip.UnitKerjaRef != nil && arsip.UnitKerjaRef.Nama != "" {
		unitNama = arsip.UnitKerjaRef.Nama
	}

	nasib := arsip.NasibAkhir
	if nasib == nil && arsip.JangkaSimpanNasibAkhir != nil {
		val := *arsip.JangkaSimpanNasibAkhir
		if strings.EqualFold(val, "MUSNAH") {
			v := string(model.NasibAkhirMusnah)
			nasib = &v
		} else if strings.EqualFold(val, "PERMANEN") {
			v := string(model.NasibAkhirPermanen)
			nasib = &v
		}
	}

	return dto.ArsipResponse{
		ID:                     arsip.ID,
		NoItemArsip:            arsip.NoItemArsip,
		NoFolder:               arsip.NoFolder,
		KodeKlasifikasi:        arsip.KodeKlasifikasi,
		UraianInformasi:        arsip.UraianInformasiArsip,
		UraianInformasiArsip:   arsip.UraianInformasiArsip,
		JenisArsip:             arsip.JenisArsip,
		KurunWaktu:             arsip.KurunWaktu,
		TingkatPerkembangan:    arsip.TingkatPerkembangan,
		Jumlah:                 arsip.Jumlah,
		Keterangan:             arsip.Keterangan,
		NomorDefinitif:         arsip.NomorDefinitif,
		LokasiSimpanID:         arsip.LokasiSimpanID,
		LokasiSimpan:           lokasiNama,
		LokasiSimpanNama:       lokasiNama,
		NoBoks:                 arsip.NoBoks,
		NoRak:                  arsip.NoRak,
		JangkaSimpanInaktif:    arsip.JangkaSimpanInaktif,
		NasibAkhir:             nasib,
		JangkaSimpanNasibAkhir: arsip.JangkaSimpanNasibAkhir,
		SKKAD:                  arsip.SKKAD,
		KategoriArsip:          arsip.KategoriArsip,
		UnitKerjaID:            arsip.UnitKerjaID,
		UnitKerjaNama:          unitNama,
		UnorID:                 arsip.UnorID,
		PenciptaArsip:          unitNama,
		PetugasInputID:         arsip.PetugasInputID,
		PetugasInput:           petugasNama,
		Status:                 arsip.Status,
		CreatedAt:              arsip.CreatedAt,
		UpdatedAt:              arsip.UpdatedAt,
	}
}

func strPtrVal(s string) *string {
	return &s
}
