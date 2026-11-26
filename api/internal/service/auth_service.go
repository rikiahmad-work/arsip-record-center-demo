package service

import (
	"context"
	"errors"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
	"gorm.io/gorm"
)

// AuthService mendefinisikan interface untuk business logic autentikasi.
type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, tokenJTI string) error
	Me(ctx context.Context, userID uint) (*dto.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error
}

type authServiceImpl struct {
	userRepo repository.UserRepository
	jwtSvc   JWTService
	hashHelp *helper.HashHelper
	mfaSvc   MfaService
	cfg      *config.Config
}

// NewAuthService membuat instance AuthService.
func NewAuthService(
	userRepo repository.UserRepository,
	jwtSvc JWTService,
	hashHelp *helper.HashHelper,
	mfaSvc MfaService,
	cfg *config.Config,
) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
		hashHelp: hashHelp,
		mfaSvc:   mfaSvc,
		cfg:      cfg,
	}
}

func (s *authServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	identifier := req.Identifier()
	if identifier == "" {
		return nil, helper.NewAppError(422, "NIP atau Username wajib diisi", nil)
	}

	// Cari user berdasarkan NIP terlebih dahulu, lalu fallback ke Username
	user, err := s.userRepo.FindByNIP(ctx, identifier)
	if err != nil {
		user, err = s.userRepo.FindByUsername(ctx, identifier)
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewAppError(401, "NIP/Username atau password salah", nil)
		}
		return nil, err
	}

	// Cek apakah user aktif
	if !user.IsActive {
		return nil, helper.ErrAccountInactive
	}

	// Verifikasi password
	if !s.hashHelp.Verify(req.Password, user.Password) {
		return nil, helper.NewAppError(401, "Username atau password salah", nil)
	}

	// Jika MFA diwajibkan oleh konfigurasi sistem, inisiasi tantangan MFA (2-step login)
	if s.cfg != nil && s.cfg.MFA.Required && s.mfaSvc != nil {
		tokenHours := s.cfg.JWT.AccessExpireHours
		if req.RememberMe {
			tokenHours = s.cfg.JWT.RefreshExpireHours
		}
		challenge, err := s.mfaSvc.StartChallenge(ctx, user, tokenHours, req.ClientIP, req.UserAgent)
		if err != nil {
			return nil, err
		}
		return &dto.LoginResponse{
			MfaRequired: true,
			Challenge:   challenge.Challenge,
			Mode:        challenge.Mode,
			Secret:      challenge.Secret,
			OtpauthURI:  challenge.OtpauthURI,
		}, nil
	}

	// Generate JWT token (fallback tanpa MFA)
	token, expiresAt, err := s.jwtSvc.GenerateToken(user.ID, user.Username, user.NIP, user.Role, req.RememberMe)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat token", err)
	}

	// Update last_login_at
	// Note: update dilakukan async untuk mengurangi latency response
	go func() {
		s.userRepo.Update(context.Background(), user)
	}()

	// Ambil daftar permissions berdasarkan peran pengguna
	permissions, _ := s.userRepo.GetPermissionsByRole(ctx, user.Role)

	// Buat response
	expTime := time.Unix(expiresAt, 0)
	resp := &dto.LoginResponse{
		User: &dto.UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			Name:        user.Name,
			NIP:         user.NIP,
			Email:       user.Email,
			Role:        user.Role,
			RoleLabel:   model.UserRole(user.Role).Label(),
			Permissions: permissions,
			UnitKerja:   user.UnitKerja,
			Jabatan:     user.Jabatan,
			IsSSO:       user.IsSSO,
			IsActive:    user.IsActive,
			LastLoginAt: user.LastLoginAt,
		},
		Token:     token,
		ExpiresAt: &expTime,
	}

	return resp, nil
}

func (s *authServiceImpl) Logout(ctx context.Context, tokenJTI string) error {
	// Blacklist JWT token di tabel blacklisted_tokens hingga masa berlaku habis (default +30 jam)
	expiresAt := time.Now().Add(30 * time.Hour).Unix()
	return s.jwtSvc.BlacklistToken(tokenJTI, expiresAt)
}

func (s *authServiceImpl) Me(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, err
	}

	// Ambil daftar permissions berdasarkan peran pengguna
	permissions, _ := s.userRepo.GetPermissionsByRole(ctx, user.Role)

	mfaEnabled := false
	if s.mfaSvc != nil {
		mfaEnabled = s.mfaSvc.IsUserMfaEnabled(ctx, user.ID)
	}

	resp := &dto.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Name:        user.Name,
		NIP:         user.NIP,
		Email:       user.Email,
		Role:        user.Role,
		RoleLabel:   model.UserRole(user.Role).Label(),
		Permissions: permissions,
		UnitKerja:   user.UnitKerja,
		Jabatan:     user.Jabatan,
		IsSSO:       user.IsSSO,
		IsActive:    user.IsActive,
		MfaEnabled:  mfaEnabled,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return resp, nil
}

func (s *authServiceImpl) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.ErrNotFound
		}
		return err
	}

	if user.IsSSO {
		return helper.NewAppError(400, "Akun SSO tidak dapat mengubah kata sandi lokal", nil)
	}

	if !s.hashHelp.Verify(req.OldPassword, user.Password) {
		return helper.NewAppError(400, "Kata sandi saat ini tidak sesuai", nil)
	}

	hashed, err := s.hashHelp.Hash(req.NewPassword)
	if err != nil {
		return helper.NewAppError(500, "Gagal memproses kata sandi baru", err)
	}

	user.Password = hashed
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
