package secstring

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strings"
	"syscall"
)

type SecString struct {
	String  []byte
	Length  int
	cipher  cipher.Block
	iv      []byte
	Padding int
}

func memset(s []byte, c byte) {
	for i := 0; i < len(s); i++ {
		s[i] = c
	}
}

func NewSecString(str []byte) (*SecString, error) {
	ret := &SecString{Length: len(str)}
	var err error

	if padding := ret.Length % aes.BlockSize; padding != 0 {
		ret.Padding = aes.BlockSize - padding
	}

	ret.String, err = syscall.Mmap(-1, 0, ret.Length+ret.Padding, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON)
	if err != nil {
		return nil, err
	}

	if err := syscall.Mlock(ret.String); err != nil {
		syscall.Munmap(ret.String)
		return nil, err
	}

	for i := 0; i < ret.Length; i++ {
		ret.String[i] = str[i]
		str[i] = 0
	}

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		memset(ret.String, 0)
		syscall.Munmap(ret.String)
		return nil, err
	}

	if ret.cipher, err = aes.NewCipher(key); err != nil {
		memset(ret.String, 0)
		syscall.Munmap(ret.String)
		return nil, err
	}

	if err := ret.Encrypt(); err != nil {
		memset(ret.String, 0)
		syscall.Munmap(ret.String)
		return nil, err
	}

	return ret, nil
}

func (s *SecString) Encrypt() error {
	if err := syscall.Mprotect(s.String, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		return err
	}

	s.iv = make([]byte, aes.BlockSize)
	if _, err := rand.Read(s.iv); err != nil {
		return err
	}

	encrypter := cipher.NewOFB(s.cipher, s.iv)
	encrypter.XORKeyStream(s.String, s.String)

	if err := syscall.Mprotect(s.String, syscall.PROT_READ); err != nil {
		return err
	}

	return nil
}

func (s *SecString) Decrypt() error {
	if err := syscall.Mprotect(s.String, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		return err
	}

	decrypter := cipher.NewOFB(s.cipher, s.iv)
	decrypter.XORKeyStream(s.String, s.String)

	if err := syscall.Mprotect(s.String, syscall.PROT_READ); err != nil {
		return err
	}

	return nil
}

func FromString(str *string) (*SecString, error) {
	b := make([]byte, len(*str))
	for i := 0; i < len(*str); i++ {
		b[i] = (*str)[i]
	}
	*str = strings.Repeat("x", len(*str))
	return NewSecString(b)
}

func (s *SecString) Destroy() error {
	if err := syscall.Mprotect(s.String, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		return err
	}

	memset(s.String, 0)

	if err := syscall.Munlock(s.String); err != nil {
		return err
	}

	if err := syscall.Munmap(s.String); err != nil {
		return err
	}

	s.String = nil
	return nil
}
