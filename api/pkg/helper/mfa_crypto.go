package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	base32Encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
	recoveryRegex  = regexp.MustCompile(`^[A-F0-9]{20}$`)
	otpRegex       = regexp.MustCompile(`^[0-9]{6}$`)
)

// GenerateTotpSecret membuat kunci rahasia TOTP acak 160-bit (20 byte) dalam format Base32 (32 karakter).
func GenerateTotpSecret() (string, error) {
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("gagal menghasilkan random bytes: %w", err)
	}
	return base32Encoding.EncodeToString(bytes), nil
}

// ProvisioningURI menghasilkan URI QR code sesuai standar otpauth:// (RFC 6238).
func ProvisioningURI(issuer, account, secret string) string {
	encodedIssuer := url.QueryEscape(issuer)
	encodedAccount := url.QueryEscape(account)
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=6&period=30",
		encodedIssuer, encodedAccount, secret, encodedIssuer)
}

// CalculateTotp menghitung kode TOTP 6-digit untuk time step tertentu berdasarkan RFC 6238.
func CalculateTotp(secret string, timeStep uint64) (string, error) {
	key, err := base32Encoding.DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
	if err != nil {
		// Fallback dengan padding bila ada
		key, err = base32.StdEncoding.DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
		if err != nil {
			return "", fmt.Errorf("secret Base32 tidak valid: %w", err)
		}
	}

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, timeStep)

	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	h := mac.Sum(nil)

	offset := h[len(h)-1] & 0x0f
	binaryCode := ((int(h[offset]) & 0x7f) << 24) |
		((int(h[offset+1]) & 0xff) << 16) |
		((int(h[offset+2]) & 0xff) << 8) |
		(int(h[offset+3]) & 0xff)

	otp := binaryCode % 1000000
	return fmt.Sprintf("%06d", otp), nil
}

// VerifyTotp memverifikasi kode OTP 6-digit dengan toleransi pergeseran waktu (window) dan anti-replay.
func VerifyTotp(secret, code string, lastUsedCounter uint64, window int) (uint64, bool) {
	trimmedCode := strings.TrimSpace(code)
	if !otpRegex.MatchString(trimmedCode) {
		return 0, false
	}

	nowStep := uint64(time.Now().Unix() / 30)
	start := nowStep
	if start >= uint64(window) {
		start -= uint64(window)
	} else {
		start = 0
	}
	end := nowStep + uint64(window)

	for step := start; step <= end; step++ {
		// Anti-replay: counter harus lebih besar dari counter yang sudah pernah dipakai
		if lastUsedCounter > 0 && step <= lastUsedCounter {
			continue
		}

		expected, err := CalculateTotp(secret, step)
		if err != nil {
			continue
		}

		if subtle.ConstantTimeCompare([]byte(expected), []byte(trimmedCode)) == 1 {
			return step, true
		}
	}

	return 0, false
}

// deriveAESKey menghasilkan kunci AES-256 (32 byte) dari string kunci apa pun via SHA-256.
func deriveAESKey(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

// EncryptSecret mengenkripsi plaintext secret menggunakan AES-256-GCM.
func EncryptSecret(plaintext, key string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(deriveAESKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSecret mendekripsi ciphertext secret yang dienkripsi AES-256-GCM.
func DecryptSecret(ciphertextBase64, key string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(deriveAESKey(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext terlalu pendek")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("gagal dekripsi secret (kunci salah atau data korup): %w", err)
	}

	return string(plaintext), nil
}

// CreateRecoveryCodes menghasilkan 10 kode pemulihan acak berformat XXXXX-XXXXX-XXXXX-XXXXX
// dan hash SHA-256-nya untuk disimpan ke database.
func CreateRecoveryCodes(count int) (codes []string, hashes []string, err error) {
	if count <= 0 {
		count = 10
	}

	codes = make([]string, 0, count)
	hashes = make([]string, 0, count)

	for i := 0; i < count; i++ {
		bytes := make([]byte, 10)
		if _, err := rand.Read(bytes); err != nil {
			return nil, nil, err
		}

		rawHex := strings.ToUpper(hex.EncodeToString(bytes)) // 20 hex chars
		formatted := fmt.Sprintf("%s-%s-%s-%s", rawHex[0:5], rawHex[5:10], rawHex[10:15], rawHex[15:20])

		h := sha256.Sum256([]byte(rawHex))
		hashHex := hex.EncodeToString(h[:])

		codes = append(codes, formatted)
		hashes = append(hashes, hashHex)
	}

	return codes, hashes, nil
}

// NormalizeRecoveryCode membersihkan tanda hubung dan spasi serta memastikan format 20 karakter hex.
func NormalizeRecoveryCode(code string) string {
	cleaned := strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(code), "-", ""), " ", ""))
	if !recoveryRegex.MatchString(cleaned) {
		return ""
	}
	return cleaned
}

// VerifyRecoveryCode mencari apakah kode pemulihan cocok dengan salah satu hash SHA-256 yang tersimpan.
// Mengembalikan index elemen yang cocok (untuk dihapus / single-use) dan true bila valid.
func VerifyRecoveryCode(code string, storedHashes []string) (int, bool) {
	normalized := NormalizeRecoveryCode(code)
	if normalized == "" {
		return -1, false
	}

	h := sha256.Sum256([]byte(normalized))
	targetHashHex := hex.EncodeToString(h[:])

	for i, candidate := range storedHashes {
		if subtle.ConstantTimeCompare([]byte(candidate), []byte(targetHashHex)) == 1 {
			return i, true
		}
	}

	return -1, false
}

// GenerateChallengeToken membuat random token 64 karakter hex dan hash SHA-256-nya.
func GenerateChallengeToken() (plainText string, hashHex string, err error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	plainText = hex.EncodeToString(bytes) // 64 hex characters
	h := sha256.Sum256([]byte(plainText))
	hashHex = hex.EncodeToString(h[:])
	return plainText, hashHex, nil
}

// HashSHA256 mengembalikan representasi hex SHA-256 dari string input.
func HashSHA256(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}
