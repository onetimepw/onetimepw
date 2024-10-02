package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
)

func generateKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))

	return hash[:]
}

func encrypt(plaintext string, password string) (string, error) {
	blockData, err := aes.NewCipher(generateKey(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockData)
	if err != nil {
		return "", err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

func decrypt(ciphertext string, password string) (string, error) {
	blockData, err := aes.NewCipher(generateKey(password))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockData)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func fromBase64(text string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(text)
}

func randomPassword(passwordLength int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, passwordLength)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}

	return string(b), nil
}
