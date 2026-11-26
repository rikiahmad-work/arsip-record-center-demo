package repository

import (
	"context"
	"encoding/json"
	"math"
	"strings"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
)

// UserRepository mendefinisikan interface untuk akses database user.
type UserRepository interface {
	FindAll(ctx context.Context, req dto.UserListRequest) ([]model.User, *dto.PaginatedResponse, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
	FindByNIP(ctx context.Context, nip string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) (*model.User, error)
	CountByRole(ctx context.Context) (map[string]int64, error)
	CountActiveInactive(ctx context.Context) (aktif, nonaktif int64, err error)
	CountTotal(ctx context.Context) (int64, error)
	GetPermissionsByRole(ctx context.Context, roleKey string) ([]string, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository membuat instance UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) FindAll(ctx context.Context, req dto.UserListRequest) ([]model.User, *dto.PaginatedResponse, error) {
	var users []model.User
	var total int64

	query := r.db.WithContext(ctx).Model(&model.User{})

	// Filter soft-deleted: "only" = hanya yang terhapus, "with" = termasuk yang terhapus
	switch req.Trashed {
	case "only":
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	case "with":
		query = query.Unscoped()
	}

	// Filter search
	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where(
			"name LIKE ? OR nip LIKE ? OR email LIKE ? OR username LIKE ? OR unit_kerja LIKE ? OR jabatan LIKE ?",
			search, search, search, search, search, search,
		)
	}

	// Filter role
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// Filter status
	if req.Status == "active" {
		query = query.Where("is_active = ?", true)
	} else if req.Status == "inactive" {
		query = query.Where("is_active = ?", false)
	}

	// Hitung total sebelum pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	// Default & whitelist sort column
	allowedSortColumns := map[string]bool{
		"id":            true,
		"name":          true,
		"nip":           true,
		"email":         true,
		"role":          true,
		"unit_kerja":    true,
		"jabatan":       true,
		"is_active":     true,
		"created_at":    true,
		"last_login_at": true,
	}
	if !allowedSortColumns[req.SortBy] {
		req.SortBy = "created_at"
	}
	if req.SortOrder != "asc" && req.SortOrder != "desc" {
		req.SortOrder = "desc"
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 || req.PerPage > 100 {
		req.PerPage = 10
	}

	offset := (req.Page - 1) * req.PerPage
	lastPage := int(math.Ceil(float64(total) / float64(req.PerPage)))
	from := offset + 1
	to := offset + len(users)

	if err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Offset(offset).
		Limit(req.PerPage).
		Find(&users).Error; err != nil {
		return nil, nil, err
	}

	// Recalculate 'to' after fetching
	if len(users) > 0 {
		to = offset + len(users)
	} else {
		to = 0
		from = 0
	}

	meta := &dto.PaginatedResponse{
		CurrentPage: req.Page,
		LastPage:    lastPage,
		PerPage:     req.PerPage,
		Total:       total,
		From:        from,
		To:          to,
	}

	return users, meta, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByNIP(ctx context.Context, nip string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("nip = ?", nip).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepositoryImpl) Restore(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = r.db.WithContext(ctx).Unscoped().Model(&user).Update("deleted_at", nil).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) CountByRole(ctx context.Context) (map[string]int64, error) {
	result := make(map[string]int64)

	type CountResult struct {
		Role  string
		Count int64
	}

	var counts []CountResult
	err := r.db.WithContext(ctx).Model(&model.User{}).
		Select("role, COUNT(*) as count").
		Group("role").
		Find(&counts).Error
	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		result[c.Role] = c.Count
	}

	return result, nil
}

func (r *userRepositoryImpl) CountActiveInactive(ctx context.Context) (aktif, nonaktif int64, err error) {
	err = r.db.WithContext(ctx).Model(&model.User{}).
		Select("COUNT(*) as count, is_active").
		Group("is_active").
		Find(&[]struct {
			Count    int64
			IsActive bool
		}{}).Error
	if err != nil {
		return 0, 0, err
	}

	// Query terpisah untuk menghitung aktif/nonaktif
	r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", true).Count(&aktif)
	r.db.WithContext(ctx).Model(&model.User{}).Where("is_active = ?", false).Count(&nonaktif)

	return aktif, nonaktif, nil
}

func (r *userRepositoryImpl) CountTotal(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
	return count, err
}

func (r *userRepositoryImpl) GetPermissionsByRole(ctx context.Context, roleKey string) ([]string, error) {
	norm := strings.ToLower(strings.TrimSpace(roleKey))
	normKey := strings.ReplaceAll(norm, "_", "-")

	var role model.RolePermission
	err := r.db.WithContext(ctx).Where("`key` = ? OR `key` = ?", normKey, norm).First(&role).Error
	if err == nil && role.Permissions != "" {
		var perms []string
		if jsonErr := json.Unmarshal([]byte(role.Permissions), &perms); jsonErr == nil {
			// Administrator selalu memiliki hak akses inti sistem
			if normKey == "administrator" || norm == "administrator" || normKey == "super-admin" {
				permMap := make(map[string]bool)
				for _, p := range perms {
					permMap[p] = true
				}
				permMap["dashboard.view"] = true
				permMap["users.view"] = true
				permMap["roles.manage"] = true

				finalList := make([]string, 0, len(permMap))
				for k := range permMap {
					finalList = append(finalList, k)
				}
				return finalList, nil
			}
			return perms, nil
		}
	}

	// Fallback jika belum tersimpan di tabel role_permissions
	switch normKey {
	case "administrator", "super-admin":
		return []string{
			"dashboard.view",
			"archives.view", "archives.create", "archives.edit", "archives.delete", "archives.import", "archives.export",
			"transfers.view", "transfers.create", "transfers.verify",
			"storage.view", "storage.manage",
			"loans.view", "loans.create", "loans.approve", "loans.return",
			"disposal.view", "disposal.propose", "disposal.approve", "disposal.execute",
			"reports.view", "reports.export",
			"users.view", "users.create", "users.edit", "users.delete", "roles.manage",
			"settings.view", "settings.manage",
		}, nil
	case "operator", "admin-arsip":
		return []string{
			"dashboard.view",
			"archives.view", "archives.create", "archives.edit", "archives.import", "archives.export",
			"transfers.view", "transfers.create", "transfers.verify",
			"storage.view", "storage.manage",
			"loans.view", "loans.create", "loans.approve", "loans.return",
			"disposal.view", "disposal.propose",
			"reports.view", "reports.export",
			"settings.view",
		}, nil
	case "approver":
		return []string{
			"dashboard.view",
			"archives.view", "archives.export",
			"transfers.view", "transfers.verify",
			"storage.view",
			"loans.view", "loans.approve",
			"disposal.view", "disposal.approve", "disposal.execute",
			"reports.view", "reports.export",
			"settings.view",
		}, nil
	case "viewer":
		return []string{
			"dashboard.view",
			"archives.view", "archives.export",
			"transfers.view",
			"storage.view",
			"loans.view",
			"disposal.view",
			"reports.view", "reports.export",
			"settings.view",
		}, nil
	default:
		return []string{}, nil
	}
}
