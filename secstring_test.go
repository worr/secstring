package secstring

import "testing"

func TestNew(t *testing.T) {
	b := []byte("test")
	c := []byte("test")
	s, err := NewSecString(b)
	if s == nil {
		t.Error("New: Expected non-nil SecString. Got nil")
	}

	if err != nil {
		t.Errorf("New: Expected nil err. Got %v", err)
	}
	defer s.Destroy()

	for i := 0; i < len(b); i++ {
		if b[i] != 0 {
			t.Error("New: input should have been cleared")
			t.Errorf("offending byte %v at position %v", b[i], i)
		}
	}

	count := 0
	for i := 0; i < len(c); i++ {
		if c[i] == s.String[i] {
			count++
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

func TestFromString(t *testing.T) {
	b := "test"
	c := []byte("test")
	s, err := FromString(&b)
	if s == nil {
		t.Error("FromString: Expected non-nil SecString. Got nil")
	}

	if err != nil {
		t.Errorf("FromString: Expected nil err. Got %v", err)
	}
	defer s.Destroy()

	for i := 0; i < s.Length; i++ {
		if b[i] != uint8('x') {
			t.Error("FromString: input should have been cleared")
			t.Errorf("offending byte %v at position %v", b[i], i)
		}
	}

	if s.Length != 4 {
		t.Errorf("FromString: Expected length 4. Got %v", s.Length)
	}

	count := 0
	for i := 0; i < len(c); i++ {
		if c[i] == s.String[i] {
			count++
		}
	}
}
