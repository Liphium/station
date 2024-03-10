package connection

import (
	"github.com/Liphium/station/pipes"
	"github.com/Liphium/station/pipes/util"
)

func Encrypt(node string, msg []byte) ([]byte, error) {

	// Get key and encrypt message
	key := pipes.GetNode(node).Cipher
	encrypted, err := util.EncryptAES(key, msg)
	if err != nil {
		return nil, err
	}

	// Add general prefix
	encrypted = append(GeneralPrefix, encrypted...)

	return encrypted, nil
}

func Decrypt(node string, msg []byte) ([]byte, error) {

	// Get key and decrypt message
	key := pipes.GetNode(node).Cipher
	decrypted, err := util.DecryptAES(key, msg)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
