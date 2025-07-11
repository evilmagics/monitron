package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"

	"monitron-server/config"
)

// Encrypt encrypts data using AES-256 GCM.
func Encrypt(data []byte, cfg *config.Config) (string, error) {
	key := []byte(cfg.EncryptionKey)
	block, err := aes.NewCipher(key)
		if err != nil {
		return "", fmt.Errorf("could not create new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
		if err != nil {
		return "", fmt.Errorf("could not create new GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("could not read nonce: %w", err)
	}

	sealed := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Decrypt decrypts data using AES-256 GCM.
func Decrypt(encryptedData string, cfg *config.Config) ([]byte, error) {
	key := []byte(cfg.EncryptionKey)
	data, err := base64.StdEncoding.DecodeString(encryptedData)
		if err != nil {
		return nil, fmt.Errorf("could not decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
		if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
		if err != nil {
		return nil, fmt.Errorf("could not create new GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
		if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
		return nil, fmt.Errorf("could not decrypt data: %w", err)
	}

	return plaintext, nil
}

// GenerateRandomKey generates a random AES-256 key.
func GenerateRandomKey() string {
	key := make([]byte, 32) // AES-256 requires a 32-byte key
		if _, err := io.ReadFull(rand.Reader, key); err != nil {
			log.Fatal().Err(err).Msg("Failed to generate random key")
	return base64.StdEncoding.EncodeToString(key)
}

