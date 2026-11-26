package dto

// StoreArsipRequest digunakan untuk request tambah arsip baru (20 kolom DAI + No. Folder).
// Mendukung baik nama field kanonik (uraian_informasi, nasib_akhir, unit_kerja_id, lokasi_simpan_id)
// maupun nama field legacy (uraian_informasi_arsip, jangka_simpan_nasib_akhir, unor_id, lokasi_simpan).
type StoreArsipRequest struct {
	NoItemArsip            string  `json:"no_item_arsip"              binding:"omitempty,max=100"`
	NoFolder               *string `json:"no_folder"                  binding:"omitempty,max=50"`
	KodeKlasifikasi        string  `json:"kode_klasifikasi"           binding:"required,max=50"`
	UraianInformasi        string  `json:"uraian_informasi"`
	UraianInformasiArsip   string  `json:"uraian_informasi_arsip"`
	JenisArsip             string  `json:"jenis_arsip"                binding:"required"`
	KurunWaktu             string  `json:"kurun_waktu"                binding:"required,max=50"`
	TingkatPerkembangan    string  `json:"tingkat_perkembangan"       binding:"required"`
	Jumlah                 string  `json:"jumlah"                     binding:"required,max=100"`
	Keterangan             string  `json:"keterangan"                 binding:"omitempty,max=255"`
	NomorDefinitif         *string `json:"nomor_definitif"            binding:"omitempty,max=100"`
	LokasiSimpanID         *uint   `json:"lokasi_simpan_id"`
	LokasiSimpan           string  `json:"lokasi_simpan"              binding:"omitempty,max=255"`
	NoBoks                 string  `json:"no_boks"                    binding:"omitempty,max=50"`
	NoRak                  string  `json:"no_rak"                     binding:"omitempty,max=50"`
	JangkaSimpanInaktif    string  `json:"jangka_simpan_inaktif"      binding:"omitempty,max=100"`
	JangkaSimpanNasibAkhir *string `json:"jangka_simpan_nasib_akhir"`
	NasibAkhir             *string `json:"nasib_akhir"`
	SKKAD                  string  `json:"skkad"                      binding:"omitempty,max=50"`
	KategoriArsip          string  `json:"kategori_arsip"             binding:"omitempty,max=100"`
	UnitKerjaID            *uint   `json:"unit_kerja_id"`
	UnorID                 string  `json:"unor_id"                    binding:"omitempty,max=36"`
	PenciptaArsip          string  `json:"pencipta_arsip"             binding:"omitempty,max=255"`
	Status                 string  `json:"status"                     binding:"omitempty,oneof=draft final"`
}

func (r StoreArsipRequest) ResolvedUraian() string {
	if r.UraianInformasi != "" {
		return r.UraianInformasi
	}
	return r.UraianInformasiArsip
}

// UpdateArsipRequest digunakan untuk request update arsip.
type UpdateArsipRequest struct {
	NoItemArsip            string  `json:"no_item_arsip"              binding:"omitempty,max=100"`
	NoFolder               *string `json:"no_folder"                  binding:"omitempty,max=50"`
	KodeKlasifikasi        string  `json:"kode_klasifikasi"           binding:"required,max=50"`
	UraianInformasi        string  `json:"uraian_informasi"`
	UraianInformasiArsip   string  `json:"uraian_informasi_arsip"`
	JenisArsip             string  `json:"jenis_arsip"                binding:"required"`
	KurunWaktu             string  `json:"kurun_waktu"                binding:"required,max=50"`
	TingkatPerkembangan    string  `json:"tingkat_perkembangan"       binding:"required"`
	Jumlah                 string  `json:"jumlah"                     binding:"required,max=100"`
	Keterangan             string  `json:"keterangan"                 binding:"omitempty,max=255"`
	NomorDefinitif         *string `json:"nomor_definitif"            binding:"omitempty,max=100"`
	LokasiSimpanID         *uint   `json:"lokasi_simpan_id"`
	LokasiSimpan           string  `json:"lokasi_simpan"              binding:"omitempty,max=255"`
	NoBoks                 string  `json:"no_boks"                    binding:"omitempty,max=50"`
	NoRak                  string  `json:"no_rak"                     binding:"omitempty,max=50"`
	JangkaSimpanInaktif    string  `json:"jangka_simpan_inaktif"      binding:"omitempty,max=100"`
	JangkaSimpanNasibAkhir *string `json:"jangka_simpan_nasib_akhir"`
	NasibAkhir             *string `json:"nasib_akhir"`
	SKKAD                  string  `json:"skkad"                      binding:"omitempty,max=50"`
	KategoriArsip          string  `json:"kategori_arsip"             binding:"omitempty,max=100"`
	UnitKerjaID            *uint   `json:"unit_kerja_id"`
	UnorID                 string  `json:"unor_id"                    binding:"omitempty,max=36"`
	PenciptaArsip          string  `json:"pencipta_arsip"             binding:"omitempty,max=255"`
	Status                 string  `json:"status"                     binding:"omitempty,oneof=draft final"`
}

func (r UpdateArsipRequest) ResolvedUraian() string {
	if r.UraianInformasi != "" {
		return r.UraianInformasi
	}
	return r.UraianInformasiArsip
}

// ArsipListRequest digunakan untuk parameter query list arsip.
// Mendukung parameter kanonik (`q`, `kode`, `kategori`, `skkad`, `tahun`, `unit_kerja`, `lokasi`)
// sekaligus parameter legacy (`search`, `kode_klasifikasi`, `kategori_arsip`, `lokasi_simpan`).
type ArsipListRequest struct {
	Page            int    `form:"page"              binding:"omitempty,min=1"`
	PerPage         int    `form:"per_page"          binding:"omitempty,min=1,max=500"`
	Q               string `form:"q"`
	Search          string `form:"search"`
	Kode            string `form:"kode"`
	KodeKlasifikasi string `form:"kode_klasifikasi"`
	JenisArsip      string `form:"jenis_arsip"`
	Kategori        string `form:"kategori"`
	KategoriArsip   string `form:"kategori_arsip"`
	SKKAD           string `form:"skkad"`
	Tahun           string `form:"tahun"`
	UnitKerja       string `form:"unit_kerja"`
	UnitKerjaID     uint   `form:"unit_kerja_id"`
	LokasiSimpan    string `form:"lokasi_simpan"`
	LokasiSimpanID  uint   `form:"lokasi_simpan_id"`
	Status          string `form:"status"            binding:"omitempty,oneof=draft final"`
	SortBy          string `form:"sort_by"`
	SortOrder       string `form:"sort_order"        binding:"omitempty,oneof=asc desc"`
}

func (r ArsipListRequest) ResolvedSearch() string {
	if r.Q != "" {
		return r.Q
	}
	return r.Search
}

func (r ArsipListRequest) ResolvedKode() string {
	if r.Kode != "" {
		return r.Kode
	}
	return r.KodeKlasifikasi
}

func (r ArsipListRequest) ResolvedKategori() string {
	if r.Kategori != "" {
		return r.Kategori
	}
	return r.KategoriArsip
}
