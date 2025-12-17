package tui

import (
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jvanrhyn/woordsoek/internal/woordsoek"
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

func TestInputFlowIncludesLanguage(t *testing.T) {
	m := InitializeModel(Flags{})

	// single letter
	m.inputs[0].SetValue("a")
	res, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)
	if m.currentState != inputSixCharString {
		t.Fatalf("expected state %v after first enter; got %v", inputSixCharString, m.currentState)
	}
	if m.focusedInput != 1 {
		t.Fatalf("expected focused input 1; got %d", m.focusedInput)
	}

	// six char string
	m.inputs[1].SetValue("abc")
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)
	if m.currentState != inputLength {
		t.Fatalf("expected state %v after second enter; got %v", inputLength, m.currentState)
	}
	if m.focusedInput != 2 {
		t.Fatalf("expected focused input 2; got %d", m.focusedInput)
	}

	// length
	m.inputs[2].SetValue("3")
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)
	if m.currentState != inputLang {
		t.Fatalf("expected state %v after third enter; got %v", inputLang, m.currentState)
	}
	if m.focusedInput != 3 {
		t.Fatalf("expected focused input 3; got %d", m.focusedInput)
	}

	// create a small dictionary for testlang
	dictDir := "dictionaries"
	_ = os.MkdirAll(dictDir, 0o755)
	testFile := filepath.Join(dictDir, "testlang.txt")
	defer func() { _ = os.Remove(testFile) }()
	if err := os.WriteFile(testFile, []byte("abc\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	// clear cache so subsequent load picks up the new file contents
	woordsoek.ClearDictionaryCache(testFile)

	// language
	m.inputs[3].SetValue("testlang")
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)
	if m.currentState != done {
		t.Fatalf("expected done state after language enter; got %v", m.currentState)
	}
	if len(m.results) == 0 {
		t.Fatalf("expected results to be non-empty after search; got %v", m.results)
	}
}
