package helper

import (
	"testing"
	"time"
)

func TestGenerateTotpSecret(t *testing.T) {
	secret, err := GenerateTotpSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secret) != 32 {
		t.Fatalf("expected 32 characters Base32 secret, got %d", len(secret))
	}
}

func TestCalculateAndVerifyTotp(t *testing.T) {
	secret, err := GenerateTotpSecret()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	step := uint64(time.Now().Unix() / 30)
	code, err := CalculateTotp(secret, step)
	if err != nil {
		t.Fatalf("failed to calculate TOTP: %v", err)
	}

	if len(code) != 6 {
		t.Fatalf("expected 6 digit OTP, got %s", code)
	}

	// Verifikasi berhasil
	usedStep, ok := VerifyTotp(secret, code, 0, 1)
	if !ok || usedStep != step {
		t.Fatalf("expected verification to succeed at step %d, got ok=%v, step=%d", step, ok, usedStep)
	}

	// Anti-replay: counter tidak boleh <= lastUsedCounter
	_, replayOk := VerifyTotp(secret, code, step, 1)
	if replayOk {
		t.Fatalf("expected replay attack with same counter to be rejected")
	}

	// Kode salah harus ditolak
	_, wrongOk := VerifyTotp(secret, "000000", 0, 1)
	if wrongOk && code != "000000" {
		t.Fatalf("expected wrong code to be rejected")
	}

	// Format salah harus ditolak
	_, invalidOk := VerifyTotp(secret, "abc", 0, 1)
	if invalidOk {
		t.Fatalf("expected non-digit code to be rejected")
	}
}

func TestEncryptDecryptSecret(t *testing.T) {
	key := "my-very-secure-32-byte-aes-key!!"
	secret := "JBSWY3DPEHPK3PXP"

	ciphertext, err := EncryptSecret(secret, key)
	if err != nil {
		t.Fatalf("failed to encrypt secret: %v", err)
	}
	if ciphertext == secret || len(ciphertext) == 0 {
		t.Fatalf("expected encrypted ciphertext, got %s", ciphertext)
	}

	decrypted, err := DecryptSecret(ciphertext, key)
	if err != nil {
		t.Fatalf("failed to decrypt secret: %v", err)
	}
	if decrypted != secret {
		t.Fatalf("expected decrypted '%s', got '%s'", secret, decrypted)
	}

	// Kunci salah harus gagal
	_, errWrongKey := DecryptSecret(ciphertext, "wrong-key")
	if errWrongKey == nil {
		t.Fatalf("expected error when decrypting with wrong key")
	}
}

func TestRecoveryCodes(t *testing.T) {
	codes, hashes, err := CreateRecoveryCodes(10)
	if err != nil {
		t.Fatalf("failed to create recovery codes: %v", err)
	}
	if len(codes) != 10 || len(hashes) != 10 {
		t.Fatalf("expected 10 codes and 10 hashes, got %d and %d", len(codes), len(hashes))
	}

	// Verifikasi kode pertama
	idx, ok := VerifyRecoveryCode(codes[0], hashes)
	if !ok || idx != 0 {
		t.Fatalf("expected code[0] to match index 0, got ok=%v, idx=%d", ok, idx)
	}

	// Verifikasi kode ketiga
	idx3, ok3 := VerifyRecoveryCode(codes[3], hashes)
	if !ok3 || idx3 != 3 {
		t.Fatalf("expected code[3] to match index 3, got ok=%v, idx=%d", ok3, idx3)
	}

	// Kode salah
	_, wrongOk := VerifyRecoveryCode("AAAAA-BBBBB-CCCCC-DDDDD", hashes)
	if wrongOk {
		t.Fatalf("expected invalid recovery code to be rejected")
	}
}

func TestGenerateChallengeToken(t *testing.T) {
	plain, hash, err := GenerateChallengeToken()
	if err != nil {
		t.Fatalf("failed to generate challenge: %v", err)
	}
	if len(plain) != 64 {
		t.Fatalf("expected 64 char hex plain challenge, got %d", len(plain))
	}
	if len(hash) != 64 {
		t.Fatalf("expected 64 char hex challenge hash, got %d", len(hash))
	}
	if HashSHA256(plain) != hash {
		t.Fatalf("expected hash to match SHA256 of plain")
	}
}
