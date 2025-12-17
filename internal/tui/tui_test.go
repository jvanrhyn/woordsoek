package tui

import (
	"os"
	"path/filepath"
	"strings"
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

func TestInitializeModelAcceptsLangFlag(t *testing.T) {
	flags := Flags{Lang: "en"}
	m := InitializeModel(flags)
	if m.flags.Lang != "en" {
		t.Fatalf("expected model to have Lang 'en'; got '%s'", m.flags.Lang)
	}
}

func TestCLIProvidedLangPreservedWhenTUIEmpty(t *testing.T) {
	// Prepare dictionaries and env so we can deterministically check which
	// language was used during the search.
	dictDir := "dictionaries"
	_ = os.MkdirAll(dictDir, 0o755)
	enFile := filepath.Join(dictDir, "en.txt")
	otherFile := filepath.Join(dictDir, "otherlang.txt")
	defer func() {
		_ = os.Remove(enFile)
		_ = os.Remove(otherFile)
	}()

	if err := os.WriteFile(enFile, []byte("alpha\n"), 0644); err != nil {
		t.Fatalf("failed to write en file: %v", err)
	}
	if err := os.WriteFile(otherFile, []byte("other\n"), 0644); err != nil {
		t.Fatalf("failed to write other file: %v", err)
	}
	// ensure cache won't interfere
	woordsoek.ClearDictionaryCache(enFile)
	woordsoek.ClearDictionaryCache(otherFile)

	// set env to otherlang to ensure env would win if CLI flag is ignored
	old := os.Getenv("WBLANG")
	_ = os.Setenv("WBLANG", "otherlang")
	defer func() { _ = os.Setenv("WBLANG", old) }()

	// initialize model with CLI-provided lang 'en'
	flags := Flags{Lang: "en"}
	m := InitializeModel(flags)

	// go through inputs: single letter and six-char selection to reach length
	m.inputs[0].SetValue("a")
	res, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)

	m.inputs[1].SetValue("lph")
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)

	// set length to 0 so default minDefaultLength applies and matches 'alpha'
	m.inputs[2].SetValue("")
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)

	// now we're at language input; leave it empty to simulate user skipping it
	// pressing enter should preserve the CLI-provided Lang ('en')
	res, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = res.(Model)

	// verify that results come from the 'en' dictionary (alpha present)
	foundAlpha := false
	for _, w := range m.results {
		if w == "alpha" {
			foundAlpha = true
		}
		if w == "other" {
			t.Fatalf("unexpectedly found 'other' from env override when CLI flag 'en' should have been used")
		}
	}
	if !foundAlpha {
		t.Fatalf("expected 'alpha' in results; got %v", m.results)
	}
}

func TestViewDisplaysCurrentLanguage(t *testing.T) {
	// case: CLI-provided lang
	m := InitializeModel(Flags{Lang: "en"})
	v := m.View()
	if !strings.Contains(v, "Language: en") {
		t.Fatalf("expected view to show 'Language: en'; got:\n%s", v)
	}

	// case: env-provided lang when flag empty
	old := os.Getenv("WBLANG")
	_ = os.Setenv("WBLANG", "otherlang")
	defer func() { _ = os.Setenv("WBLANG", old) }()
	m2 := InitializeModel(Flags{})
	v2 := m2.View()
	if !strings.Contains(v2, "Language: otherlang") {
		t.Fatalf("expected view to show 'Language: otherlang'; got:\n%s", v2)
	}

	// case: default when neither flag nor env set
	_ = os.Unsetenv("WBLANG")
	m3 := InitializeModel(Flags{})
	v3 := m3.View()
	if !strings.Contains(v3, "Language: af-za") {
		t.Fatalf("expected view to show default 'Language: af-za'; got:\n%s", v3)
	}
}
