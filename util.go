package gomail

import (
	"bytes"
	"fmt"
	"net/mail"

	"github.com/digisan/gotk/crypto"
)

func genCode(str string, key []byte) string {
	if len(key) > 16 {
		key = key[:16]
	}
	if len(key) < 16 {
		zero := bytes.Repeat([]byte{0}, 16)
		copy(zero, key)
	}
	return fmt.Sprintf("%x", crypto.Encrypt(str, key))
}

func translateKey(code string, key []byte) string {
	if len(key) > 16 {
		key = key[:16]
	}
	if len(key) < 16 {
		zero := bytes.Repeat([]byte{0}, 16)
		copy(zero, key)
	}
	data := []byte{}
	fmt.Sscanf(code, "%x", &data)
	return crypto.Decrypt(data, key)
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
