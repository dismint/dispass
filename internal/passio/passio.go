package passio

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/dismint/dispass/internal/state"
	"github.com/dismint/dispass/internal/uconst"
)

func SecretFromString(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

func Encrypt(key, plaintext []byte) ([]byte, error) {
	// create aes block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// generate a random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// encrypt and prepend nonce to ciphertext
	ciphertext := aesgcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func Decrypt(key, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, errors.New("incorrect password or corrupted data")
	}

	return plaintext, nil
}

func WriteStateCreds(sm *state.Model) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(sm.KeyToCredInfo); err != nil {
		log.Fatalf("failed to encode: %v", err)
	}

	bytes, err := Encrypt(sm.Secret, buf.Bytes())
	if err != nil {
		log.Fatalf("failed to encrypt: %v", err)
	}

	if err := os.WriteFile(uconst.DataFileName, bytes, 0644); err != nil {
		log.Fatalf("failed to write to %v: %v", uconst.DataFileName, err)
	}
}

func ReadStateCreds(sm *state.Model) error {
	dat, err := os.ReadFile(uconst.DataFileName)
	if err != nil {
		log.Fatalf("failed to write to %v: %v", uconst.DataFileName, err)
	}

	if len(dat) > 0 {
		d, err := Decrypt(sm.Secret, dat)
		if err != nil {
			log.Warnf("failed to decrypt: %v", err)
			return err
		}
		buf := bytes.NewBuffer(d)
		dec := gob.NewDecoder(buf)
		if err = dec.Decode(&sm.KeyToCredInfo); err != nil {
			log.Fatalf("failed to decode: %v", err)
		}
	}

	return nil
}
