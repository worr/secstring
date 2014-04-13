package secstring

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"strings"
	"syscall"
)

type SecString struct {
	String  []byte // Encrypted/Decrypted string
	Length  int    // Length of the target string
	Padding int    // Length of padding added
	iv      []byte
	cipher  cipher.Block
	encrypted bool
}

func memset(s []byte, c byte) {
	for i := 0; i < len(s); i++ {
		s[i] = c
	}
}

// Takes a []byte and builds a SecString out of it, wiping str in the
// process. The SecString.String is encrypted after this function.
//
// A SecString should be destroyed when it's no longer needed to prevent memory leaks.
// It is probably a good idea to defer SecString.Destroy()
func NewSecString(str []byte) (*SecString, error) {
	ret := &SecString{Length: len(str)}
	var err error

	if padding := ret.Length % aes.BlockSize; padding != 0 {
		ret.Padding = aes.BlockSize - padding
	}

	ret.String, err = syscall.Mmap(0, 0, ret.Length+ret.Padding, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON|syscall.MAP_PRIVATE)
	if err != nil {
		memset(str, 0)
		return nil, err
	}

	if err := syscall.Mlock(ret.String); err != nil {
		memset(str, 0)
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

// Makes a new SecString from a string reference. Destroys str after creating
// the secstring
func FromString(str *string) (*SecString, error) {
	b := make([]byte, len(*str))
	for i := 0; i < len(*str); i++ {
		b[i] = (*str)[i]
	}
	*str = strings.Repeat("x", len(*str))
	return NewSecString(b)
}

// Encrypts SecString.String
func (s *SecString) Encrypt() error {
	if s.encrypted {
		return errors.New("String already encrypted")
	}

	if err := syscall.Mprotect(s.String, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		return err
	}

	s.iv = make([]byte, aes.BlockSize)
	if _, err := rand.Read(s.iv); err != nil {
		return err
	}

	encrypter := cipher.NewOFB(s.cipher, s.iv)
	encrypter.XORKeyStream(s.String, s.String)
	s.encrypted = true

	if err := syscall.Mprotect(s.String, syscall.PROT_READ); err != nil {
		return err
	}

	return nil
}

// Decrypt SecString.String for use
func (s *SecString) Decrypt() error {
	if ! s.encrypted {
		return errors.New("String already decrypted")
	}

	if err := syscall.Mprotect(s.String, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		return err
	}

	decrypter := cipher.NewOFB(s.cipher, s.iv)
	decrypter.XORKeyStream(s.String, s.String)
	s.encrypted = false

	if err := syscall.Mprotect(s.String, syscall.PROT_READ); err != nil {
		return err
	}

	return nil
}

// Destroys the s. *MUST* be called to prevent memory leaks. Probably best to
// be called in a defer
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
