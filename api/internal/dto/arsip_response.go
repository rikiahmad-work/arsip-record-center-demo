package dto

import "time"

// ArsipResponse merepresentasikan data arsip yang dikembalikan ke client (20 field DAI + No. Folder).
type ArsipResponse struct {
	ID                     uint      `json:"id"`
	NoItemArsip            string    `json:"no_item_arsip"`
	NoFolder               *string   `json:"no_folder"`
	KodeKlasifikasi        string    `json:"kode_klasifikasi"`
	UraianInformasi        string    `json:"uraian_informasi"`
	UraianInformasiArsip   string    `json:"uraian_informasi_arsip"`
	JenisArsip             string    `json:"jenis_arsip"`
	KurunWaktu             string    `json:"kurun_waktu"`
	TingkatPerkembangan    string    `json:"tingkat_perkembangan"`
	Jumlah                 string    `json:"jumlah"`
	Keterangan             string    `json:"keterangan"`
	NomorDefinitif         *string   `json:"nomor_definitif"`
	LokasiSimpanID         *uint     `json:"lokasi_simpan_id"`
	LokasiSimpan           string    `json:"lokasi_simpan"`
	LokasiSimpanNama       string    `json:"lokasi_simpan_nama"`
	NoBoks                 string    `json:"no_boks"`
	NoRak                  string    `json:"no_rak"`
	JangkaSimpanInaktif    string    `json:"jangka_simpan_inaktif"`
	NasibAkhir             *string   `json:"nasib_akhir"`
	JangkaSimpanNasibAkhir *string   `json:"jangka_simpan_nasib_akhir"`
	SKKAD                  string    `json:"skkad"`
	KategoriArsip          string    `json:"kategori_arsip"`
	UnitKerjaID            *uint     `json:"unit_kerja_id"`
	UnitKerjaNama          string    `json:"unit_kerja_nama"`
	UnorID                 string    `json:"unor_id"`
	PenciptaArsip          string    `json:"pencipta_arsip"`
	PetugasInputID         uint      `json:"petugas_input_id"`
	PetugasInput           string    `json:"petugas_input,omitempty"`
	Status                 string    `json:"status"`
	TanggalJatuhTempo      string    `json:"tanggal_jatuh_tempo"`
	IsJatuhTempo           bool      `json:"is_jatuh_tempo"`
	RekomendasiRetensi     string    `json:"rekomendasi_retensi"`
	SedangDipinjam         bool      `json:"sedang_dipinjam"`
	StatusPinjam           string    `json:"status_pinjam,omitempty"`
	DalamPenyusutan        bool      `json:"dalam_penyusutan"`
	StatusPenyusutan       string    `json:"status_penyusutan,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// ArsipStatsResponse digunakan untuk statistik arsip per kategori & SKKAD.
type ArsipStatsResponse struct {
	Total           int64            `json:"total"`
	ArsipInaktif    int64            `json:"arsip_inaktif"`
	ArsipVital      int64            `json:"arsip_vital"`
	ArsipTerjaga    int64            `json:"arsip_terjaga"`
	ArsipUsulMusnah int64            `json:"arsip_usul_musnah"`
	ArsipUsulSerah  int64            `json:"arsip_usul_serah"`
	ArsipAktif      int64            `json:"arsip_aktif"` // alias kompatibilitas
	ArsipDraft      int64            `json:"arsip_draft"`
	ArsipFinal      int64            `json:"arsip_final"`
	BySKKAD         map[string]int64 `json:"by_skkad"`
}

// ArsipImportResultResponse merepresentasikan hasil import CSV/Excel arsip.
type ArsipImportResultResponse struct {
	TotalRows    int      `json:"total_rows"`
	ImportedRows int      `json:"imported_rows"`
	FailedRows   int      `json:"failed_rows"`
	Errors       []string `json:"errors,omitempty"`
}
