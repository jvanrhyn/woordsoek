package tui

import (
	"testing"
	"unicode/utf8"
)

func TestValidateSingleLetter(t *testing.T) {
	m := InitializeModel(Flags{})

	// invalid: empty
	m.inputs[0].SetValue("")
	if utf8.RuneCountInString(m.inputs[0].Value()) != 0 {
		t.Fatalf("expected empty input")
	}

	// invalid: two characters
	m.inputs[0].SetValue("ab")
	if utf8.RuneCountInString(m.inputs[0].Value()) != 2 {
		t.Fatalf("expected two runes")
	}

	// valid single rune
	m.inputs[0].SetValue("x")
	if utf8.RuneCountInString(m.inputs[0].Value()) != 1 {
		t.Fatalf("expected one rune")
	}
}
