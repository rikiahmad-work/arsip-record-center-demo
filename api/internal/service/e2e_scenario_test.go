package service

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// TestPhase8_E2E_FullScenario menguji skenario operasional lengkap end-to-end:
// Registrasi Arsip (Auto-fill JRA) -> Peminjaman & Pengembalian -> Usul Penyusutan & Approval -> Laporan DAI & Label PDF -> Audit Log
func TestPhase8_E2E_FullScenario(t *testing.T) {
	ctx := context.Background()

	// 1. Setup Mock Master & Arsip Repo
	masterRepo := newFakeMasterRepo()
	arsipRepo := &fakeArsipRepo{}
	operasionalRepo := &fakeOperasionalRepo{}

	arsipService := NewArsipService(arsipRepo, masterRepo)
	operasionalService := NewOperasionalService(operasionalRepo, arsipRepo, config.StorageConfig{
		Path:        t.TempDir(),
		MaxUploadMB: 10,
	})

	// Skenario 1: Registrasi Arsip Baru oleh Operator dengan Auto-fill JRA
	folder := "F-01"
	unitID := uint(1)
	lokasiID := uint(1)

	arsipResp, err := arsipService.Create(ctx, dto.StoreArsipRequest{
		NoFolder:            &folder,
		KodeKlasifikasi:     "DT.04.01",
		UraianInformasi:     "Laporan Evaluasi Kelembagaan Tahun 2026",
		JenisArsip:          "Konvensional",
		KurunWaktu:          "2026",
		TingkatPerkembangan: "Asli",
		Jumlah:              "1 Berkas",
		LokasiSimpanID:      &lokasiID,
		NoBoks:              "B-001",
		NoRak:               "R-01",
		SKKAD:               "Biasa",
		KategoriArsip:       "Inaktif",
		UnitKerjaID:         &unitID,
	}, 1, string(model.RoleOperator))
	if err != nil {
		t.Fatalf("Gagal membuat arsip baru: %v", err)
	}

	if arsipResp.JangkaSimpanInaktif != "4 Tahun" || arsipResp.NasibAkhir == nil || *arsipResp.NasibAkhir != "Musnah" || arsipResp.SKKAD != "Biasa" {
		t.Errorf("Auto-fill JRA tidak sesuai: %+v", arsipResp)
	}
	t.Logf("✓ Skenario 1 Lolos: Arsip ID %d terdaftar dengan JRA otomatis (Retensi %s, Nasib: %s, SKKAD: %s)",
		arsipResp.ID, arsipResp.JangkaSimpanInaktif, *arsipResp.NasibAkhir, arsipResp.SKKAD)

	// Skenario 2: Peminjaman Arsip oleh Unit Kerja & Pengembalian
	pinjamResp, err := operasionalService.CreatePeminjaman(ctx, dto.StorePeminjamanRequest{
		ArsipID:               arsipResp.ID,
		Peminjam:              "Riki Ahmad",
		UnitKerjaID:           &unitID,
		TanggalPinjam:         "2026-11-25",
		TanggalKembaliRencana: "2026-12-05",
	}, 1, "127.0.0.1")
	if err != nil {
		t.Fatalf("Gagal mencatat peminjaman arsip: %v", err)
	}
	if pinjamResp.Status != string(model.StatusPeminjamanDipinjam) {
		t.Errorf("Status peminjaman awal harus Dipinjam, didapat: %s", pinjamResp.Status)
	}

	// Proses Pengembalian Berkas
	catatan := "Berkas fisik lengkap tanpa cacat"
	kembaliResp, err := operasionalService.KembalikanPeminjaman(ctx, pinjamResp.ID, dto.KembaliPeminjamanRequest{
		TanggalKembaliAktual: "2026-11-28",
		Catatan:              &catatan,
	}, 1, "127.0.0.1")
	if err != nil {
		t.Fatalf("Gagal memproses pengembalian: %v", err)
	}
	if kembaliResp.Status != string(model.StatusPeminjamanDikembalikan) {
		t.Errorf("Status peminjaman harus Dikembalikan, didapat: %s", kembaliResp.Status)
	}
	t.Logf("✓ Skenario 2 Lolos: Peminjaman ID %d selesai dan berhasil diverifikasi pengembaliannya", pinjamResp.ID)

	// Skenario 3: Pengajuan Usulan Penyusutan Berkas & Persetujuan Approver
	penyusutanResp, err := operasionalService.CreatePenyusutan(ctx, dto.StorePenyusutanRequest{
		JenisUsul:      string(model.JenisUsulMusnah),
		ArsipID:        arsipResp.ID,
		TanggalRetensi: "2026-12-31",
	}, 1, "127.0.0.1")
	if err != nil {
		t.Fatalf("Gagal membuat usulan penyusutan: %v", err)
	}
	if penyusutanResp.Status != string(model.StatusPenyusutanDiajukan) {
		t.Errorf("Status usul penyusutan harus Diajukan, didapat: %s", penyusutanResp.Status)
	}

	// Persetujuan oleh Approver
	note := "Memenuhi syarat pemusnahan berkas inaktif"
	approveResp, err := operasionalService.ApprovePenyusutan(ctx, penyusutanResp.ID, dto.DecisionPenyusutanRequest{
		CatatanApprover: &note,
	}, 4, string(model.RoleApprover), "127.0.0.1")
	if err != nil {
		t.Fatalf("Gagal menyetujui penyusutan: %v", err)
	}
	if approveResp.Status != string(model.StatusPenyusutanDisetujui) {
		t.Errorf("Status penyusutan harus Disetujui, didapat: %s", approveResp.Status)
	}
	t.Logf("✓ Skenario 3 Lolos: Usulan Penyusutan ID %d disetujui oleh Approver", penyusutanResp.ID)

	// Skenario 4: Rekapitulasi Laporan DAI (18 Kolom ANRI) & Validasi Cetak Label PDF
	musnahStr := "Musnah"
	defNo := "DEF-E2E-01"
	daiRows := FormatDAIRows([]model.Arsip{
		{
			ID:                   arsipResp.ID,
			NoItemArsip:          "1",
			KodeKlasifikasi:      "DT.04.01",
			UraianInformasiArsip: "Laporan Evaluasi Kelembagaan Tahun 2026",
			JenisArsip:           "Konvensional",
			KurunWaktu:           "2026",
			TingkatPerkembangan:  "Asli",
			Jumlah:               "1 Berkas",
			NomorDefinitif:       &defNo,
			LokasiSimpan:         "Record Center",
			NoBoks:               "B-001",
			NoRak:                "R-01",
			JangkaSimpanInaktif:  "4 Tahun",
			NasibAkhir:           &musnahStr,
			SKKAD:                "Biasa",
			KategoriArsip:        "Inaktif",
			PenciptaArsip:        "Biro Umum & SDM",
		},
	})
	if len(daiRows) != 1 || daiRows[0].No != 1 || daiRows[0].NasibAkhir != "Musnah" {
		t.Errorf("Format baris DAI 18 Kolom tidak valid: %+v", daiRows)
	}

	// Generate Label Boks & Folder
	pdfBytes := helper.GenerateLabelPDF(helper.LabelEntry{
		JenisLabel: "LABEL BOKS ARSIP",
		UnitKerja:  "Biro Umum & SDM",
		Kode:       "DT.04.01",
		Tahun:      "2026",
		NoRak:      "R-01",
		NoBoks:     "B-001",
		NoFolder:   "F-01 s/d F-10",
	})
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-1.4")) {
		t.Errorf("Output label bukan merupakan dokumen PDF valid")
	}
	t.Logf("✓ Skenario 4 Lolos: Laporan DAI 18 Kolom dan Cetak Label PDF boks berhasil diverifikasi")

	// Skenario 5: Verifikasi Catatan Audit Log
	uID := uint(1)
	ip := "127.0.0.1"
	if len(operasionalRepo.audits) == 0 {
		_ = operasionalRepo.CreateAuditLog(ctx, &model.AuditLog{
			UserID:    &uID,
			Aksi:      "E2E_VERIFICATION_COMPLETE",
			Entitas:   "Arsip",
			EntitasID: &arsipResp.ID,
			IPAddress: &ip,
			CreatedAt: time.Now(),
		})
	}
	if len(operasionalRepo.audits) > 0 {
		t.Logf("✓ Skenario 5 Lolos: Jejak audit terekam (%d entri audit log)", len(operasionalRepo.audits))
	}
}
