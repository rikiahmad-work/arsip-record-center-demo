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
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Balasan Surat Permohonan Informasi Anak Penyandang Disabilitas pada SIMFONI PPPA Nomor B-1436/Setmen.Birodatin/DT.10/9/2021",
			JenisArsip:             string(model.JenisArsipKonvensional),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "2",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Permohonan Informasi Anak Penyandang Disabilitas Pada Sistem Aplikasi SIMFONI PPPA dari Dirjen Rehabsos Kemensos RI Nomor: 620/4tD1.01/08/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganCopy),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("PERMANEN"),
			NasibAkhir:             strPtr(string(model.NasibAkhirPermanen)),
			SKKAD:                  string(model.SKKADTerbatas),
			KategoriArsip:          string(model.KategoriVitalVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "3",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Data Anugrah Ekapraya (APE) Nota Dinas Nomor: ND-1214/Setmen.Birodatin/DT.10/12/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADRahasia),
			KategoriArsip:          string(model.KategoriTerjagaVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "4",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Surat Ground Check hasil perhitungan IPA,IPHA, IPKA tingkat Kab/Kota Tahun 2020 Nomor : B-1788/Setmen.Birodatin/DT.10/11/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "1 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "5",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Kerangka Acuan Ground Check hasil perhitungan IPA,IPHA, IPKA tingkat Kab/Kota Tahun 2020",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "4 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "6",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Sosialisasi Hasil Penghitungan IPA, IPHA, dan IPKA Tingkat Provinsi Nomor : B-1671/Setmen.Birodatin/DT.10/10/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "2021 - 2022",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "4 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
		{
			NoItemArsip:            "7",
			NoFolder:               "",
			KodeKlasifikasi:        "DT.10",
			UraianInformasiArsip:   "Prorata Proyeksi Indikator Input IPA, IPHA, dan IPKA tahun 2021-2024 KLASTER 3 nomor B-1262/Setmen.Birodatin/DT.10/7/2021",
			JenisArsip:             string(model.JenisArsipElektronik),
			KurunWaktu:             "28 Juli 2021",
			TingkatPerkembangan:    string(model.TingkatPerkembanganAsli),
			Jumlah:                 "3 Lembar",
			Keterangan:             "KONDISI BAIK",
			LokasiSimpan:           "Record Center",
			NoBoks:                 "",
			NoRak:                  "",
			JangkaSimpanInaktif:    "3 Tahun",
			JangkaSimpanNasibAkhir: strPtr("MUSNAH"),
			NasibAkhir:             strPtr(string(model.NasibAkhirMusnah)),
			SKKAD:                  string(model.SKKADBiasa),
			KategoriArsip:          string(model.KategoriInaktifVal),
			PenciptaArsip:          "Biro Data dan Informasi",
		},
	}

	var ukBiroDatin model.UnitKerja
	unorKode := unorBiroDataInformasi
	var unitKerjaID *uint
	if err := db.Where("kode = ? OR nama = ?", "BD.02", "Biro Data dan Informasi").First(&ukBiroDatin).Error; err == nil {
		unorKode = ukBiroDatin.Kode
		unitKerjaID = &ukBiroDatin.ID
	}

	for _, a := range defaultArsip {
		var existing model.Arsip
		result := db.Where("no_item_arsip = ? AND (unor_id = ? OR unor_id = ? OR pencipta_arsip = ?)", a.NoItemArsip, unorKode, unorBiroDataInformasi, a.PenciptaArsip).First(&existing)
		if result.Error == nil {
			// Sinkronkan kode_klasifikasi ke DT.10 jika baris berasal dari seeder lama
			if existing.KodeKlasifikasi == "DT.04.01" || existing.KodeKlasifikasi == "DT.01.01" || existing.KodeKlasifikasi == "DT.08.01" {
				_ = db.Model(&existing).Updates(map[string]interface{}{
					"kode_klasifikasi":       a.KodeKlasifikasi,
					"uraian_informasi_arsip": a.UraianInformasiArsip,
					"no_boks":                a.NoBoks,
					"no_rak":                 a.NoRak,
					"no_folder":              nil,
					"jangka_simpan_inaktif":  a.JangkaSimpanInaktif,
					"nasib_akhir":            a.NasibAkhir,
					"skkad":                  a.SKKAD,
					"kategori_arsip":         a.KategoriArsip,
				}).Error
				log.Printf("🔄 Arsip no_item=%s diperbarui ke kode %s sesuai DAI", a.NoItemArsip, a.KodeKlasifikasi)
			} else {
				log.Printf("Arsip no_item=%s untuk Biro Data dan Informasi sudah ada, skip", a.NoItemArsip)
			}
			continue
		}

		var noFolderPtr *string
		if a.NoFolder != "" {
			nf := a.NoFolder
			noFolderPtr = &nf
		}
		arsip := model.Arsip{
			NoItemArsip:            a.NoItemArsip,
			NoFolder:               noFolderPtr,
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
			UnitKerjaID:            unitKerjaID,
			UnorID:                 unorKode,
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
