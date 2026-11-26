package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// PermissionGroupItem merepresentasikan 1 grup checklist hak akses di UI (persis 9 grup di referensi Laravel).
type PermissionGroupItem struct {
	Group string            `json:"group"`
	Icon  string            `json:"icon"`
	Items map[string]string `json:"items"`
}

// DefaultPermissionGroups adalah 9 grup checklist permission standar SiARSIP.
var DefaultPermissionGroups = []PermissionGroupItem{
	{
		Group: "Dashboard",
		Icon:  "ti-dashboard",
		Items: map[string]string{
			"dashboard.view": "Lihat Halaman Dashboard",
		},
	},
	{
		Group: "Data Arsip",
		Icon:  "ti-folders",
		Items: map[string]string{
			"archives.view":   "Lihat Daftar Arsip",
			"archives.create": "Tambah Data Arsip",
			"archives.edit":   "Ubah Data Arsip",
			"archives.delete": "Hapus Data Arsip",
			"archives.import": "Impor Data Excel",
			"archives.export": "Ekspor Data Arsip",
		},
	},
	{
		Group: "Pemindahan Arsip",
		Icon:  "ti-arrows-transfer-down",
		Items: map[string]string{
			"transfers.view":   "Lihat Daftar Pemindahan",
			"transfers.create": "Ajukan Usulan Pemindahan",
			"transfers.verify": "Verifikasi & Terbitkan BAST",
		},
	},
	{
		Group: "Lokasi Simpan",
		Icon:  "ti-box",
		Items: map[string]string{
			"storage.view":   "Lihat Lokasi Rak & Boks",
			"storage.manage": "Kelola Rak & Penempatan Boks",
		},
	},
	{
		Group: "Peminjaman",
		Icon:  "ti-file-symlink",
		Items: map[string]string{
			"loans.view":    "Lihat Daftar Peminjaman",
			"loans.create":  "Ajukan Permohonan Pinjam",
			"loans.approve": "Persetujuan Peminjaman",
			"loans.return":  "Verifikasi Pengembalian",
		},
	},
	{
		Group: "Penyusutan",
		Icon:  "ti-archive-off",
		Items: map[string]string{
			"disposal.view":    "Lihat Daftar Penyusutan",
			"disposal.propose": "Ajukan Usul Musnah / Serah",
			"disposal.approve": "Penilaian & Persetujuan JRA",
			"disposal.execute": "Pengesahan Berita Acara",
		},
	},
	{
		Group: "Laporan & Label",
		Icon:  "ti-report",
		Items: map[string]string{
			"reports.view":   "Lihat Rekapitulasi Laporan & DAI",
			"reports.export": "Unduh Laporan & Cetak Label PDF",
		},
	},
	{
		Group: "Manajemen User",
		Icon:  "ti-users",
		Items: map[string]string{
			"users.view":   "Lihat Daftar Pengguna",
			"users.create": "Tambah Pengguna",
			"users.edit":   "Ubah Data & Status Pengguna",
			"users.delete": "Hapus Pengguna",
			"roles.manage": "Kelola Peran & Hak Akses",
		},
	},
	{
		Group: "Master Data",
		Icon:  "ti-settings",
		Items: map[string]string{
			"settings.view":   "Lihat Master Data",
			"settings.manage": "Kelola Klasifikasi & Unit Kerja",
		},
	},
}

// RoleItemResponse merepresentasikan 1 role beserta permissions dan jumlah user.
type RoleItemResponse struct {
	ID          uint     `json:"id"`
	Key         string   `json:"key"`
	Name        string   `json:"name"`
	IsSystem    bool     `json:"is_system"`
	UsersCount  int64    `json:"users_count"`
	Permissions []string `json:"permissions"`
}

// RoleHandler menangani CRUD Role & Dynamic RBAC Permission Matrix (/api/roles).
type RoleHandler struct {
	db *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

func allPermissionKeys() []string {
	var all []string
	for _, g := range DefaultPermissionGroups {
		for k := range g.Items {
			all = append(all, k)
		}
	}
	return all
}

// EnsureDefaultRoles memastikan 4 role kanonik selalu ada di tabel role_permissions.
func EnsureDefaultRoles(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	allPerms, _ := json.Marshal(allPermissionKeys())
	operatorPerms, _ := json.Marshal([]string{
		"dashboard.view",
		"archives.view", "archives.create", "archives.edit", "archives.import", "archives.export",
		"transfers.view", "transfers.create", "transfers.verify",
		"storage.view", "storage.manage",
		"loans.view", "loans.create", "loans.approve", "loans.return",
		"disposal.view", "disposal.propose",
		"reports.view", "reports.export",
		"settings.view",
	})
	approverPerms, _ := json.Marshal([]string{
		"dashboard.view",
		"archives.view", "archives.export",
		"transfers.view", "transfers.verify",
		"storage.view",
		"loans.view", "loans.approve",
		"disposal.view", "disposal.approve", "disposal.execute",
		"reports.view", "reports.export",
		"settings.view",
	})
	viewerPerms, _ := json.Marshal([]string{
		"dashboard.view",
		"archives.view", "archives.export",
		"transfers.view",
		"storage.view",
		"loans.view",
		"disposal.view",
		"reports.view", "reports.export",
		"settings.view",
	})

	defaults := []model.RolePermission{
		{Key: "administrator", Name: "Administrator", IsSystem: true, Permissions: string(allPerms)},
		{Key: "operator", Name: "Operator", IsSystem: true, Permissions: string(operatorPerms)},
		{Key: "approver", Name: "Approver", IsSystem: true, Permissions: string(approverPerms)},
		{Key: "viewer", Name: "Viewer", IsSystem: true, Permissions: string(viewerPerms)},
	}

	for _, d := range defaults {
		var existing model.RolePermission
		if err := db.Where("`key` = ?", d.Key).First(&existing).Error; err != nil {
			if errCreate := db.Create(&d).Error; errCreate != nil {
				return errCreate
			}
		}
	}
	return nil
}

// Index menangani GET /api/roles (mengembalikan daftar role + 9 grup permission matrix).
func (h *RoleHandler) Index(c *gin.Context) {
	_ = EnsureDefaultRoles(h.db)

	var roles []model.RolePermission
	if err := h.db.Order("id ASC").Find(&roles).Error; err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data role")
		return
	}

	// Hitung jumlah user per role
	type roleCount struct {
		Role  string
		Total int64
	}
	var counts []roleCount
	_ = h.db.Model(&model.User{}).Select("role, COUNT(*) as total").Group("role").Scan(&counts).Error
	countMap := make(map[string]int64)
	for _, rc := range counts {
		countMap[rc.Role] = rc.Total
	}

	roleResponses := make([]RoleItemResponse, 0, len(roles))
	for _, r := range roles {
		var perms []string
		_ = json.Unmarshal([]byte(r.Permissions), &perms)
		if perms == nil {
			perms = []string{}
		}
		roleResponses = append(roleResponses, RoleItemResponse{
			ID:          r.ID,
			Key:         r.Key,
			Name:        r.Name,
			IsSystem:    r.IsSystem,
			UsersCount:  countMap[r.Key],
			Permissions: perms,
		})
	}

	helper.SuccessResponse(c, "", gin.H{
		"roles":             roleResponses,
		"permission_groups": DefaultPermissionGroups,
	})
}

// Store menangani POST /api/roles (tambah peran kustom baru).
func (h *RoleHandler) Store(c *gin.Context) {
	var req struct {
		Name        string   `json:"name"        binding:"required,max=100"`
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	name := strings.TrimSpace(req.Name)
	key := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	var existing model.RolePermission
	if err := h.db.Where("`key` = ? OR name = ?", key, name).First(&existing).Error; err == nil {
		helper.ErrorResponse(c, http.StatusConflict, "Nama peran tersebut sudah terdaftar")
		return
	}

	if len(req.Permissions) == 0 {
		req.Permissions = []string{"dashboard.view"}
	}
	permsJSON, _ := json.Marshal(req.Permissions)

	role := model.RolePermission{
		Key:         key,
		Name:        name,
		IsSystem:    false,
		Permissions: string(permsJSON),
	}
	if err := h.db.Create(&role).Error; err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menambahkan peran baru")
		return
	}

	helper.CreatedResponse(c, "Peran "+role.Name+" berhasil ditambahkan", RoleItemResponse{
		ID:          role.ID,
		Key:         role.Key,
		Name:        role.Name,
		IsSystem:    role.IsSystem,
		UsersCount:  0,
		Permissions: req.Permissions,
	})
}

// UpdatePermissions menangani PUT /api/roles/:id/permissions.
func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID peran tidak valid")
		return
	}

	var req struct {
		Permissions []string `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, map[string]string{"validation": err.Error()})
		return
	}

	var role model.RolePermission
	if err := h.db.First(&role, uint(id)).Error; err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Peran tidak ditemukan")
		return
	}

	// Pastikan permission inti Administrator tidak pernah hilang (sama seperti RoleController.php di Laravel)
	permSet := make(map[string]bool)
	for _, p := range req.Permissions {
		permSet[p] = true
	}
	if role.Key == "administrator" {
		permSet["dashboard.view"] = true
		permSet["users.view"] = true
		permSet["roles.manage"] = true
	}

	finalPerms := make([]string, 0, len(permSet))
	for k := range permSet {
		finalPerms = append(finalPerms, k)
	}
	permsJSON, _ := json.Marshal(finalPerms)

	role.Permissions = string(permsJSON)
	if err := h.db.Save(&role).Error; err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan pengaturan hak akses")
		return
	}

	helper.SuccessResponse(c, "Pengaturan hak akses untuk peran "+role.Name+" berhasil disimpan", RoleItemResponse{
		ID:          role.ID,
		Key:         role.Key,
		Name:        role.Name,
		IsSystem:    role.IsSystem,
		Permissions: finalPerms,
	})
}

// Destroy menangani DELETE /api/roles/:id (hanya peran kustom yang bisa dihapus).
func (h *RoleHandler) Destroy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "Parameter ID peran tidak valid")
		return
	}

	var role model.RolePermission
	if err := h.db.First(&role, uint(id)).Error; err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "Peran tidak ditemukan")
		return
	}

	if role.IsSystem {
		helper.ErrorResponse(c, http.StatusBadRequest, "Peran utama "+role.Name+" tidak dapat dihapus")
		return
	}

	if err := h.db.Delete(&role).Error; err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus peran")
		return
	}

	helper.SuccessResponse(c, "Peran "+role.Name+" berhasil dihapus", nil)
}
