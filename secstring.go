package secstring

import "syscall"

type SecString struct {
	String []byte
	Length int
}

func memset(s []byte, c byte) {
	for i := 0; i < len(s); i++ {
		s[i] = c
	}
}

func NewSecString(str []byte) (*SecString, error) {
	ret := &SecString{Length: len(str)}
	var err error
	ret.String, err = syscall.Mmap(-1, 0, ret.Length, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_ANON)
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

	if err := syscall.Mprotect(ret.String, syscall.PROT_READ); err != nil {
		syscall.Munmap(ret.String)
		memset(ret.String, 0)
		return nil, err
	}

	return ret, nil
}

func (s *SecString) Destroy() (error) {
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

	return nil
}
