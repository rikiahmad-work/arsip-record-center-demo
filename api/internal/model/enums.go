package model

// ================================================================
// 1. UserRole — Role Pengguna Kanonik (4 Role)
// ================================================================

type UserRole string

const (
	RoleAdministrator UserRole = "administrator"
	RoleOperator      UserRole = "operator"
	RoleViewer        UserRole = "viewer"
	RoleApprover      UserRole = "approver"

	// Alias kompatibilitas kode lama (dipetakan ke role kanonik)
	RoleSuperAdmin UserRole = RoleAdministrator
	RoleAdminArsip UserRole = RoleOperator
)

// LegacyRoleMap memetakan nilai role lama ke role kanonik.
var LegacyRoleMap = map[string]UserRole{
	"super-admin": RoleAdministrator,
	"admin-arsip": RoleOperator,
}

// NormalizeRole mengembalikan role kanonik dari string role (termasuk nilai legacy).
func NormalizeRole(raw string) UserRole {
	if mapped, ok := LegacyRoleMap[raw]; ok {
		return mapped
	}
	return UserRole(raw)
}

// Label mengembalikan label role dalam bahasa Indonesia.
func (r UserRole) Label() string {
	switch NormalizeRole(string(r)) {
	case RoleAdministrator:
		return "Administrator"
	case RoleOperator:
		return "Operator"
	case RoleViewer:
		return "Viewer"
	case RoleApprover:
		return "Approver"
	default:
		return string(r)
	}
}

// IsValid memeriksa apakah role termasuk 4 role kanonik.
func (r UserRole) IsValid() bool {
	switch NormalizeRole(string(r)) {
	case RoleAdministrator, RoleOperator, RoleViewer, RoleApprover:
		return true
	default:
		return false
	}
}

// CanWrite menandakan apakah role boleh menambah/mengubah data arsip.
func (r UserRole) CanWrite() bool {
	norm := NormalizeRole(string(r))
	return norm == RoleAdministrator || norm == RoleOperator
}

// CanApprove menandakan apakah role boleh menyetujui/menolak penyusutan.
func (r UserRole) CanApprove() bool {
	norm := NormalizeRole(string(r))
	return norm == RoleAdministrator || norm == RoleApprover
}

// AllRoles mengembalikan seluruh role kanonik.
func AllRoles() []UserRole {
	return []UserRole{RoleAdministrator, RoleOperator, RoleViewer, RoleApprover}
}

// ================================================================
// 2. SKKAD — Sistem Klasifikasi Keamanan dan Akses Arsip Dinamis
// ================================================================

type SKKAD string

const (
	SKKADBiasa    SKKAD = "Biasa"
	SKKADTerbatas SKKAD = "Terbatas"
	SKKADRahasia  SKKAD = "Rahasia"
)

func (s SKKAD) IsValid() bool {
	switch s {
	case SKKADBiasa, SKKADTerbatas, SKKADRahasia:
		return true
	default:
		return false
	}
}

// ================================================================
// 3. KategoriArsip — Status / Kategori Arsip Record Center
// ================================================================

type KategoriArsip string

const (
	KategoriInaktifVal    KategoriArsip = "Inaktif"
	KategoriVitalVal      KategoriArsip = "Vital"
	KategoriTerjagaVal    KategoriArsip = "Terjaga"
	KategoriUsulMusnahVal KategoriArsip = "Usul Musnah"
	KategoriUsulSerahVal  KategoriArsip = "Usul Serah"
)

func (k KategoriArsip) IsValid() bool {
	switch k {
	case KategoriInaktifVal, KategoriVitalVal, KategoriTerjagaVal, KategoriUsulMusnahVal, KategoriUsulSerahVal:
		return true
	default:
		return false
	}
}

// ================================================================
// 4. NasibAkhir — Nasib Akhir Retensi JRA
// ================================================================

type NasibAkhir string

const (
	NasibAkhirMusnah   NasibAkhir = "Musnah"
	NasibAkhirPermanen NasibAkhir = "Permanen"
)

func (n NasibAkhir) IsValid() bool {
	switch n {
	case NasibAkhirMusnah, NasibAkhirPermanen:
		return true
	default:
		return false
	}
}

// ================================================================
// 5. JenisArsip — Media Arsip
// ================================================================

type JenisArsip string

const (
	JenisArsipKonvensional JenisArsip = "Konvensional"
	JenisArsipElektronik   JenisArsip = "Elektronik"
)

func (j JenisArsip) IsValid() bool {
	switch j {
	case JenisArsipKonvensional, JenisArsipElektronik:
		return true
	default:
		return false
	}
}

// ================================================================
// 6. TingkatPerkembangan — Tingkat Perkembangan Dokumen Arsip
// ================================================================

type TingkatPerkembangan string

const (
	TingkatPerkembanganAsli     TingkatPerkembangan = "Asli"
	TingkatPerkembanganCopy     TingkatPerkembangan = "Copy"
	TingkatPerkembanganSalinan  TingkatPerkembangan = "Salinan"
	TingkatPerkembanganTembusan TingkatPerkembangan = "Tembusan"
)

func (t TingkatPerkembangan) IsValid() bool {
	switch t {
	case TingkatPerkembanganAsli, TingkatPerkembanganCopy, TingkatPerkembanganSalinan, TingkatPerkembanganTembusan:
		return true
	default:
		return false
	}
}

// ================================================================
// 7. StatusArsip — Status Penyimpanan Form Arsip
// ================================================================

type StatusArsip string

const (
	StatusArsipDraft StatusArsip = "draft"
	StatusArsipFinal StatusArsip = "final"
)

func (s StatusArsip) IsValid() bool {
	return s == StatusArsipDraft || s == StatusArsipFinal
}

// ================================================================
// 8. StatusPeminjaman — Status Peminjaman Arsip
// ================================================================

type StatusPeminjaman string

const (
	StatusPeminjamanDipinjam     StatusPeminjaman = "Dipinjam"
	StatusPeminjamanDikembalikan StatusPeminjaman = "Dikembalikan"
	StatusPeminjamanTerlambat    StatusPeminjaman = "Terlambat"
)

func (s StatusPeminjaman) IsValid() bool {
	switch s {
	case StatusPeminjamanDipinjam, StatusPeminjamanDikembalikan, StatusPeminjamanTerlambat:
		return true
	default:
		return false
	}
}

// ================================================================
// 9. JenisUsul & StatusPenyusutan — Penyusutan Arsip
// ================================================================

type JenisUsul string

const (
	JenisUsulMusnah JenisUsul = "Musnah"
	JenisUsulSerah  JenisUsul = "Serah"
)

func (j JenisUsul) IsValid() bool {
	return j == JenisUsulMusnah || j == JenisUsulSerah
}

type StatusPenyusutan string

const (
	StatusPenyusutanDiajukan  StatusPenyusutan = "Diajukan"
	StatusPenyusutanDisetujui StatusPenyusutan = "Disetujui"
	StatusPenyusutanDitolak   StatusPenyusutan = "Ditolak"
	StatusPenyusutanSelesai   StatusPenyusutan = "Selesai"
)

func (s StatusPenyusutan) IsValid() bool {
	switch s {
	case StatusPenyusutanDiajukan, StatusPenyusutanDisetujui, StatusPenyusutanDitolak, StatusPenyusutanSelesai:
		return true
	default:
		return false
	}
}
