package f

import (
	"errors"
	"testing"
)

func TestAbc(t *testing.T) {
	if x := Abc(2); x != 4 {
		t.Error("ERROR")
	}
}

func TestDef(t *testing.T) {
	if x := Def(2); x != 4 {
		t.Error("ERROR")
	}
}
func TestNewEqual(t *testing.T) {
	if errors.New("abc") == errors.New("abc") {
		t.Errorf(`New("abc") == New("abc")`)
	}
}
