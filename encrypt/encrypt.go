package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ozzadar/platformer_mission_server/models"
)

func DecryptMetaFile(ciphertext []byte) *models.Meta {
	hasher := md5.New()
	hasher.Write([]byte("thing"))

	key := hex.EncodeToString(hasher.Sum(nil))

	plaintext := decryptAes128Ecb(ciphertext, []byte(key))

	plaintext = bytes.Trim(plaintext, "\x00")
	plaintext = bytes.Trim(plaintext, "\x01")

	meta := &models.Meta{}

	err := json.Unmarshal(plaintext, meta)

	if err != nil {
		fmt.Printf("Failed to unmarshal json: %v", err)
		return nil
	}

	return meta
}

func decryptAes128Ecb(data, key []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(key))
	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		cipher.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted[32:]
}
