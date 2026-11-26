package seeders

import (
	"log"

	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
)

// UNOR ID untuk Biro Data dan Informasi
const unorBiroDataInformasi = "a1a2a3a4-b5c6-7890-abcd-ef1234567890"

// SeedArsipDefault membuat data arsip awal yang menggunakan enum kanonik dan kode klasifikasi valid.
func SeedArsipDefault(db *gorm.DB) error {
	type arsipSeed struct {
		NoItemArsip            string
		NoFolder               string
		KodeKlasifikasi        string
		UraianInformasiArsip   string
		JenisArsip             string
		KurunWaktu             string
		TingkatPerkembangan    string
		Jumlah                 string
		Keterangan             string
		LokasiSimpan           string
		NoBoks                 string
		NoRak                  string
		JangkaSimpanInaktif    string
		JangkaSimpanNasibAkhir *string
		NasibAkhir             *string
		SKKAD                  string
		KategoriArsip          string
		PenciptaArsip          string
	}

	defaultArsip := []arsipSeed{
		{
			NoItemArsip:            "1",
			NoFolder:               "F-01",
			KodeKlasifikasi:        "DT.04.01",
			UraianInformasiArsip:   "Balasan Surat Permohonan Informasi Anak Penyandang Disabilitas pada SIMFONI PPPA Nomor B-1436/Setmen.Birodatin/DT.04.01/9/2021",
			JenisArsip:             string(model.JenisArsipKonvensional),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "B-001",
			NoRak:                  "R-01",
			JangkaSimpanInaktif:    "4 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "2",
			NoFolder:               "F-01",
			KodeKlasifikasi:        "DT.01.01",
			UraianInformasiArsip:   "Permohonan Informasi Anak Penyandang Disabilitas Pada Sistem Aplikasi SIMFONI PPPA dari Dirjen Rehabsos Kemensos RI Nomor: 620/4tD1.01/08/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganCopy),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "B-001",
			NoRak:                  "R-01",
			JangkaSimpanInaktif:    "4 Tahun",
			JangkaSimpanNasibAkhir: strPtr("PERMANEN"),
			NasibAkhir:             strPtr(string(model.NasibAkhirPermanen)),
			SKKAD:                  string(model.SKKADTerbatas),
			KategoriArsip:          string(model.KategoriVitalVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "3",
			NoFolder:               "F-02",
			KodeKlasifikasi:        "DT.08.01",
			UraianInformasiArsip:   "Data Anugrah Ekapraya (APE) Nota Dinas Nomor: ND-1214/Setmen.Birodatin/DT.08.01/12/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "B-001",
			NoRak:                  "R-01",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("PERMANEN"),
			NasibAkhir:             strPtr(string(model.NasibAkhirPermanen)),
			SKKAD:                  string(model.SKKADRahasia),
			KategoriArsip:          string(model.KategoriTerjagaVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
	}

	for _, a := range defaultArsip {
		var existing model.Arsip
		result := db.Where("no_item_arsip = ? AND unor_id = ?", a.NoItemArsip, unorBiroDataInformasi).First(&existing)
		if result.Error == nil {
			log.Printf("Arsip no_item=%s untuk Biro Data dan Informasi sudah ada, skip", a.NoItemArsip)
			continue
		}

		noFolder := a.NoFolder
		arsip := model.Arsip{
			NoItemArsip:            a.NoItemArsip,
			NoFolder:               &noFolder,
			KodeKlasifikasi:        a.KodeKlasifikasi,
			UraianInformasiArsip:   a.UraianInformasiArsip,
			JenisArsip:             a.JenisArsip,
			KurunWaktu:             a.KurunWaktu,
			TingkatPerkembangan:    a.TingkatPerkembangan,
			Jumlah:                 a.Jumlah,
			Keterangan:             a.Keterangan,
			NomorDefinitif:         nil,
			LokasiSimpan:           a.LokasiSimpan,
			NoBoks:                 a.NoBoks,
			NoRak:                  a.NoRak,
			JangkaSimpanInaktif:    a.JangkaSimpanInaktif,
			JangkaSimpanNasibAkhir: a.JangkaSimpanNasibAkhir,
			NasibAkhir:             a.NasibAkhir,
			SKKAD:                  a.SKKAD,
			KategoriArsip:          a.KategoriArsip,
			UnorID:                 unorBiroDataInformasi,
			PenciptaArsip:          a.PenciptaArsip,
			PetugasInputID:         1,
			Status:                 string(model.StatusArsipFinal),
		}

		if err := db.Create(&arsip).Error; err != nil {
			log.Printf("Gagal insert arsip no_item=%s: %v", a.NoItemArsip, err)
			return err
		}
		log.Printf("✅ Arsip no_item=%s (%s) berhasil ditambahkan", a.NoItemArsip, a.PenciptaArsip)
	}

	return nil
}

// strPtr helper untuk membuat *string dari string literal.
func strPtr(s string) *string {
	return &s
}
