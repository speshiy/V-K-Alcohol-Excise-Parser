package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
)

const fish = "ffklkfd3o24mf9rksnfSDFSDf421"

func createHash(key string) string {
	hasher := md5.New()
	_, _ = hasher.Write([]byte(key + fish))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

//SetReadOnly set 0444 permission to ifle
func SetReadOnly(filepath string) error {
	err := os.Chmod(filepath, 0444)
	return err
}

func decryptConfigFile(filename string, passphrase string) []byte {
	_ = SetReadOnly(filename)
	data, _ := ioutil.ReadFile(filename)
	return decrypt(data, passphrase)
}
