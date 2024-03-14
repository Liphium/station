package pipes

import (
	"crypto/cipher"

	"github.com/Liphium/station/pipes/util"
)

func (local *LocalNode) Encrypt(node string, msg []byte) ([]byte, error) {

	// Get key and encrypt message
	var key cipher.Block
	if node == local.ID {
		key = local.Cipher
	} else {
		key = local.GetNode(node).Cipher
	}
	encrypted, err := util.EncryptAES(key, msg)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

func (local *LocalNode) Decrypt(node string, msg []byte) ([]byte, error) {

	// Get key and decrypt message
	var key cipher.Block
	if node == local.ID {
		key = local.Cipher
	} else {
		key = local.GetNode(node).Cipher
	}
	decrypted, err := util.DecryptAES(key, msg)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
