package service

import (
	"context"
	"strings"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// LaporanService menyediakan business logic untuk Dashboard, Laporan Rekap, Rekap Format DAI (18 kolom), dan Cetak Label PDF.
type LaporanService interface {
	GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
	GetDashboardChart(ctx context.Context, tahun int) ([]dto.DashboardChartPoint, error)
	GetLaporanRekap(ctx context.Context, req dto.LaporanFilterRequest) (*dto.LaporanRekapResponse, error)
	GetLaporanDAI(ctx context.Context, req dto.LaporanFilterRequest) ([]dto.DAIRowResponse, error)
	ExportLaporanDAICSV(ctx context.Context, req dto.LaporanFilterRequest) ([]byte, error)
	ExportLaporanDAIExcel(ctx context.Context, req dto.LaporanFilterRequest) ([]byte, error)
	GenerateLabel(ctx context.Context, jenis string, req dto.LabelQueryRequest) ([]byte, error)
}

type laporanServiceImpl struct {
	db *gorm.DB
}

func NewLaporanService(db *gorm.DB) LaporanService {
	return &laporanServiceImpl{db: db}
}

var namaBulanSingkat = []string{
	"Jan", "Feb", "Mar", "Apr", "Mei", "Jun",
	"Jul", "Agu", "Sep", "Okt", "Nov", "Des",
}

func (s *laporanServiceImpl) GetDashboardStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	var total, inaktif, vital, terjaga, usul, dipinjam int64

	s.db.WithContext(ctx).Model(&model.Arsip{}).Count(&total)
	s.db.WithContext(ctx).Model(&model.Arsip{}).Where("UPPER(kategori_arsip) IN ('INAKTIF', 'AKTIF')").Count(&inaktif)
	s.db.WithContext(ctx).Model(&model.Arsip{}).Where("UPPER(kategori_arsip) = 'VITAL'").Count(&vital)
	s.db.WithContext(ctx).Model(&model.Arsip{}).Where("UPPER(kategori_arsip) = 'TERJAGA'").Count(&terjaga)
	s.db.WithContext(ctx).Model(&model.Arsip{}).Where("UPPER(kategori_arsip) IN ('USUL MUSNAH', 'USUL SERAH')").Count(&usul)

	s.db.WithContext(ctx).Model(&model.Peminjaman{}).Where("status IN ('Dipinjam', 'Terlambat')").Count(&dipinjam)

	var penyusutanPending int64
	s.db.WithContext(ctx).Model(&model.Penyusutan{}).Where("status = 'Diajukan'").Count(&penyusutanPending)

	// Hitung arsip yang jatuh tempo retensi per 31 Desember sesuai standar ANRI
	var allArsip []model.Arsip
	s.db.WithContext(ctx).Model(&model.Arsip{}).
		Select("id, kurun_waktu, jangka_simpan_inaktif, nasib_akhir, jangka_simpan_nasib_akhir, kategori_arsip").
		Find(&allArsip)

	now := time.Now()
	var jatuhTempoANRI int64
	for _, a := range allArsip {
		nasib := "Musnah"
		if a.NasibAkhir != nil && *a.NasibAkhir != "" {
			nasib = *a.NasibAkhir
		} else if a.JangkaSimpanNasibAkhir != nil && *a.JangkaSimpanNasibAkhir != "" {
			nasib = *a.JangkaSimpanNasibAkhir
		}
		_, _, _, isDue, _ := helper.HitungRetensiANRI(a.KurunWaktu, a.JangkaSimpanInaktif, nasib, a.KategoriArsip, now)
		if isDue {
			jatuhTempoANRI++
		}
	}

	jatuhTempo := jatuhTempoANRI
	if usul+penyusutanPending > jatuhTempo {
		jatuhTempo = usul + penyusutanPending
	}

	return &dto.DashboardStatsResponse{
		TotalArsip:        total,
		ArsipInaktif:      inaktif,
		ArsipVital:        vital,
		ArsipTerjaga:      terjaga,
		JatuhTempoRetensi: jatuhTempo,
		ArsipDipinjam:     dipinjam,
		UsulPenyusutan:    usul,
	}, nil
}

func (s *laporanServiceImpl) GetDashboardChart(ctx context.Context, tahun int) ([]dto.DashboardChartPoint, error) {
	if tahun <= 0 {
		tahun = time.Now().Year()
	}

	type monthRow struct {
		Bulan  int
		Jumlah int64
	}
	var rows []monthRow
	_ = s.db.WithContext(ctx).Model(&model.Arsip{}).
		Select("MONTH(created_at) as bulan, COUNT(*) as jumlah").
		Where("YEAR(created_at) = ?", tahun).
		Group("MONTH(created_at)").
		Scan(&rows).Error

	byMonth := make(map[int]int64, 12)
	for _, r := range rows {
		byMonth[r.Bulan] = r.Jumlah
	}

	res := make([]dto.DashboardChartPoint, 0, 12)
	for m := 1; m <= 12; m++ {
		res = append(res, dto.DashboardChartPoint{
			Bulan:  namaBulanSingkat[m-1],
			Jumlah: byMonth[m],
		})
	}
	return res, nil
}

func (s *laporanServiceImpl) buildFilteredArsipQuery(ctx context.Context, req dto.LaporanFilterRequest) *gorm.DB {
	q := s.db.WithContext(ctx).Model(&model.Arsip{})
	if kode := req.ResolvedKode(); kode != "" {
		q = q.Where("kode_klasifikasi = ?", kode)
	}
	if req.Tahun != "" {
		q = q.Where("kurun_waktu LIKE ?", "%"+req.Tahun+"%")
	}
	if req.Kategori != "" {
		q = q.Where("UPPER(kategori_arsip) = ?", strings.ToUpper(req.Kategori))
	}
	if req.SKKAD != "" {
		q = q.Where("UPPER(skkad) = ?", strings.ToUpper(req.SKKAD))
	}
	if req.UnitKerjaID > 0 {
		q = q.Where("unit_kerja_id = ?", req.UnitKerjaID)
	} else if req.UnitKerja != "" {
		q = q.Where("pencipta_arsip LIKE ?", "%"+req.UnitKerja+"%")
	}
	return q
}

func (s *laporanServiceImpl) GetLaporanRekap(ctx context.Context, req dto.LaporanFilterRequest) (*dto.LaporanRekapResponse, error) {
	var totalArsip, totalPeminjaman, totalPenyusutan int64

	baseQ := s.buildFilteredArsipQuery(ctx, req)
	baseQ.Count(&totalArsip)
	s.db.WithContext(ctx).Model(&model.Peminjaman{}).Count(&totalPeminjaman)
	s.db.WithContext(ctx).Model(&model.Penyusutan{}).Count(&totalPenyusutan)

	var perKategori, perSKKAD, perKode, perUnit []dto.LaporanRekapItem

	s.buildFilteredArsipQuery(ctx, req).
		Select("kategori_arsip as label, COUNT(*) as jumlah").
		Group("kategori_arsip").
		Scan(&perKategori)

	s.buildFilteredArsipQuery(ctx, req).
		Select("skkad as label, COUNT(*) as jumlah").
		Group("skkad").
		Scan(&perSKKAD)

	s.buildFilteredArsipQuery(ctx, req).
		Select("kode_klasifikasi as label, COUNT(*) as jumlah").
		Group("kode_klasifikasi").
		Order("jumlah DESC").
		Limit(20).
		Scan(&perKode)

	s.buildFilteredArsipQuery(ctx, req).
		Select("pencipta_arsip as label, COUNT(*) as jumlah").
		Group("pencipta_arsip").
		Order("jumlah DESC").
		Scan(&perUnit)

	return &dto.LaporanRekapResponse{
		TotalArsip:      totalArsip,
		TotalPeminjaman: totalPeminjaman,
		TotalPenyusutan: totalPenyusutan,
		PerKategori:     perKategori,
		PerSKKAD:        perSKKAD,
		PerKode:         perKode,
		PerUnitKerja:    perUnit,
	}, nil
}

// GetLaporanDAI mengembalikan daftar arsip dalam format tepat 18 kolom Daftar Arsip Inaktif (DAI).
func (s *laporanServiceImpl) GetLaporanDAI(ctx context.Context, req dto.LaporanFilterRequest) ([]dto.DAIRowResponse, error) {
	var list []model.Arsip
	if err := s.buildFilteredArsipQuery(ctx, req).
		Order("id ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}

	return FormatDAIRows(list), nil
}

// FormatDAIRows memetakan slice model.Arsip menjadi 18 kolom kanonik Format Daftar Arsip Inaktif.
func FormatDAIRows(list []model.Arsip) []dto.DAIRowResponse {
	rows := make([]dto.DAIRowResponse, 0, len(list))
	for i, a := range list {
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

		rows = append(rows, dto.DAIRowResponse{
			No:                  i + 1,
			NoItemArsip:         a.NoItemArsip,
			KodeKlasifikasi:     a.KodeKlasifikasi,
			UraianInformasi:     a.UraianInformasiArsip,
			JenisArsip:          a.JenisArsip,
			KurunWaktu:          helper.NormalizeKurunWaktuExcel(a.KurunWaktu),
			TingkatPerkembangan: a.TingkatPerkembangan,
			Jumlah:              a.Jumlah,
			Keterangan:          a.Keterangan,
			NomorDefinitif:      noDef,
			LokasiSimpan:        a.LokasiSimpan,
			NoBoks:              a.NoBoks,
			NoRak:               a.NoRak,
			JangkaSimpanInaktif: a.JangkaSimpanInaktif,
			NasibAkhir:          nasib,
			SKKAD:               a.SKKAD,
			KategoriArsip:       a.KategoriArsip,
			PenciptaArsip:       a.PenciptaArsip,
		})
	}
	return rows
}

func (s *laporanServiceImpl) ExportLaporanDAICSV(ctx context.Context, req dto.LaporanFilterRequest) ([]byte, error) {
	rows, err := s.GetLaporanDAI(ctx, req)
	if err != nil {
		return nil, err
	}
	return helper.GenerateDAICSV(toHelperDAIRows(rows))
}

func (s *laporanServiceImpl) ExportLaporanDAIExcel(ctx context.Context, req dto.LaporanFilterRequest) ([]byte, error) {
	rows, err := s.GetLaporanDAI(ctx, req)
	if err != nil {
		return nil, err
	}
	return helper.GenerateDAIExcel(toHelperDAIRows(rows), req.Tahun)
}

func toHelperDAIRows(rows []dto.DAIRowResponse) []helper.DAIRowData {
	res := make([]helper.DAIRowData, 0, len(rows))
	for _, r := range rows {
		res = append(res, helper.DAIRowData{
			No:                  r.No,
			NoItemArsip:         r.NoItemArsip,
			KodeKlasifikasi:     r.KodeKlasifikasi,
			UraianInformasi:     r.UraianInformasi,
			JenisArsip:          r.JenisArsip,
			KurunWaktu:          helper.NormalizeKurunWaktuExcel(r.KurunWaktu),
			TingkatPerkembangan: r.TingkatPerkembangan,
			Jumlah:              r.Jumlah,
			Keterangan:          r.Keterangan,
			NomorDefinitif:      r.NomorDefinitif,
			LokasiSimpan:        r.LokasiSimpan,
			NoBoks:              r.NoBoks,
			NoRak:               r.NoRak,
			JangkaSimpanInaktif: r.JangkaSimpanInaktif,
			NasibAkhir:          r.NasibAkhir,
			SKKAD:               r.SKKAD,
			KategoriArsip:       r.KategoriArsip,
			PenciptaArsip:       r.PenciptaArsip,
		})
	}
	return res
}

func (s *laporanServiceImpl) GenerateLabel(ctx context.Context, jenis string, req dto.LabelQueryRequest) ([]byte, error) {
	title := "LABEL BOKS ARSIP"
	switch strings.ToLower(jenis) {
	case "folder":
		title = "LABEL FOLDER ARSIP"
	case "inaktif":
		title = "LABEL ARSIP INAKTIF"
	}

	pdfBytes := helper.GenerateLabelPDF(helper.LabelEntry{
		JenisLabel:  title,
		UnitKerja:   req.UnitKerja,
		Kode:        req.Kode,
		Tahun:       req.Tahun,
		NoRak:       req.NoRak,
		NoBoks:      req.NoBoks,
		NoFolder:    req.NoFolder,
		RentangItem: req.RentangItem,
		Uraian:      req.Uraian,
	})
	return pdfBytes, nil
}
