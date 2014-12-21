package secstring

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/sys/unix"
)

type SecString struct {
	String []byte // Protected string
	Length int    // Length of the target string
}

func memset(s []byte, c byte) {
	for i := 0; i < len(s); i++ {
		s[i] = c
	}
}

// Takes a []byte and builds a SecString out of it, wiping str in the
// process.
//
// A SecString should be destroyed when it's no longer needed to prevent memory leaks.
// It is probably a good idea to defer SecString.Destroy()
func NewSecString(str []byte) (*SecString, error) {
	ret := &SecString{Length: len(str)}
	var err error

	ret.String, err = unix.Mmap(-1, 0, ret.Length, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE)
	if err != nil {
		memset(str, 0)
		return nil, err
	}

	if err := unix.Mlock(ret.String); err != nil {
		memset(str, 0)
		unix.Munmap(ret.String)
		return nil, err
	}

	for i := 0; i < ret.Length; i++ {
		ret.String[i] = str[i]
		str[i] = 0
	}

	if err := unix.Mprotect(ret.String, unix.PROT_READ); err != nil {
		memset(str, 0)
		memset(ret.String, 0)
		unix.Munmap(ret.String)
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

func (s *SecString) Clone() (*SecString, error) {
	var ret *SecString

	str := make([]byte, len(s.String), cap(s.String))
	if copied := copy(str, s.String); copied != len(s.String) {
		return nil, errors.New(fmt.Sprintf("Only %v copied", copied))
	}

	var err error
	if ret, err = NewSecString(str); err != nil {
		return nil, err
	}

	return ret, nil
}

// Destroys the s. *MUST* be called to prevent memory leaks. Probably best to
// be called in a defer
func (s *SecString) Destroy() error {
	if err := unix.Mprotect(s.String, unix.PROT_READ|unix.PROT_WRITE); err != nil {
		return err
	}

	memset(s.String, 0)

	if err := unix.Munlock(s.String); err != nil {
		return err
	}

	if err := unix.Munmap(s.String); err != nil {
		return err
	}

	s.String = nil
	return nil
}
