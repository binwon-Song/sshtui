package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var secretKey []byte

// generateRandomKey generates a random 32-byte key
func generateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// getKeyFilePath returns the path to the secret key file
func getKeyFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".sshtui")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "secret.key"), nil
}

// loadOrCreateKey loads the existing key or creates a new one
func loadOrCreateKey() error {
	keyPath, err := getKeyFilePath()
	if err != nil {
		return err
	}

	// Try to load existing key
	keyBytes, err := os.ReadFile(keyPath)
	if err == nil && len(keyBytes) == 64 { // 32 bytes in hex = 64 characters
		// Decode hex string to bytes
		secretKey = make([]byte, 32)
		_, err = hex.Decode(secretKey, keyBytes)
		if err == nil {
			return nil
		}
	}

	// Generate new key if loading failed
	newKey, err := generateRandomKey()
	if err != nil {
		return err
	}

	// Convert key to hex string for storage
	keyHex := make([]byte, hex.EncodedLen(len(newKey)))
	hex.Encode(keyHex, newKey)

	// Save new key
	if err := os.WriteFile(keyPath, keyHex, 0600); err != nil {
		return err
	}

	secretKey = newKey
	return nil
}

// init initializes the secret key when the package is loaded
func init() {
	if err := loadOrCreateKey(); err != nil {
		panic("Failed to initialize encryption key: " + err.Error())
	}
}

// Encrypt 함수는 평문 비밀번호를 암호화합니다.
func Encrypt(password string) (string, error) {
	if len(secretKey) == 0 {
		return "", errors.New("encryption key is not initialized")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(password)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 함수는 암호화된 비밀번호를 복호화합니다.
func Decrypt(encrypted string) (string, error) {
	if len(secretKey) == 0 {
		return "", errors.New("encryption key is not initialized")
	}

	if encrypted == "" {
		return "", errors.New("암호화된 비밀번호가 없습니다")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.URLEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("암호화된 데이터가 너무 짧습니다")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}
