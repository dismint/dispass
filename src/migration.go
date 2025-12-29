package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	creds        = make(map[string]string)
	jcpwFileName = "pass.dat"
)

func textToSecret(text string) []byte {
	textBytes := []byte(text)
	pad := "aaaaaaaaaaaaaaaa"
	if len(textBytes) == 16 {
		return textBytes
	} else if len(textBytes) > 16 {
		return textBytes[:16]
	} else {
		return []byte(text + string(pad))[:16]
	}
}

func (m *model) jcpwDecrypt() error {
	secret := textToSecret(m.entryModel.passwordInput.Value())
	fileBytes, err := os.ReadFile(jcpwFileName)
	block, err := aes.NewCipher(secret)
	if err != nil {
		return err
	}
	if len(fileBytes) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	iv := fileBytes[:aes.BlockSize]
	fileBytes = fileBytes[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(fileBytes, fileBytes)
	buf := bytes.NewBuffer(fileBytes)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&creds); err != nil {
		panic("invalid key")
	}

	m.keyToCredInfo = make(map[string]credInfo)
	m.secret = secretFromString(m.entryModel.passwordInput.Value())
	for meta, password := range creds {
		splitMeta := strings.Split(meta, " ")
		newUUID := uuid.NewString()
		if len(splitMeta) == 1 {
			m.keyToCredInfo[newUUID] = credInfo{
				Source:   "",
				Username: splitMeta[0],
				Password: password,
			}
		} else if len(splitMeta) > 1 {
			m.keyToCredInfo[newUUID] = credInfo{
				Source:   strings.Join(splitMeta[:len(splitMeta)-1], " "),
				Username: splitMeta[len(splitMeta)-1],
				Password: password,
			}
		}
	}
	m.writePass()
	m.initFuzzy()
	return nil
}
