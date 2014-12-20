package secstring

import (
	"code.google.com/p/gomock/gomock"
	"errors"
	"golang.org/x/sys/unix" // mock
	"testing"
)

func TestNewBadMmap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	unix.MOCK().SetController(ctrl)
	unix.EXPECT().Mmap(0, int64(0), 16, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE).Return(nil, errors.New("error"))
	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadMmap: Expected non-nil error. Got nil")
	}
}

func TestNewBadMlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := []byte("test")
	unix.MOCK().SetController(ctrl)
	unix.EXPECT().Mmap(0, int64(0), 16, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE).Return([]byte("test"), nil)
	unix.EXPECT().Mlock([]byte("test")).Return(errors.New("error"))
	unix.EXPECT().Munmap([]byte("test")).Return(nil)
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
	unix.MOCK().SetController(ctrl)
	unix.EXPECT().Mmap(0, int64(0), 16, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE).Return([]byte("test"), nil)
	unix.EXPECT().Mlock([]byte("test")).Return(nil)
	unix.EXPECT().Munmap(make([]byte, 4)).Return(nil)
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
	unix.MOCK().SetController(ctrl)

	unix.EXPECT().Mmap(0, int64(0), 16, unix.PROT_READ|unix.PROT_WRITE, unix.MAP_ANON|unix.MAP_PRIVATE).Return([]byte("test"), nil)
	unix.EXPECT().Mlock(b).Return(nil)
	unix.EXPECT().Munmap(make([]byte, 4)).Return(nil)
	unix.EXPECT().Mprotect([]byte("test"), 3).Return(errors.New("error"))

	if _, err := NewSecString(b); err == nil {
		t.Error("NewBadMprotect: Expected non-nil error. Got nil.")
	}
	checkScrubbed("NewBadRand", b, t)
}
