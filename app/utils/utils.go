package utils

import (
	"crypto/rand"
)

var (
	randomStringCharset = `0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
)

func RandString(size uint) string {
	buf := make([]byte, size)
	rand.Read(buf)

	lr := uint(len(randomStringCharset))

	for i := uint(0); i < size; i++ {
		buf[i] = randomStringCharset[buf[i]%byte(lr)]
	}
	return string(buf)
}
