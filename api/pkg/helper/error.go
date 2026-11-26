package helper

import (
	"errors"
	"fmt"
)

// AppError adalah custom error untuk konsistensi response API.
type AppError struct {
	Code     int               // HTTP status code
	Message  string            // User-facing message
	Internal error             // Original error (untuk logging)
	Details  map[string]string // Field-level errors
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

// Predefined errors yang bisa digunakan di seluruh aplikasi.
var (
	ErrNotFound          = &AppError{Code: 404, Message: "Data tidak ditemukan"}
	ErrUnauthorized      = &AppError{Code: 401, Message: "Token tidak valid atau sesi telah berakhir"}
	ErrTokenBlacklisted  = &AppError{Code: 401, Message: "Sesi telah berakhir, silakan login kembali"}
	ErrForbidden         = &AppError{Code: 403, Message: "Anda tidak memiliki akses ke sumber daya ini"}
	ErrDuplicateUsername = &AppError{Code: 409, Message: "Username sudah terdaftar"}
	ErrDuplicateNIP      = &AppError{Code: 409, Message: "NIP sudah terdaftar"}
	ErrDuplicateEmail    = &AppError{Code: 409, Message: "Email sudah terdaftar"}
	ErrSelfDelete        = &AppError{Code: 400, Message: "Tidak dapat menghapus akun Anda sendiri"}
	ErrSelfDeactivate    = &AppError{Code: 400, Message: "Tidak dapat menonaktifkan akun Anda sendiri"}
	ErrAccountInactive   = &AppError{Code: 403, Message: "Akun Anda telah dinonaktifkan. Silakan hubungi administrator."}
)

// NewValidationError membuat AppError untuk validasi gagal (422).
func NewValidationError(details map[string]string) *AppError {
	return &AppError{
		Code:    422,
		Message: "Validasi gagal",
		Details: details,
	}
}

// NewAppError membuat AppError dengan kode dan pesan kustom.
func NewAppError(code int, message string, internal error) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Internal: internal,
	}
}

// IsAppError memeriksa apakah error adalah *AppError.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError mengekstrak *AppError dari error chain.
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
