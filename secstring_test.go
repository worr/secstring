package secstring

import "testing"

func TestNew(t *testing.T) {
	b := []byte("test")
	s, err := NewSecString(b)
	if s == nil {
		t.Error("New: Expected non-nil SecString. Got nil")
	}

	if err != nil {
		t.Errorf("New: Expected nil err. Got %v", err)
	}

	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			t.Error("New: input should have been cleared")
			t.Errorf("offending byte at offset %v: %v", i, b[i])
		}
	}

	if s.Length != 4 {
		t.Errorf("New: Expected length 4. Got %v", s.Length)
	}
}

func TestDestroy(t *testing.T) {
	s, err := NewSecString([]byte("test"))
	err = s.Destroy()

	if err != nil {
		t.Errorf("Destroy: Expected nil err. Got %v", err)
	}

	if s.String != nil {
		t.Errorf("Destroy: Expected nil s.String. got %v", s.String)
	}
}
