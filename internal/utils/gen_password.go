package utils

import (
	"math/rand"
	"time"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialBytes = "!@#$%^&*()_+-=[]{}\\|;':\",.<>/?`~"
	numBytes     = "0123456789"
)

func GeneratePassword(length int, useLetters bool, useSpecial bool, useNum bool) string {
	b := make([]byte, length)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		if useLetters {
			b[i] = letterBytes[rnd.Intn(len(letterBytes))]
		} else if useSpecial {
			b[i] = specialBytes[rnd.Intn(len(specialBytes))]
		} else if useNum {
			b[i] = numBytes[rnd.Intn(len(numBytes))]
		}
	}
	return string(b)
}
