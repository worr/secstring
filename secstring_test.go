package secstring

import "testing"

func TestNew(t *testing.T) {
	b := []byte("test")
	s, err := NewSecString(b)
	if s == nil {
		t.Error("NewSecString should not return nil")
	}

	if err != nil {
		t.Error("err should be nil")
		t.Errorf("err was %v", err)
	}

	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			t.Error("input should have been cleared")
		}
	}
}
