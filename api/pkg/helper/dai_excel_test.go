package helper

import (
	"testing"
	"time"
)

func TestNormalizeKurunWaktuExcel(t *testing.T) {
	if got := NormalizeKurunWaktuExcel("44310"); got != "24 April 2021" {
		t.Fatalf("expected '24 April 2021', got %q", got)
	}
	if got := NormalizeKurunWaktuExcel("44292"); got != "6 April 2021" {
		t.Fatalf("expected '6 April 2021', got %q", got)
	}
	if got := NormalizeKurunWaktuExcel("2021 - 2022"); got != "2021 - 2022" {
		t.Fatalf("expected '2021 - 2022', got %q", got)
	}
}

func TestHitungRetensiANRI_Cutoff31Desember(t *testing.T) {
	now2026 := time.Date(2026, 7, 10, 10, 0, 0, 0, time.UTC)

	// 1. Arsip tahun 2021 - 2022 dengan retensi inaktif 3 Tahun -> jatuh tempo 31 Desember 2025 (Sudah Jatuh Tempo di 2026)
	thAkhir, dur, tgl, isDue, rek := HitungRetensiANRI("2021 - 2022", "3 Tahun", "Musnah", "Inaktif", now2026)
	if thAkhir != 2022 || dur != 3 || tgl != "31 Desember 2025" || !isDue || rek != "Siap Usul Musnah" {
		t.Fatalf("unexpected retensi ANRI result: thAkhir=%d dur=%d tgl=%s isDue=%v rek=%s", thAkhir, dur, tgl, isDue, rek)
	}

	// 2. Arsip tahun 2021 - 2022 dengan Nasib Akhir Permanen -> Siap Usul Serah
	_, _, _, isDuePerm, rekPerm := HitungRetensiANRI("2021 - 2022", "3 Tahun", "Permanen", "Vital", now2026)
	if !isDuePerm || rekPerm != "Siap Usul Serah" {
		t.Fatalf("expected Siap Usul Serah, got isDue=%v rek=%s", isDuePerm, rekPerm)
	}

	// 3. Arsip tahun 2025 dengan retensi 3 Tahun -> jatuh tempo 31 Desember 2028 (Belum Jatuh Tempo di 2026)
	_, _, tglFuture, isDueFuture, _ := HitungRetensiANRI("2025", "3 Tahun", "Musnah", "Inaktif", now2026)
	if isDueFuture || tglFuture != "31 Desember 2028" {
		t.Fatalf("expected not due until 31 Desember 2028, got isDue=%v tgl=%s", isDueFuture, tglFuture)
	}
}

func TestGenerateAndParseDAIExcelRoundTrip(t *testing.T) {
	sample := []DAIRowData{
		{
			No:                  1,
			NoItemArsip:         "1",
			KodeKlasifikasi:     "DT.10",
			UraianInformasi:     "Balasan Surat Permohonan Informasi Anak Penyandang Disabilitas",
			JenisArsip:          "KONVENSIONAL",
			KurunWaktu:          "2021 - 2022",
			TingkatPerkembangan: "ASLI",
			Jumlah:              "1 Lembar",
			Keterangan:          "KONDISI BAIK",
			NomorDefinitif:      "",
			LokasiSimpan:        "RECORDS CENTER",
			NoBoks:              "",
			NoRak:               "",
			JangkaSimpanInaktif: "3 Tahun",
			NasibAkhir:          "Musnah",
			SKKAD:               "Biasa",
			KategoriArsip:       "INAKTIF",
			PenciptaArsip:       "Biro Data dan Informasi",
		},
	}

	xlsxBytes, err := GenerateDAIExcel(sample, "2021 - 2022")
	if err != nil {
		t.Fatalf("GenerateDAIExcel failed: %v", err)
	}
	if len(xlsxBytes) < 100 {
		t.Fatalf("expected valid xlsx zip bytes, got len=%d", len(xlsxBytes))
	}

	// Pastikan file valid XLSX dengan magic bytes PK\x03\x04
	if len(xlsxBytes) < 4 || xlsxBytes[0] != 'P' || xlsxBytes[1] != 'K' || xlsxBytes[2] != 3 || xlsxBytes[3] != 4 {
		t.Fatalf("file output bukan zip/xlsx yang valid")
	}

	parsedRows, err := ParseDAIImportFile(xlsxBytes)
	if err != nil {
		t.Fatalf("ParseDAIImportFile failed on generated xlsx: %v", err)
	}
	if len(parsedRows) < 5 {
		t.Fatalf("expected at least 5 rows (title, tahun, 2 headers, 1 data), got %d", len(parsedRows))
	}
	dataRow := parsedRows[len(parsedRows)-1]
	if dataRow[2] != "DT.10" || dataRow[17] != "Biro Data dan Informasi" {
		t.Fatalf("unexpected parsed data row: %+v", dataRow)
	}
}
