package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// DAIRowData merepresentasikan 18 kolom baku Format Daftar Arsip Inaktif (DAI) KemenPPPA & ANRI.
type DAIRowData struct {
	No                  int
	NoItemArsip         string
	KodeKlasifikasi     string
	UraianInformasi     string
	JenisArsip          string
	KurunWaktu          string
	TingkatPerkembangan string
	Jumlah              string
	Keterangan          string
	NomorDefinitif      string
	LokasiSimpan        string
	NoBoks              string
	NoRak               string
	NoFolder            string
	JangkaSimpanInaktif string
	NasibAkhir          string
	SKKAD               string
	KategoriArsip       string
	PenciptaArsip       string
}

var (
	reFourDigitYear = regexp.MustCompile(`\b(19\d{2}|20\d{2})\b`)
	reFirstDigits   = regexp.MustCompile(`\d+`)
	bulanIndoFull   = []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
)

// NormalizeKurunWaktuExcel mengonversi angka serial tanggal Excel (misal "44310") menjadi format tanggal Indonesia ("24 April 2021").
func NormalizeKurunWaktuExcel(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if len(trimmed) == 5 {
		if serial, err := strconv.Atoi(trimmed); err == nil && serial >= 25000 && serial <= 60000 {
			// Excel base date: 1899-12-30 (memperhitungkan leap year bug 1900 Excel)
			base := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
			dt := base.AddDate(0, 0, serial)
			return fmt.Sprintf("%d %s %d", dt.Day(), bulanIndoFull[dt.Month()-1], dt.Year())
		}
	}
	return trimmed
}

// HitungRetensiANRI menghitung jatuh tempo retensi arsip inaktif sesuai standar ANRI (Perka ANRI No. 37/2016)
// dengan perhitungan batas akhir per 31 Desember.
func HitungRetensiANRI(kurunWaktu, jangkaSimpanInaktif, nasibAkhir, kategoriArsip string, now time.Time) (tahunAkhir int, durasiTahun int, tanggalJatuhTempo string, isJatuhTempo bool, rekomendasi string) {
	normKurun := NormalizeKurunWaktuExcel(kurunWaktu)

	// 1. Cari tahun terakhir pada kurun_waktu (misal "2021 - 2022" -> 2022, "28 Juli 2021" -> 2021)
	matches := reFourDigitYear.FindAllString(normKurun, -1)
	for _, m := range matches {
		if y, err := strconv.Atoi(m); err == nil && y > tahunAkhir {
			tahunAkhir = y
		}
	}
	if tahunAkhir == 0 {
		tahunAkhir = now.Year()
	}

	// 2. Cari angka tahun pada jangka_simpan_inaktif (misal "3 Tahun" -> 3)
	if dStr := reFirstDigits.FindString(jangkaSimpanInaktif); dStr != "" {
		if d, err := strconv.Atoi(dStr); err == nil && d >= 0 {
			durasiTahun = d
		}
	}
	if durasiTahun <= 0 {
		durasiTahun = 3
	}

	// 3. Perhitungan akhir per 31 Desember (Tahun Akhir + Durasi Retensi Inaktif)
	tahunJatuhTempo := tahunAkhir + durasiTahun
	tanggalJatuhTempo = fmt.Sprintf("31 Desember %d", tahunJatuhTempo)

	cutoffDate := time.Date(tahunJatuhTempo, 12, 31, 23, 59, 59, 0, now.Location())
	isJatuhTempo = !now.Before(cutoffDate) || now.Year() > tahunJatuhTempo

	katUpper := strings.ToUpper(strings.TrimSpace(kategoriArsip))
	nasibUpper := strings.ToUpper(strings.TrimSpace(nasibAkhir))

	switch {
	case katUpper == "USUL MUSNAH":
		isJatuhTempo = true
		rekomendasi = "Proses Usul Musnah"
	case katUpper == "USUL SERAH":
		isJatuhTempo = true
		rekomendasi = "Proses Usul Serah"
	case isJatuhTempo:
		if strings.Contains(nasibUpper, "PERMANEN") {
			rekomendasi = "Siap Usul Serah"
		} else {
			rekomendasi = "Siap Usul Musnah"
		}
	default:
		rekomendasi = fmt.Sprintf("Inaktif s/d 31 Des %d", tahunJatuhTempo)
	}

	return tahunAkhir, durasiTahun, tanggalJatuhTempo, isJatuhTempo, rekomendasi
}

// GenerateDAICSV menghasilkan file CSV 18 kolom sesuai Format Daftar Arsip Inaktif (DAI) KemenPPPA.
func GenerateDAICSV(rows []DAIRowData) ([]byte, error) {
	var buf bytes.Buffer
	// Tambahkan UTF-8 BOM agar Microsoft Excel di Windows otomatis membaca karakter UTF-8 dengan benar
	buf.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(&buf)
	headers := []string{
		"NO.",
		"NO. ITEM ARSIP",
		"KODE KLASIFIKASI",
		"URAIAN INFORMASI ARSIP",
		"JENIS ARSIP",
		"KURUN WAKTU",
		"TINGKAT PERKEMBANGAN",
		"JUMLAH",
		"KETERANGAN",
		"NOMOR DEFINITIF",
		"LOKASI SIMPAN",
		"NO.BOKS",
		"NO. RAK",
		"JANGKA SIMPAN INAKTIF",
		"NASIB AKHIR",
		"SKKAD",
		"KATEGORI ARSIP",
		"PENCIPTA ARSIP",
	}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for i, r := range rows {
		no := r.No
		if no <= 0 {
			no = i + 1
		}
		record := []string{
			strconv.Itoa(no),
			r.NoItemArsip,
			r.KodeKlasifikasi,
			r.UraianInformasi,
			strings.ToUpper(r.JenisArsip),
			r.KurunWaktu,
			strings.ToUpper(r.TingkatPerkembangan),
			r.Jumlah,
			r.Keterangan,
			r.NomorDefinitif,
			r.LokasiSimpan,
			r.NoBoks,
			r.NoRak,
			r.JangkaSimpanInaktif,
			r.NasibAkhir,
			r.SKKAD,
			strings.ToUpper(r.KategoriArsip),
			r.PenciptaArsip,
		}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	return buf.Bytes(), writer.Error()
}

// GenerateDAIExcel menghasilkan file Excel (.xlsx) asli dengan struktur 18 kolom dan header 2 tingkat (merged cells)
// persis mengikuti template "DAFTAR ARSIP INAKTIF RECORDS CENTER KEMEN PPPA" (sheet DAI).
func GenerateDAIExcel(rows []DAIRowData, tahunLabel string) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheet := "DAI"
	_ = f.SetSheetName("Sheet1", sheet)

	// Lebar 18 Kolom Sesuai Standar Template DAI KemenPPPA
	colWidths := map[string]float64{
		"A": 6,
		"B": 15,
		"C": 18,
		"D": 50,
		"E": 16,
		"F": 16,
		"G": 20,
		"H": 14,
		"I": 18,
		"J": 16,
		"K": 18,
		"L": 12,
		"M": 12,
		"N": 14,
		"O": 14,
		"P": 14,
		"Q": 16,
		"R": 28,
	}
	for col, w := range colWidths {
		_ = f.SetColWidth(sheet, col, col, w)
	}

	// Styles
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 13, Family: "Calibri"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	subTitleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 10, Family: "Calibri"},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 10, Family: "Calibri"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1A3A52"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "D0D0D0", Style: 1},
			{Type: "right", Color: "D0D0D0", Style: 1},
			{Type: "top", Color: "D0D0D0", Style: 1},
			{Type: "bottom", Color: "D0D0D0", Style: 1},
		},
	})
	dataStyleCenter, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "top", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "E0E0E0", Style: 1},
			{Type: "right", Color: "E0E0E0", Style: 1},
			{Type: "top", Color: "E0E0E0", Style: 1},
			{Type: "bottom", Color: "E0E0E0", Style: 1},
		},
	})
	dataStyleLeft, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Size: 10, Family: "Calibri"},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "top", WrapText: true},
		Border: []excelize.Border{
			{Type: "left", Color: "E0E0E0", Style: 1},
			{Type: "right", Color: "E0E0E0", Style: 1},
			{Type: "top", Color: "E0E0E0", Style: 1},
			{Type: "bottom", Color: "E0E0E0", Style: 1},
		},
	})

	// Row 1: Judul Utama
	_ = f.SetRowHeight(sheet, 1, 24)
	_ = f.SetCellValue(sheet, "A1", "DAFTAR ARSIP INAKTIF RECORDS CENTER KEMEN PPPA")
	_ = f.MergeCell(sheet, "A1", "J1")
	_ = f.SetCellStyle(sheet, "A1", "J1", titleStyle)

	// Row 2: TAHUN
	_ = f.SetRowHeight(sheet, 2, 18)
	tahunText := "TAHUN:"
	if strings.TrimSpace(tahunLabel) != "" {
		tahunText = "TAHUN: " + strings.TrimSpace(tahunLabel)
	}
	_ = f.SetCellValue(sheet, "A2", tahunText)
	_ = f.MergeCell(sheet, "A2", "H2")
	_ = f.SetCellStyle(sheet, "A2", "H2", subTitleStyle)

	// Row 3 & 4: Headers Tingkat 1 & Tingkat 2
	_ = f.SetRowHeight(sheet, 3, 22)
	_ = f.SetRowHeight(sheet, 4, 20)

	headers := map[string]string{
		"A3": "NO.",
		"B3": "NO. ITEM ARSIP",
		"C3": "KODE KLASIFIKASI",
		"D3": "URAIAN INFORMASI ARSIP",
		"E3": "JENIS ARSIP",
		"F3": "KURUN WAKTU",
		"G3": "TINGKAT PERKEMBANGAN",
		"H3": "JUMLAH",
		"I3": "KETERANGAN",
		"J3": "NOMOR DEFINITIF",
		"K3": "LOKASI SIMPAN",
		"L3": "NO.BOKS",
		"M3": "NO. RAK",
		"N3": "JANGKA SIMPAN",
		"N4": "INAKTIF",
		"O4": "NASIB AKHIR",
		"P3": "SKKAD",
		"Q3": "KATEGORI ARSIP",
		"R3": "PENCIPTA ARSIP",
	}
	for cell, title := range headers {
		_ = f.SetCellValue(sheet, cell, title)
	}

	merges := []string{
		"A3:A4", "B3:B4", "C3:C4", "D3:D4", "E3:E4", "F3:F4",
		"G3:G4", "H3:H4", "I3:I4", "J3:J4", "K3:K4", "L3:L4", "M3:M4",
		"N3:O3",
		"P3:P4", "Q3:Q4", "R3:R4",
	}
	for _, m := range merges {
		parts := strings.Split(m, ":")
		if len(parts) == 2 {
			_ = f.MergeCell(sheet, parts[0], parts[1])
		}
	}
	_ = f.SetCellStyle(sheet, "A3", "R4", headerStyle)

	// Row 5+: Baris Data
	for i, r := range rows {
		rowNum := i + 5
		no := r.No
		if no <= 0 {
			no = i + 1
		}
		lokasi := r.LokasiSimpan
		if strings.EqualFold(lokasi, "Record Center") {
			lokasi = "RECORDS CENTER"
		}

		_ = f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), no)
		_ = f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), r.NoItemArsip)
		_ = f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), r.KodeKlasifikasi)
		_ = f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), r.UraianInformasi)
		_ = f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), strings.ToUpper(r.JenisArsip))
		_ = f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), r.KurunWaktu)
		_ = f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), strings.ToUpper(r.TingkatPerkembangan))
		_ = f.SetCellValue(sheet, fmt.Sprintf("H%d", rowNum), r.Jumlah)
		_ = f.SetCellValue(sheet, fmt.Sprintf("I%d", rowNum), r.Keterangan)
		_ = f.SetCellValue(sheet, fmt.Sprintf("J%d", rowNum), r.NomorDefinitif)
		_ = f.SetCellValue(sheet, fmt.Sprintf("K%d", rowNum), lokasi)
		_ = f.SetCellValue(sheet, fmt.Sprintf("L%d", rowNum), r.NoBoks)
		_ = f.SetCellValue(sheet, fmt.Sprintf("M%d", rowNum), r.NoRak)
		_ = f.SetCellValue(sheet, fmt.Sprintf("N%d", rowNum), r.JangkaSimpanInaktif)
		_ = f.SetCellValue(sheet, fmt.Sprintf("O%d", rowNum), r.NasibAkhir)
		_ = f.SetCellValue(sheet, fmt.Sprintf("P%d", rowNum), r.SKKAD)
		_ = f.SetCellValue(sheet, fmt.Sprintf("Q%d", rowNum), strings.ToUpper(r.KategoriArsip))
		_ = f.SetCellValue(sheet, fmt.Sprintf("R%d", rowNum), r.PenciptaArsip)

		// Styling cell data
		colsCenter := []string{"A", "B", "C", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q"}
		for _, col := range colsCenter {
			cellRef := fmt.Sprintf("%s%d", col, rowNum)
			_ = f.SetCellStyle(sheet, cellRef, cellRef, dataStyleCenter)
		}
		for _, col := range []string{"D", "R"} {
			cellRef := fmt.Sprintf("%s%d", col, rowNum)
			_ = f.SetCellStyle(sheet, cellRef, cellRef, dataStyleLeft)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseDAIImportFile membaca file upload baik berformat Excel (.xlsx) maupun CSV (.csv)
// menjadi matriks baris dan kolom ([][]string).
func ParseDAIImportFile(fileBytes []byte) ([][]string, error) {
	if len(fileBytes) == 0 {
		return nil, NewAppError(422, "File import kosong", nil)
	}

	// Cek apakah file adalah ZIP/XLSX (magic bytes PK\x03\x04)
	if len(fileBytes) >= 4 && fileBytes[0] == 'P' && fileBytes[1] == 'K' && fileBytes[2] == 0x03 && fileBytes[3] == 0x04 {
		return parseXLSXRows(fileBytes)
	}

	// Bersihkan UTF-8 BOM jika ada
	cleaned := bytes.TrimPrefix(fileBytes, []byte{0xEF, 0xBB, 0xBF})
	text := string(cleaned)

	// Deteksi delimiter CSV (koma atau titik-koma)
	firstLine := text
	if idx := strings.IndexAny(text, "\r\n"); idx >= 0 {
		firstLine = text[:idx]
	}
	commaCount := strings.Count(firstLine, ",")
	semiCount := strings.Count(firstLine, ";")

	csvReader := csv.NewReader(strings.NewReader(text))
	csvReader.FieldsPerRecord = -1
	csvReader.LazyQuotes = true
	if semiCount > commaCount {
		csvReader.Comma = ';'
	}

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, NewAppError(422, "Format file CSV/Excel tidak valid", nil)
	}
	return records, nil
}

func parseXLSXRows(fileBytes []byte) ([][]string, error) {
	xlsx, err := excelize.OpenReader(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, NewAppError(422, "Gagal membuka file Excel (.xlsx): file corrupt atau format tidak valid", err)
	}
	defer xlsx.Close()

	sheets := xlsx.GetSheetList()
	if len(sheets) == 0 {
		return nil, NewAppError(422, "Sheet tidak ditemukan di dalam file Excel (.xlsx)", nil)
	}

	targetSheet := sheets[0]
	for _, s := range sheets {
		if strings.EqualFold(s, "DAI") || strings.EqualFold(s, "Sheet1") {
			targetSheet = s
			break
		}
	}

	rows, err := xlsx.GetRows(targetSheet)
	if err != nil {
		return nil, NewAppError(422, "Gagal membaca baris pada sheet Excel", err)
	}
	return rows, nil
}
