package f

import "testing"

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
