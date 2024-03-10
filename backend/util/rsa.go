package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

const StandardKeySize = 2048

// Generate a new RSA key pair.
func GenerateRSAKey(keySize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, nil, err
	}
	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// Sign a message with a private key.
func SignRSA(privateKey *rsa.PrivateKey, message string) (string, error) {
	hashed := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Verify a signature with a public key. (valid signature = nil error)
func VerifyRSASignature(signature string, publicKey *rsa.PublicKey, message string) error {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(message))
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], sig)
}

// Encrypt a message with a public key. (can't be infinitely long)
func EncryptRSA(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, message)
}

// Decrypt a message with a private key. (can't be infinitely long)
func DecryptRSA(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, ciphertext)
}

// Unpackage public key (in our own format that is kinda crappy but who cares)
func UnpackageRSAPublicKey(pub string) (*rsa.PublicKey, error) {
	parts := strings.Split(pub, ":")
	if len(parts) != 2 {
		return nil, errors.New("invalid public key format")
	}

	modulus, valid := new(big.Int).SetString(parts[0], 36)
	if !valid {
		return nil, errors.New("couldn't parse modulus")
	}

	exponent, valid := new(big.Int).SetString(parts[1], 36)
	if !valid {
		return nil, errors.New("couldn't parse exponent")
	}

	return &rsa.PublicKey{
		N: modulus,
		E: int(exponent.Int64()),
	}, nil
}

// PackageRSAPublicKey packages a public key into our own format.
// Packaging order: modulus, public exponent
func PackageRSAPublicKey(publicKey *rsa.PublicKey) string {
	modulus := publicKey.N.Text(36)
	exponent := big.NewInt(int64(publicKey.E)).Text(36)
	return modulus + ":" + exponent
}

// Unpackage private key (in our own format that is kinda crappy but who cares)
func UnpackageRSAPrivateKey(priv string) (*rsa.PrivateKey, error) {
	parts := strings.Split(priv, ":")
	if len(parts) != 5 {
		return nil, errors.New("invalid private key format")
	}

	n, valid := new(big.Int).SetString(parts[0], 36)
	if !valid {
		return nil, errors.New("couldn't parse modulus")
	}

	e, valid := new(big.Int).SetString(parts[1], 36)
	if !valid {
		return nil, errors.New("couldn't parse public exponent")
	}

	d, valid := new(big.Int).SetString(parts[2], 36)
	if !valid {
		return nil, errors.New("couldn't parse private exponent")
	}

	p, valid := new(big.Int).SetString(parts[3], 36)
	if !valid {
		return nil, errors.New("couldn't parse prime p")
	}

	q, valid := new(big.Int).SetString(parts[4], 36)
	if !valid {
		return nil, errors.New("couldn't parse prime q")
	}

	return &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		},
		D:      d,
		Primes: []*big.Int{p, q},
	}, nil
}

// PackageRSAPrivateKey packages a private key into our own format.
// Packaging order: modulus, public exponent, private exponent, p, q
func PackageRSAPrivateKey(privateKey *rsa.PrivateKey) string {
	d := privateKey.D.Text(36)
	e := big.NewInt(int64(privateKey.E)).Text(36)
	n := privateKey.N.Text(36)
	p := privateKey.Primes[0].Text(36)
	q := privateKey.Primes[1].Text(36)

	return n + ":" + e + ":" + d + ":" + p + ":" + q
}
