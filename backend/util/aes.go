package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
)

// Encrypt encrypts the given plaintext using AES-GCM.
func EncryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(nonce, ciphertext...)

	return ciphertext, nil
}

// Decrypt decrypts the given ciphertext using AES-GCM.
func DecryptAES(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aesGCM.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:aesGCM.NonceSize()]
	ciphertext = ciphertext[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// Generate a new AES key with a length of 32
func NewAESKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func TestAES() {

	key, _ := base64.StdEncoding.DecodeString("EKc1lOIKwZXb/BSiHfilfCd+ptNMDRg7eleShj/VYrE=")
	msg, _ := base64.StdEncoding.DecodeString("4MlmspkWoWemUoKQKHoUfEXfs+bdURGG/ve05n49A62l6ax+dBk6")

	decrypted, err := DecryptAES(key, msg)
	if err != nil {
		log.Fatalln("Error while decrypting with AES:", err)
	}
	log.Println(string(decrypted))

}
