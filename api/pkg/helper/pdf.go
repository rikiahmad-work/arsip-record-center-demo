package helper

import (
	"bytes"
	"fmt"
	"strings"
)

// LabelEntry merepresentasikan data kartu label arsip (Boks / Folder / Inaktif).
type LabelEntry struct {
	JenisLabel  string
	Instansi    string
	UnitKerja   string
	Kode        string
	Tahun       string
	NoRak       string
	NoBoks      string
	NoFolder    string
	RentangItem string
	Uraian      string
}

func pdfEscape(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `(`, `\(`, `)`, `\)`)
	return r.Replace(s)
}

// GenerateLabelPDF menghasilkan dokumen PDF 1.4 standar (tanpa library eksternal)
// untuk Cetak Label Boks, Label Folder, atau Label Arsip Inaktif KemenPPPA.
func GenerateLabelPDF(entry LabelEntry) []byte {
	if entry.Instansi == "" {
		entry.Instansi = "KEMENTERIAN PEMBERDAYAAN PEREMPUAN DAN PERLINDUNGAN ANAK"
	}
	if entry.UnitKerja == "" {
		entry.UnitKerja = "BIRO UMUM DAN SDM"
	}
	if entry.Tahun == "" {
		entry.Tahun = "2026"
	}
	if entry.NoBoks == "" {
		entry.NoBoks = "1"
	}
	if entry.NoFolder == "" {
		entry.NoFolder = "1 - 10"
	}
	if entry.NoRak == "" {
		entry.NoRak = "R-01"
	}
	if entry.Kode == "" {
		entry.Kode = "DT.04.01"
	}
	if entry.RentangItem == "" {
		entry.RentangItem = "1 - 20"
	}

	var stream bytes.Buffer
	// Bingkai luar & header kartu label
	stream.WriteString("1.5 w\n")
	stream.WriteString("40 480 515 280 re S\n")
	stream.WriteString("40 700 515 60 re S\n")
	stream.WriteString("40 540 515 160 re S\n")
	stream.WriteString("380 480 175 160 re S\n")

	// Teks Header Instansi & Jenis Label
	stream.WriteString("BT\n")
	stream.WriteString("/F1 11 Tf\n")
	stream.WriteString(fmt.Sprintf("55 738 Td (%s) Tj\n", pdfEscape(entry.Instansi)))
	stream.WriteString("/F1 13 Tf\n")
	stream.WriteString(fmt.Sprintf("0 -22 Td (RECORDS CENTER - %s) Tj\n", pdfEscape(strings.ToUpper(entry.JenisLabel))))
	stream.WriteString("ET\n")

	// Isi Metadata Utama (Unit Pencipta, Kode Klasifikasi, Tahun, Uraian)
	stream.WriteString("BT\n")
	stream.WriteString("/F1 11 Tf\n")
	stream.WriteString(fmt.Sprintf("55 672 Td (UNIT PENCIPTA : %s) Tj\n", pdfEscape(strings.ToUpper(entry.UnitKerja))))
	stream.WriteString(fmt.Sprintf("0 -26 Td (KODE KLASIFIKASI : %s) Tj\n", pdfEscape(entry.Kode)))
	stream.WriteString(fmt.Sprintf("0 -26 Td (KURUN WAKTU / TAHUN : %s) Tj\n", pdfEscape(entry.Tahun)))
	stream.WriteString(fmt.Sprintf("0 -26 Td (RENTANG FOLDER / ITEM : %s (Item: %s)) Tj\n", pdfEscape(entry.NoFolder), pdfEscape(entry.RentangItem)))
	if entry.Uraian != "" {
		uraianShort := entry.Uraian
		if len(uraianShort) > 55 {
			uraianShort = uraianShort[:55] + "..."
		}
		stream.WriteString(fmt.Sprintf("0 -26 Td (URAIAN : %s) Tj\n", pdfEscape(uraianShort)))
	}
	stream.WriteString("ET\n")

	// Kotak Nomor Boks & Rak di sisi kanan/bawah
	stream.WriteString("BT\n")
	stream.WriteString("/F1 12 Tf\n")
	stream.WriteString(fmt.Sprintf("55 505 Td (NO. RAK : %s   |   NO. FOLDER : %s) Tj\n", pdfEscape(entry.NoRak), pdfEscape(entry.NoFolder)))
	stream.WriteString("ET\n")

	stream.WriteString("BT\n")
	stream.WriteString("/F1 12 Tf\n")
	stream.WriteString("415 610 Td (NOMOR BOKS) Tj\n")
	stream.WriteString("/F1 28 Tf\n")
	stream.WriteString(fmt.Sprintf("15 -50 Td (%s) Tj\n", pdfEscape(entry.NoBoks)))
	stream.WriteString("ET\n")

	contentBytes := stream.Bytes()

	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n")

	offsets := make([]int, 6)

	// Object 1: Catalog
	offsets[1] = pdf.Len()
	pdf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	// Object 2: Pages
	offsets[2] = pdf.Len()
	pdf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// Object 3: Page
	offsets[3] = pdf.Len()
	pdf.WriteString("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 595 842] /Resources << /Font << /F1 4 0 R >> >> /Contents 5 0 R >>\nendobj\n")

	// Object 4: Font Helvetica-Bold
	offsets[4] = pdf.Len()
	pdf.WriteString("4 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>\nendobj\n")

	// Object 5: Content Stream
	offsets[5] = pdf.Len()
	pdf.WriteString(fmt.Sprintf("5 0 obj\n<< /Length %d >>\nstream\n", len(contentBytes)))
	pdf.Write(contentBytes)
	pdf.WriteString("\nendstream\nendobj\n")

	xrefOffset := pdf.Len()
	pdf.WriteString("xref\n0 6\n")
	pdf.WriteString("0000000000 65535 f \n")
	for i := 1; i <= 5; i++ {
		pdf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	pdf.WriteString(fmt.Sprintf("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n%d\n%%%%EOF\n", xrefOffset))

	return pdf.Bytes()
}
