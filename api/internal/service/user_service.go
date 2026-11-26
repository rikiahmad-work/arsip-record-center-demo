package service

import (
	"context"
	"errors"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// UserService mendefinisikan interface untuk business logic manajemen user.
type UserService interface {
	GetAll(ctx context.Context, req dto.UserListRequest) ([]dto.UserResponse, *dto.PaginatedResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	Create(ctx context.Context, req dto.StoreUserRequest, createdBy uint) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, req dto.UpdateUserRequest, currentUserID ...uint) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint, currentUserID uint) error
	Restore(ctx context.Context, id uint) (*dto.UserResponse, error)
	GetStats(ctx context.Context) (*dto.UserStatsResponse, error)
}

type userServiceImpl struct {
	userRepo repository.UserRepository
	hashHelp *helper.HashHelper
	mfaRepo  repository.MfaRepository
}

// NewUserService membuat instance UserService.
func NewUserService(userRepo repository.UserRepository, hashHelp *helper.HashHelper, mfaRepo ...repository.MfaRepository) UserService {
	var repo repository.MfaRepository
	if len(mfaRepo) > 0 {
		repo = mfaRepo[0]
	}
	return &userServiceImpl{
		userRepo: userRepo,
		hashHelp: hashHelp,
		mfaRepo:  repo,
	}
}

func (s *userServiceImpl) GetAll(ctx context.Context, req dto.UserListRequest) ([]dto.UserResponse, *dto.PaginatedResponse, error) {
	users, meta, err := s.userRepo.FindAll(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	var mfaMap map[uint]bool
	if s.mfaRepo != nil && len(users) > 0 {
		ids := make([]uint, len(users))
		for i, u := range users {
			ids[i] = u.ID
		}
		mfaMap, _ = s.mfaRepo.GetMfaStatusesByUserIDs(ctx, ids)
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		resp := s.toResponse(user)
		if mfaMap != nil {
			resp.MfaEnabled = mfaMap[user.ID]
		}
		responses = append(responses, resp)
	}

	return responses, meta, nil
}

func (s *userServiceImpl) GetByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	resp := s.toResponse(*user)
	if s.mfaRepo != nil {
		resp.MfaEnabled = s.mfaRepo.IsUserMfaEnabled(ctx, id)
	}
	return &resp, nil
}

func (s *userServiceImpl) Create(ctx context.Context, req dto.StoreUserRequest, createdBy uint) (*dto.UserResponse, error) {
	if req.Username == "" {
		req.Username = req.NIP
	}

	// Cek unique Username
	existingUsername, _ := s.userRepo.FindByUsername(ctx, req.Username)
	if existingUsername != nil {
		return nil, helper.NewAppError(409, "Username sudah terdaftar", nil)
	}

	// Cek unique NIP
	existingNIP, _ := s.userRepo.FindByNIP(ctx, req.NIP)
	if existingNIP != nil {
		return nil, helper.ErrDuplicateNIP
	}

	// Cek unique Email
	existingEmail, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existingEmail != nil {
		return nil, helper.ErrDuplicateEmail
	}

	// Default is_sso
	isSSO := false
	if req.IsSSO != nil {
		isSSO = *req.IsSSO
	}

	// Hash password (wajib untuk manual, opsional untuk SSO)
	var hashedPassword string
	if req.Password != "" {
		pwd, err := s.hashHelp.Hash(req.Password)
		if err != nil {
			return nil, helper.NewAppError(500, "Gagal memproses password", err)
		}
		hashedPassword = pwd
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var unitKerja *string
	if req.UnitKerja != "" {
		unitKerja = &req.UnitKerja
	}

	var jabatan *string
	if req.Jabatan != "" {
		jabatan = &req.Jabatan
	}

	user := &model.User{
		Username:  req.Username,
		Name:      req.Name,
		NIP:       req.NIP,
		Email:     req.Email,
		Password:  hashedPassword,
		Role:      req.Role,
		UnitKerja: unitKerja,
		Jabatan:   jabatan,
		IsSSO:     isSSO,
		IsActive:  isActive,
		CreatedBy: &createdBy,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	resp := s.toResponse(*user)
	return &resp, nil
}

func (s *userServiceImpl) Update(ctx context.Context, id uint, req dto.UpdateUserRequest, currentUserID ...uint) (*dto.UserResponse, error) {
	// Cegah user menonaktifkan akun miliknya sendiri
	if len(currentUserID) > 0 && id == currentUserID[0] && req.IsActive != nil && !*req.IsActive {
		return nil, helper.ErrSelfDeactivate
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	if req.Username == "" {
		req.Username = req.NIP
	}

	// Cek unique Username (kecuali milik user sendiri)
	if req.Username != user.Username {
		existingUsername, _ := s.userRepo.FindByUsername(ctx, req.Username)
		if existingUsername != nil {
			return nil, helper.NewAppError(409, "Username sudah terdaftar", nil)
		}
	}

	// Cek unique NIP (kecuali milik user sendiri)
	if req.NIP != user.NIP {
		existingNIP, _ := s.userRepo.FindByNIP(ctx, req.NIP)
		if existingNIP != nil {
			return nil, helper.ErrDuplicateNIP
		}
	}

	// Cek unique Email (kecuali milik user sendiri)
	if req.Email != user.Email {
		existingEmail, _ := s.userRepo.FindByEmail(ctx, req.Email)
		if existingEmail != nil {
			return nil, helper.ErrDuplicateEmail
		}
	}

	// Update field
	user.Username = req.Username
	user.Name = req.Name
	user.NIP = req.NIP
	user.Email = req.Email
	user.Role = req.Role

	if req.UnitKerja != "" {
		user.UnitKerja = &req.UnitKerja
	} else {
		user.UnitKerja = nil
	}

	if req.Jabatan != "" {
		user.Jabatan = &req.Jabatan
	} else {
		user.Jabatan = nil
	}

	if req.IsSSO != nil {
		user.IsSSO = *req.IsSSO
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	// Update password jika diisi
	if req.Password != "" {
		hashedPassword, err := s.hashHelp.Hash(req.Password)
		if err != nil {
			return nil, helper.NewAppError(500, "Gagal memproses password", err)
		}
		user.Password = hashedPassword
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	resp := s.toResponse(*user)
	return &resp, nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id uint, currentUserID uint) error {
	// Cegah user menghapus dirinya sendiri
	if id == currentUserID {
		return helper.ErrSelfDelete
	}

	// Cek user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *userServiceImpl) Restore(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.Restore(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	resp := s.toResponse(*user)
	return &resp, nil
}

func (s *userServiceImpl) GetStats(ctx context.Context) (*dto.UserStatsResponse, error) {
	total, err := s.userRepo.CountTotal(ctx)
	if err != nil {
		return nil, err
	}

	aktif, nonaktif, err := s.userRepo.CountActiveInactive(ctx)
	if err != nil {
		return nil, err
	}

	byRole, err := s.userRepo.CountByRole(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.UserStatsResponse{
		Total:         total,
		Aktif:         aktif,
		Nonaktif:      nonaktif,
		Administrator: byRole[string(model.RoleAdministrator)],
		Operator:      byRole[string(model.RoleOperator)],
		Viewer:        byRole[string(model.RoleViewer)],
		Approver:      byRole[string(model.RoleApprover)],
		SuperAdmin:    byRole[string(model.RoleAdministrator)],
		AdminArsip:    byRole[string(model.RoleOperator)],
	}, nil
}

func (s *userServiceImpl) toResponse(user model.User) dto.UserResponse {
	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		t := user.DeletedAt.Time
		deletedAt = &t
	}

	return dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		NIP:         user.NIP,
		Email:       user.Email,
		Role:        user.Role,
		RoleLabel:   model.UserRole(user.Role).Label(),
		UnitKerja:   user.UnitKerja,
		Jabatan:     user.Jabatan,
		IsSSO:       user.IsSSO,
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   deletedAt,
		CreatedBy:   user.CreatedBy,
	}
}
