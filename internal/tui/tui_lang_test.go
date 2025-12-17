package tui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTUI_LangOverridesEnv(t *testing.T) {
	dictDir := "dictionaries"
	_ = os.MkdirAll(dictDir, 0o755)
	testFile := filepath.Join(dictDir, "testlang.txt")
	otherFile := filepath.Join(dictDir, "otherlang.txt")
	defer func() {
		_ = os.Remove(testFile)
		_ = os.Remove(otherFile)
	}()

	if err := os.WriteFile(testFile, []byte("alpha\nbeta\n"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	if err := os.WriteFile(otherFile, []byte("other\n"), 0644); err != nil {
		t.Fatalf("failed to write other file: %v", err)
	}

	// set environment to otherlang
	old := os.Getenv("WBLANG")
	_ = os.Setenv("WBLANG", "otherlang")
	defer func() { _ = os.Setenv("WBLANG", old) }()

	m := InitializeModel(Flags{})
	m.flags.SingleLetter = "a"
	m.flags.SixCharString = "lph"
	m.flags.Length = 0
	m.flags.Lang = "testlang"

	m = m.searchWords()
	found := false
	for _, w := range m.results {
		if w == "alpha" {
			found = true
		}
		if w == "other" {
			t.Fatalf("unexpectedly found 'other' from env override when flags.Lang was set")
		}
	}
	if !found {
		t.Fatalf("expected 'alpha' in results; got %v", m.results)
	}
}
