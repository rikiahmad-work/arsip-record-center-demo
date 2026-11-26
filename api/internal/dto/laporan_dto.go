package dto

// ================================================================
// 1. Dashboard DTO (Task 3.1)
// ================================================================

type DashboardStatsResponse struct {
	TotalArsip        int64 `json:"total_arsip"`
	ArsipInaktif      int64 `json:"arsip_inaktif"`
	ArsipVital        int64 `json:"arsip_vital"`
	ArsipTerjaga      int64 `json:"arsip_terjaga"`
	JatuhTempoRetensi int64 `json:"jatuh_tempo_retensi"`
	ArsipDipinjam     int64 `json:"arsip_dipinjam"`
	UsulPenyusutan    int64 `json:"usul_penyusutan"`
}

type DashboardChartPoint struct {
	Bulan  string `json:"bulan"`
	Jumlah int64  `json:"jumlah"`
}

// ================================================================
// 2. Laporan & Rekap DAI 18 Kolom DTO (Task 3.2)
// ================================================================

type LaporanFilterRequest struct {
	Tahun           string `form:"tahun"`
	Bulan           string `form:"bulan"`
	KodeKlasifikasi string `form:"kode"`
	KodeLegacy      string `form:"kode_klasifikasi"`
	UnitKerja       string `form:"unit_kerja"`
	UnitKerjaID     uint   `form:"unit_kerja_id"`
	Kategori        string `form:"kategori"`
	SKKAD           string `form:"skkad"`
}

func (r LaporanFilterRequest) ResolvedKode() string {
	if r.KodeKlasifikasi != "" {
		return r.KodeKlasifikasi
	}
	return r.KodeLegacy
}

type LaporanRekapItem struct {
	Label  string `json:"label"`
	Jumlah int64  `json:"jumlah"`
}

type LaporanRekapResponse struct {
	TotalArsip      int64              `json:"total_arsip"`
	TotalPeminjaman int64              `json:"total_peminjaman"`
	TotalPenyusutan int64              `json:"total_penyusutan"`
	PerKategori     []LaporanRekapItem `json:"per_kategori"`
	PerSKKAD        []LaporanRekapItem `json:"per_skkad"`
	PerKode         []LaporanRekapItem `json:"per_kode"`
	PerUnitKerja    []LaporanRekapItem `json:"per_unit_kerja"`
}

// DAIRowResponse merepresentasikan tepat 18 kolom Format Daftar Arsip Inaktif (DAI).
type DAIRowResponse struct {
	No                  int    `json:"no"`                    // 1. NO.
	NoItemArsip         string `json:"no_item_arsip"`         // 2. NO. ITEM ARSIP
	KodeKlasifikasi     string `json:"kode_klasifikasi"`      // 3. KODE KLASIFIKASI
	UraianInformasi     string `json:"uraian_informasi"`      // 4. URAIAN INFORMASI ARSIP
	JenisArsip          string `json:"jenis_arsip"`           // 5. JENIS ARSIP
	KurunWaktu          string `json:"kurun_waktu"`           // 6. KURUN WAKTU
	TingkatPerkembangan string `json:"tingkat_perkembangan"`  // 7. TINGKAT PERKEMBANGAN
	Jumlah              string `json:"jumlah"`                // 8. JUMLAH
	Keterangan          string `json:"keterangan"`            // 9. KETERANGAN
	NomorDefinitif      string `json:"nomor_definitif"`       // 10. NOMOR DEFINITIF
	LokasiSimpan        string `json:"lokasi_simpan"`         // 11. LOKASI SIMPAN
	NoBoks              string `json:"no_boks"`               // 12. NO. BOKS
	NoRak               string `json:"no_rak"`                // 13. NO. RAK
	JangkaSimpanInaktif string `json:"jangka_simpan_inaktif"` // 14. JANGKA SIMPAN INAKTIF
	NasibAkhir          string `json:"nasib_akhir"`           // 15. NASIB AKHIR
	SKKAD               string `json:"skkad"`                 // 16. SKKAD
	KategoriArsip       string `json:"kategori_arsip"`        // 17. KATEGORI ARSIP
	PenciptaArsip       string `json:"pencipta_arsip"`        // 18. PENCIPTA ARSIP
}

// ================================================================
// 3. Cetak Label PDF DTO (Task 3.3)
// ================================================================

type LabelQueryRequest struct {
	UnitKerja   string `form:"unit_kerja"`
	Tahun       string `form:"tahun"`
	NoBoks      string `form:"no_boks"`
	NoRak       string `form:"no_rak"`
	NoFolder    string `form:"no_folder"`
	Kode        string `form:"kode"`
	RentangItem string `form:"rentang_item"`
	Uraian      string `form:"uraian"`
}
