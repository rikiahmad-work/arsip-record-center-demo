package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/dto"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"github.com/kemenpppa/siarsip-api/internal/model"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

// MfaService mendefinisikan interface logika bisnis Multi-Factor Authentication.
type MfaService interface {
	StartChallenge(ctx context.Context, user *model.User, tokenLifetimeHours int, ip, userAgent string) (*dto.MfaChallengeData, error)
	Verify(ctx context.Context, req dto.MfaVerifyRequest, ip, userAgent string) (*dto.MfaVerifyData, error)
	GetStatus(ctx context.Context, userID uint) (*dto.MfaStatusResponse, error)
	ResetRecoveryCodes(ctx context.Context, userID uint, code, ip, userAgent string) ([]string, error)
	ResetSelf(ctx context.Context, user *model.User, code, ip, userAgent string) (*dto.MfaChallengeData, error)
	ResetUserByAdmin(ctx context.Context, adminUser *model.User, targetUserID uint, code, ip, userAgent string) error
	IsUserMfaEnabled(ctx context.Context, userID uint) bool
}

type mfaServiceImpl struct {
	cfg            *config.Config
	mfaRepo        repository.MfaRepository
	userRepo       repository.UserRepository
	authMiddleware *middleware.AuthMiddleware
}

// NewMfaService membuat instance implementasi MfaService.
func NewMfaService(
	cfg *config.Config,
	mfaRepo repository.MfaRepository,
	userRepo repository.UserRepository,
	authMiddleware *middleware.AuthMiddleware,
) MfaService {
	return &mfaServiceImpl{
		cfg:            cfg,
		mfaRepo:        mfaRepo,
		userRepo:       userRepo,
		authMiddleware: authMiddleware,
	}
}

func (s *mfaServiceImpl) IsUserMfaEnabled(ctx context.Context, userID uint) bool {
	return s.mfaRepo.IsUserMfaEnabled(ctx, userID)
}

func (s *mfaServiceImpl) StartChallenge(
	ctx context.Context,
	user *model.User,
	tokenLifetimeHours int,
	ip, userAgent string,
) (*dto.MfaChallengeData, error) {
	if tokenLifetimeHours <= 0 {
		tokenLifetimeHours = s.cfg.JWT.AccessExpireHours
	}

	mfa, err := s.mfaRepo.FindUserMfa(ctx, user.ID)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal memeriksa status MFA", err)
	}

	if mfa == nil {
		mfa = &model.UserMfa{
			UserID: user.ID,
		}
	}

	needsSetup := mfa.EnabledAt == nil || mfa.Secret == ""
	var plainSecret string

	if needsSetup {
		plainSecret, err = helper.GenerateTotpSecret()
		if err != nil {
			return nil, helper.NewAppError(500, "Gagal membuat kunci rahasia TOTP", err)
		}

		encryptedPending, err := helper.EncryptSecret(plainSecret, s.cfg.MFA.EncryptionKey)
		if err != nil {
			return nil, helper.NewAppError(500, "Gagal mengamankan kunci rahasia TOTP", err)
		}

		mfa.PendingSecret = encryptedPending
		if err := s.mfaRepo.SaveUserMfa(ctx, mfa); err != nil {
			return nil, helper.NewAppError(500, "Gagal menyimpan inisiasi MFA", err)
		}
	}

	plainChallenge, challengeHash, err := helper.GenerateChallengeToken()
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat tantangan MFA", err)
	}

	purpose := "login"
	mode := "verify"
	if needsSetup {
		purpose = "setup"
		mode = "setup"
	}

	challenge := &model.MfaChallenge{
		UserID:         user.ID,
		ChallengeHash:  challengeHash,
		Purpose:        purpose,
		Attempts:       0,
		IPAddress:      ip,
		UserAgentHash:  helper.HashSHA256(userAgent),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
		TokenExpiresAt: time.Now().Add(time.Duration(tokenLifetimeHours) * time.Hour),
	}

	if err := s.mfaRepo.CreateChallenge(ctx, challenge); err != nil {
		return nil, helper.NewAppError(500, "Gagal mencatat tantangan MFA", err)
	}

	resp := &dto.MfaChallengeData{
		MfaRequired: true,
		Challenge:   plainChallenge,
		Mode:        mode,
	}

	if needsSetup {
		accountIdentifier := user.NIP
		if accountIdentifier == "" {
			accountIdentifier = user.Username
		}
		resp.Secret = plainSecret
		resp.OtpauthURI = helper.ProvisioningURI(s.cfg.MFA.Issuer, accountIdentifier, plainSecret)
	}

	return resp, nil
}

func (s *mfaServiceImpl) Verify(
	ctx context.Context,
	req dto.MfaVerifyRequest,
	ip, userAgent string,
) (*dto.MfaVerifyData, error) {
	challengeHash := helper.HashSHA256(req.Challenge)
	challenge, err := s.mfaRepo.FindChallengeByHash(ctx, challengeHash)
	if err != nil || challenge == nil || challenge.ConsumedAt != nil || challenge.ExpiresAt.Before(time.Now()) {
		return nil, helper.NewAppError(422, "Challenge tidak valid atau sudah kedaluwarsa.", nil)
	}

	// Client binding check (SMP-32)
	currentUAHash := helper.HashSHA256(userAgent)
	if challenge.UserAgentHash != "" && challenge.UserAgentHash != currentUAHash {
		return nil, helper.NewAppError(422, "Sesi challenge tidak cocok dengan klien asal.", nil)
	}

	mfa, err := s.mfaRepo.FindUserMfa(ctx, challenge.UserID)
	if err != nil || mfa == nil {
		return nil, helper.NewAppError(422, "MFA belum tersedia untuk akun ini.", nil)
	}

	// Cek penguncian sementara akun (SMP-31)
	if mfa.LockedUntil != nil && mfa.LockedUntil.After(time.Now()) {
		return nil, helper.NewAppError(429, "Akun dikunci sementara setelah terlalu banyak percobaan MFA. Silakan tunggu 5 menit.", nil)
	}

	// Tentukan secret yang akan diverifikasi
	encryptedSecret := mfa.Secret
	if challenge.Purpose == "setup" {
		encryptedSecret = mfa.PendingSecret
	}
	if encryptedSecret == "" {
		return nil, helper.NewAppError(422, "MFA belum siap diverifikasi.", nil)
	}

	decryptedSecret, err := helper.DecryptSecret(encryptedSecret, s.cfg.MFA.EncryptionKey)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal memproses data keamanan MFA", err)
	}

	// 1. Verifikasi TOTP
	counter, totpOk := helper.VerifyTotp(decryptedSecret, req.Code, mfa.LastTotpCounter, 1)

	// 2. Verifikasi Recovery Code (hanya untuk mode login)
	recoveryIndex := -1
	recoveryOk := false
	var recoveryHashes []string
	if !totpOk && challenge.Purpose != "setup" && mfa.RecoveryCodes != nil && *mfa.RecoveryCodes != "" {
		if err := json.Unmarshal([]byte(*mfa.RecoveryCodes), &recoveryHashes); err == nil {
			recoveryIndex, recoveryOk = helper.VerifyRecoveryCode(req.Code, recoveryHashes)
		}
	}

	// Jika keduanya gagal
	if !totpOk && !recoveryOk {
		mfa.FailedAttempts++
		locked := false
		if mfa.FailedAttempts >= 5 {
			lockTime := time.Now().Add(5 * time.Minute)
			mfa.LockedUntil = &lockTime
			mfa.FailedAttempts = 0
			locked = true
		}
		_ = s.mfaRepo.SaveUserMfa(ctx, mfa)
		_ = s.mfaRepo.IncrementChallengeAttempts(ctx, challenge.ID)

		if challenge.Attempts+1 >= 5 {
			_ = s.mfaRepo.ConsumeChallenge(ctx, challenge.ID)
		}

		_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
			TargetUserID: challenge.UserID,
			Action:       "mfa_verify_failed",
			IPAddress:    ip,
			UserAgent:    userAgent,
		})

		if locked {
			return nil, helper.NewAppError(429, "Akun dikunci sementara setelah terlalu banyak percobaan MFA.", nil)
		}
		return nil, helper.NewAppError(422, "Kode verifikasi salah.", nil)
	}

	// Berhasil verifikasi: bersihkan failed attempts & penguncian
	mfa.FailedAttempts = 0
	mfa.LockedUntil = nil

	var createdRecoveryCodes []string
	if challenge.Purpose == "setup" {
		mfa.Secret = mfa.PendingSecret
		mfa.PendingSecret = ""
		now := time.Now()
		mfa.EnabledAt = &now

		// Buat 10 recovery codes baru
		codes, hashes, err := helper.CreateRecoveryCodes(10)
		if err == nil {
			hashesJSON, _ := json.Marshal(hashes)
			hashesStr := string(hashesJSON)
			mfa.RecoveryCodes = &hashesStr
			createdRecoveryCodes = codes
		}

		_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
			TargetUserID: challenge.UserID,
			Action:       "mfa_setup_completed",
			IPAddress:    ip,
			UserAgent:    userAgent,
		})
	} else if recoveryOk && recoveryIndex >= 0 {
		// Hapus recovery code yang sudah dipakai (single-use)
		if recoveryIndex < len(recoveryHashes) {
			recoveryHashes = append(recoveryHashes[:recoveryIndex], recoveryHashes[recoveryIndex+1:]...)
			hashesJSON, _ := json.Marshal(recoveryHashes)
			hashesStr := string(hashesJSON)
			mfa.RecoveryCodes = &hashesStr
		}

		_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
			TargetUserID: challenge.UserID,
			Action:       "recovery_code_used",
			IPAddress:    ip,
			UserAgent:    userAgent,
		})
	}

	if totpOk && counter > 0 {
		mfa.LastTotpCounter = counter
	}

	if err := s.mfaRepo.SaveUserMfa(ctx, mfa); err != nil {
		return nil, helper.NewAppError(500, "Gagal memperbarui status MFA", err)
	}

	_ = s.mfaRepo.ConsumeChallenge(ctx, challenge.ID)

	// Ambil data user
	user, err := s.userRepo.FindByID(ctx, challenge.UserID)
	if err != nil || user == nil {
		return nil, helper.NewAppError(404, "Data pengguna tidak ditemukan", nil)
	}

	// Pastikan user masih aktif
	if !user.IsActive {
		return nil, helper.ErrAccountInactive
	}

	// Update waktu login terakhir
	now := time.Now()
	user.LastLoginAt = &now
	_ = s.userRepo.Update(ctx, user)

	// Generate JWT Access Token dengan claim mfa_verified = true
	rememberMe := challenge.TokenExpiresAt.Sub(time.Now()) > 10*time.Hour
	token, expTime, err := s.authMiddleware.GenerateTokenWithMFA(
		user.ID,
		user.Username,
		user.NIP,
		user.Role,
		rememberMe,
		true,
	)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat token akses", err)
	}

	permissions, _ := s.userRepo.GetPermissionsByRole(ctx, user.Role)

	_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
		TargetUserID: user.ID,
		Action:       "mfa_verify_success",
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	userResp := &dto.UserResponse{
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
		MfaEnabled:  true,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	expiresIn := int64(time.Until(expTime).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	return &dto.MfaVerifyData{
		AccessToken:   token,
		Token:         token,
		TokenType:     "bearer",
		ExpiresIn:     expiresIn,
		ExpiresAt:     expTime,
		RecoveryCodes: createdRecoveryCodes,
		User:          userResp,
	}, nil
}

func (s *mfaServiceImpl) GetStatus(ctx context.Context, userID uint) (*dto.MfaStatusResponse, error) {
	mfa, err := s.mfaRepo.FindUserMfa(ctx, userID)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal memeriksa status MFA", err)
	}

	if mfa == nil || mfa.EnabledAt == nil || mfa.Secret == "" {
		return &dto.MfaStatusResponse{
			Enabled:                false,
			RecoveryCodesRemaining: 0,
		}, nil
	}

	remaining := 0
	if mfa.RecoveryCodes != nil && *mfa.RecoveryCodes != "" {
		var hashes []string
		if err := json.Unmarshal([]byte(*mfa.RecoveryCodes), &hashes); err == nil {
			remaining = len(hashes)
		}
	}

	return &dto.MfaStatusResponse{
		Enabled:                true,
		EnabledAt:              mfa.EnabledAt,
		RecoveryCodesRemaining: remaining,
	}, nil
}

// verifyCurrentFactor memverifikasi kode OTP aktif atau salah satu recovery code milik pengguna.
func (s *mfaServiceImpl) verifyCurrentFactor(mfa *model.UserMfa, code string) bool {
	if mfa == nil || mfa.EnabledAt == nil || mfa.Secret == "" {
		return false
	}
	if mfa.LockedUntil != nil && mfa.LockedUntil.After(time.Now()) {
		return false
	}

	decryptedSecret, err := helper.DecryptSecret(mfa.Secret, s.cfg.MFA.EncryptionKey)
	if err != nil {
		return false
	}

	counter, totpOk := helper.VerifyTotp(decryptedSecret, code, mfa.LastTotpCounter, 1)
	if totpOk {
		mfa.LastTotpCounter = counter
		mfa.FailedAttempts = 0
		mfa.LockedUntil = nil
		return true
	}

	if mfa.RecoveryCodes != nil && *mfa.RecoveryCodes != "" {
		var hashes []string
		if err := json.Unmarshal([]byte(*mfa.RecoveryCodes), &hashes); err == nil {
			idx, recOk := helper.VerifyRecoveryCode(code, hashes)
			if recOk && idx >= 0 {
				mfa.FailedAttempts = 0
				mfa.LockedUntil = nil
				return true
			}
		}
	}

	return false
}

func (s *mfaServiceImpl) ResetRecoveryCodes(
	ctx context.Context,
	userID uint,
	code, ip, userAgent string,
) ([]string, error) {
	mfa, err := s.mfaRepo.FindUserMfa(ctx, userID)
	if err != nil || mfa == nil {
		return nil, helper.NewAppError(404, "Konfigurasi MFA tidak ditemukan", nil)
	}

	if !s.verifyCurrentFactor(mfa, code) {
		mfa.FailedAttempts++
		if mfa.FailedAttempts >= 5 {
			lockTime := time.Now().Add(5 * time.Minute)
			mfa.LockedUntil = &lockTime
			mfa.FailedAttempts = 0
			_ = s.mfaRepo.SaveUserMfa(ctx, mfa)
			return nil, helper.NewAppError(429, "Terlalu banyak percobaan. Akun dikunci sementara.", nil)
		}
		_ = s.mfaRepo.SaveUserMfa(ctx, mfa)
		return nil, helper.NewAppError(422, "Kode MFA atau recovery code salah.", nil)
	}

	codes, hashes, err := helper.CreateRecoveryCodes(10)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat recovery code", err)
	}

	hashesJSON, _ := json.Marshal(hashes)
	hashesStr := string(hashesJSON)
	mfa.RecoveryCodes = &hashesStr
	if err := s.mfaRepo.SaveUserMfa(ctx, mfa); err != nil {
		return nil, helper.NewAppError(500, "Gagal menyimpan recovery code baru", err)
	}

	actorID := userID
	_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
		ActorUserID:  &actorID,
		TargetUserID: userID,
		Action:       "recovery_codes_reset",
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	return codes, nil
}

func (s *mfaServiceImpl) ResetSelf(
	ctx context.Context,
	user *model.User,
	code, ip, userAgent string,
) (*dto.MfaChallengeData, error) {
	mfa, err := s.mfaRepo.FindUserMfa(ctx, user.ID)
	if err != nil || mfa == nil {
		return nil, helper.NewAppError(404, "Konfigurasi MFA tidak ditemukan", nil)
	}

	if !s.verifyCurrentFactor(mfa, code) {
		mfa.FailedAttempts++
		if mfa.FailedAttempts >= 5 {
			lockTime := time.Now().Add(5 * time.Minute)
			mfa.LockedUntil = &lockTime
			mfa.FailedAttempts = 0
			_ = s.mfaRepo.SaveUserMfa(ctx, mfa)
			return nil, helper.NewAppError(429, "Terlalu banyak percobaan. Akun dikunci sementara.", nil)
		}
		_ = s.mfaRepo.SaveUserMfa(ctx, mfa)
		return nil, helper.NewAppError(422, "Kode MFA atau recovery code salah.", nil)
	}

	// Reset data MFA
	plainSecret, err := helper.GenerateTotpSecret()
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat kunci rahasia baru", err)
	}

	encryptedPending, err := helper.EncryptSecret(plainSecret, s.cfg.MFA.EncryptionKey)
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal mengamankan kunci rahasia baru", err)
	}

	mfa.Secret = ""
	mfa.PendingSecret = encryptedPending
	mfa.RecoveryCodes = nil
	mfa.LastTotpCounter = 0
	mfa.EnabledAt = nil
	mfa.FailedAttempts = 0
	mfa.LockedUntil = nil

	if err := s.mfaRepo.SaveUserMfa(ctx, mfa); err != nil {
		return nil, helper.NewAppError(500, "Gagal mereset status MFA", err)
	}

	// Buat challenge setup baru
	plainChallenge, challengeHash, err := helper.GenerateChallengeToken()
	if err != nil {
		return nil, helper.NewAppError(500, "Gagal membuat tantangan setup", err)
	}

	challenge := &model.MfaChallenge{
		UserID:         user.ID,
		ChallengeHash:  challengeHash,
		Purpose:        "setup",
		Attempts:       0,
		IPAddress:      ip,
		UserAgentHash:  helper.HashSHA256(userAgent),
		ExpiresAt:      time.Now().Add(5 * time.Minute),
		TokenExpiresAt: time.Now().Add(time.Duration(s.cfg.JWT.AccessExpireHours) * time.Hour),
	}
	_ = s.mfaRepo.CreateChallenge(ctx, challenge)

	actorID := user.ID
	_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
		ActorUserID:  &actorID,
		TargetUserID: user.ID,
		Action:       "mfa_reset_self",
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	accountIdentifier := user.NIP
	if accountIdentifier == "" {
		accountIdentifier = user.Username
	}

	return &dto.MfaChallengeData{
		MfaRequired: true,
		Challenge:   plainChallenge,
		Mode:        "setup",
		Secret:      plainSecret,
		OtpauthURI:  helper.ProvisioningURI(s.cfg.MFA.Issuer, accountIdentifier, plainSecret),
	}, nil
}

func (s *mfaServiceImpl) ResetUserByAdmin(
	ctx context.Context,
	adminUser *model.User,
	targetUserID uint,
	code, ip, userAgent string,
) error {
	// Verifikasi hak akses admin
	if adminUser.Role != string(model.RoleAdministrator) {
		return helper.NewAppError(403, "Akses ditolak. Hanya Administrator yang dapat mereset MFA pengguna.", nil)
	}

	// Step-Up Authentication: verifikasi kode OTP milik admin pemanggil
	adminMfa, err := s.mfaRepo.FindUserMfa(ctx, adminUser.ID)
	if err != nil || adminMfa == nil || !s.verifyCurrentFactor(adminMfa, code) {
		return helper.NewAppError(422, "Kode verifikasi Administrator tidak valid.", nil)
	}

	targetUser, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil || targetUser == nil {
		return helper.NewAppError(404, "Pengguna target tidak ditemukan", nil)
	}

	// Eksekusi reset MFA target user
	if err := s.mfaRepo.ResetUserMfa(ctx, targetUserID); err != nil {
		return helper.NewAppError(500, "Gagal mereset MFA pengguna", err)
	}

	actorID := adminUser.ID
	_ = s.mfaRepo.CreateAuditLog(ctx, &model.MfaAuditLog{
		ActorUserID:  &actorID,
		TargetUserID: targetUserID,
		Action:       "mfa_reset_admin",
		IPAddress:    ip,
		UserAgent:    userAgent,
	})

	return nil
}
