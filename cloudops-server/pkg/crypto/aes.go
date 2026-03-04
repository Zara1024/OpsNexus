package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// AESGCM 提供 AES-256-GCM 加解密能力。
type AESGCM struct {
	key []byte
}

// NewAESGCM 创建 AES-256-GCM 工具。
func NewAESGCM(secret string) *AESGCM {
	sum := sha256.Sum256([]byte(secret))
	return &AESGCM{key: sum[:]}
}

// Encrypt 加密明文并返回 Base64 字符串。
func (a *AESGCM) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt 解密 Base64 密文。
func (a *AESGCM) Decrypt(cipherBase64 string) (string, error) {
	if cipherBase64 == "" {
		return "", nil
	}
	buf, err := base64.StdEncoding.DecodeString(cipherBase64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(buf) < nonceSize {
		return "", fmt.Errorf("密文长度非法")
	}
	nonce, cipherText := buf[:nonceSize], buf[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
