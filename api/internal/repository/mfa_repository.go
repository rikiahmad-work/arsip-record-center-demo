package repository

import (
	"context"
	"errors"
	"time"

	"github.com/kemenpppa/siarsip-api/internal/model"
	"gorm.io/gorm"
)

// MfaRepository mendefinisikan interface interaksi database untuk fitur MFA.
type MfaRepository interface {
	FindUserMfa(ctx context.Context, userID uint) (*model.UserMfa, error)
	SaveUserMfa(ctx context.Context, mfa *model.UserMfa) error
	CreateChallenge(ctx context.Context, challenge *model.MfaChallenge) error
	FindChallengeByHash(ctx context.Context, challengeHash string) (*model.MfaChallenge, error)
	IncrementChallengeAttempts(ctx context.Context, challengeID uint) error
	ConsumeChallenge(ctx context.Context, challengeID uint) error
	DeleteUserChallenges(ctx context.Context, userID uint) error
	CreateAuditLog(ctx context.Context, auditLog *model.MfaAuditLog) error
	ResetUserMfa(ctx context.Context, userID uint) error
	IsUserMfaEnabled(ctx context.Context, userID uint) bool
	GetMfaStatusesByUserIDs(ctx context.Context, userIDs []uint) (map[uint]bool, error)
}

type mfaRepositoryImpl struct {
	db *gorm.DB
}

// NewMfaRepository membuat instance implementasi MfaRepository.
func NewMfaRepository(db *gorm.DB) MfaRepository {
	return &mfaRepositoryImpl{db: db}
}

func (r *mfaRepositoryImpl) FindUserMfa(ctx context.Context, userID uint) (*model.UserMfa, error) {
	var mfa model.UserMfa
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&mfa).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mfa, nil
}

func (r *mfaRepositoryImpl) SaveUserMfa(ctx context.Context, mfa *model.UserMfa) error {
	return r.db.WithContext(ctx).Save(mfa).Error
}

func (r *mfaRepositoryImpl) CreateChallenge(ctx context.Context, challenge *model.MfaChallenge) error {
	// Hapus challenge sebelumnya milik user yang belum terpakai
	_ = r.db.WithContext(ctx).
		Where("user_id = ? AND consumed_at IS NULL", challenge.UserID).
		Delete(&model.MfaChallenge{}).Error

	return r.db.WithContext(ctx).Create(challenge).Error
}

func (r *mfaRepositoryImpl) FindChallengeByHash(ctx context.Context, challengeHash string) (*model.MfaChallenge, error) {
	var challenge model.MfaChallenge
	err := r.db.WithContext(ctx).
		Where("challenge_hash = ?", challengeHash).
		First(&challenge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &challenge, nil
}

func (r *mfaRepositoryImpl) IncrementChallengeAttempts(ctx context.Context, challengeID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.MfaChallenge{}).
		Where("id = ?", challengeID).
		Update("attempts", gorm.Expr("attempts + 1")).Error
}

func (r *mfaRepositoryImpl) ConsumeChallenge(ctx context.Context, challengeID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.MfaChallenge{}).
		Where("id = ?", challengeID).
		Update("consumed_at", now).Error
}

func (r *mfaRepositoryImpl) DeleteUserChallenges(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.MfaChallenge{}).Error
}

func (r *mfaRepositoryImpl) CreateAuditLog(ctx context.Context, auditLog *model.MfaAuditLog) error {
	return r.db.WithContext(ctx).Create(auditLog).Error
}

func (r *mfaRepositoryImpl) ResetUserMfa(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Kosongkan konfigurasi MFA pengguna
		if err := tx.Model(&model.UserMfa{}).Where("user_id = ?", userID).Updates(map[string]interface{}{
			"secret":            nil,
			"pending_secret":    nil,
			"recovery_codes":    nil,
			"last_totp_counter": 0,
			"failed_attempts":   0,
			"locked_until":      nil,
			"enabled_at":        nil,
		}).Error; err != nil {
			return err
		}

		// Hapus semua challenge aktif
		return tx.Where("user_id = ?", userID).Delete(&model.MfaChallenge{}).Error
	})
}

func (r *mfaRepositoryImpl) IsUserMfaEnabled(ctx context.Context, userID uint) bool {
	var count int64
	r.db.WithContext(ctx).
		Model(&model.UserMfa{}).
		Where("user_id = ? AND enabled_at IS NOT NULL AND secret IS NOT NULL AND secret != ''", userID).
		Count(&count)
	return count > 0
}

func (r *mfaRepositoryImpl) GetMfaStatusesByUserIDs(ctx context.Context, userIDs []uint) (map[uint]bool, error) {
	result := make(map[uint]bool)
	if len(userIDs) == 0 {
		return result, nil
	}

	var mfas []model.UserMfa
	err := r.db.WithContext(ctx).
		Where("user_id IN ? AND enabled_at IS NOT NULL AND secret IS NOT NULL AND secret != ''", userIDs).
		Find(&mfas).Error
	if err != nil {
		return nil, err
	}

	for _, m := range mfas {
		result[m.UserID] = true
	}
	return result, nil
}
