package secstring

import (
	"errors"
	"crypto/rand" // mock
	"syscall" // mock
	"testing"
	"code.google.com/p/gomock/gomock"
)

func TestNewBadMmap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	syscall.MOCK().SetController(ctrl)
	syscall.EXPECT().Mmap(0, int64(0), 16, 3, 4098).Return(nil, errors.New("error"))
	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadMmap: Expected non-nil error. Got nil")
	}
}

func TestNewBadMlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	syscall.MOCK().SetController(ctrl)
	syscall.EXPECT().Mmap(0, int64(0), 16, 3, 4098).Return([]byte("test"), nil)
	syscall.EXPECT().Mlock([]byte("test")).Return(errors.New("error"))
	syscall.EXPECT().Munmap([]byte("test")).Return(nil)
	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadMlock: Expected non-nil error. Got nil")
	}
	checkScrubbed("NewBadRand", b, t)
}

func checkScrubbed(test string, val []byte, t *testing.T) {
	for i := 0; i < len(val); i++ {
		if val[i] != 0 {
			t.Errorf("%v: b not scrubbed. b: %v", test, val)
		}
	}
}

func TestNewBadRand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	syscall.MOCK().SetController(ctrl)
	rand.MOCK().SetController(ctrl)
	syscall.EXPECT().Mmap(0, int64(0), 16, 3, 4098).Return([]byte("test"), nil)
	syscall.EXPECT().Mlock([]byte("test")).Return(nil)
	syscall.EXPECT().Munmap(make([]byte, 4)).Return(nil)
	rand.EXPECT().Read(make([]byte, 32)).Return(0, errors.New("error"))
	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadRand: Expected non-nil error. Got nil")
	}
	checkScrubbed("NewBadRand", b, t)
}

func TestNewBadMprotect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	k := make([]byte, 32)
	syscall.MOCK().SetController(ctrl)
	rand.MOCK().SetController(ctrl)

	syscall.EXPECT().Mmap(0, int64(0), 16, 3, 4098).Return([]byte("test"), nil)
	syscall.EXPECT().Mlock(b).Return(nil)
	syscall.EXPECT().Munmap(make([]byte, 4)).Return(nil)
	syscall.EXPECT().Mprotect([]byte("test"), 3).Return(errors.New("error"))
	rand.EXPECT().Read(k).Return(0, nil)

	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadMprotect: Expected non-nil error. Got nil.")
	}
	checkScrubbed("NewBadRand", b, t)
}
