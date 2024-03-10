package tests

import (
	"encoding/base64"
	"fmt"
	"log"
	"node-backend/util"
	"testing"
)

func TestEncryption(t *testing.T) {

	// Generate RSA key pair
	priv, pub, err := util.GenerateRSAKey(util.StandardKeySize)
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	log.Println(priv.D.String() + " | " + fmt.Sprintf("%d", priv.PublicKey.E))

	//* Test public key packaging and unpackaging
	packaged := util.PackageRSAPublicKey(pub)
	pub, err = util.UnpackageRSAPublicKey(packaged)
	if err != nil {
		t.Errorf("Couldn't package and unpackage public key: %v", err)
	}

	//* Test private key packaging and unpackaging
	packaged = util.PackageRSAPrivateKey(priv)
	priv, err = util.UnpackageRSAPrivateKey(packaged)
	if err != nil {
		t.Errorf("Couldn't package and unpackage private key: %v", err)
	}

	//* Test encryption
	// Encrypt a sample plaintext using the public key
	plaintext := "abcdefghijklmopqrstuvxyzABC-._ßäöü!"
	ciphertext, err := util.EncryptRSA(pub, []byte(plaintext))
	if err != nil {
		t.Fatalf("Failed to encrypt plaintext: %v", err)
	}

	// Decrypt the ciphertext using the private key
	decryptedPlaintext, err := util.DecryptRSA(priv, ciphertext)
	if err != nil {
		t.Fatalf("Failed to decrypt ciphertext: %v", err)
	}

	// Check if the decrypted plaintext matches the original plaintext
	if base64.StdEncoding.EncodeToString(decryptedPlaintext) != base64.StdEncoding.EncodeToString([]byte(plaintext)) {
		t.Errorf("Decrypted plaintext does not match original plaintext. Expected: %s, Got: %s", plaintext, decryptedPlaintext)
	}

	//* Test signatures
	// Sign the plaintext using the private key
	signature, err := util.SignRSA(priv, plaintext)
	if err != nil {
		t.Fatalf("Failed to sign plaintext: %v", err)
	}

	// Verify the signature using the public key
	err = util.VerifyRSASignature(signature, pub, plaintext)
	if err != nil {
		t.Errorf("Signature verification failed")
	}
}
