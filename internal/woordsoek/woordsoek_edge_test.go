package woordsoek

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestUnicodeEdgeCases(t *testing.T) {
	LoadVowelForms()
	// combining acute accent: e + \u0301 should be considered distinct unless normalized
	words := []string{"e\u0301cole", "école", "ecole"}
	r := strings.NewReader(strings.Join(words, "\n"))
	res, err := searchFromReader(r, "e", "cole", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// normalization currently replaces precomposed accented vowels but not combining marks,
	// ensure at least one normalized form is present and all results are lowercase
	found := false
	for _, w := range res {
		if strings.ToLower(w) == "ecole" {
			found = true
		}
		if w != strings.ToLower(w) {
			t.Fatalf("result %q not lowercase", w)
		}
	}
	if !found {
		t.Fatalf("expected normalized 'ecole' in results; got %v", res)
	}
}

func TestLargeDictionaryBehavior(t *testing.T) {
	// create a large synthetic dictionary with many entries but keep it reasonable for CI
	var sb strings.Builder
	const total = 20000
	for i := 0; i < total; i++ {
		sb.WriteString(fmt.Sprintf("word%05d\n", i))
	}

	// add some matching words sprinkled in
	sb.WriteString("alpha\nalfha\nalphA\n")

	res, err := searchFromReader(strings.NewReader(sb.String()), "a", "lph", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// expect at least 'alpha' normalized to lowercase
	found := false
	for _, w := range res {
		if w == "alpha" {
			found = true
		}
		if utf8.RuneCountInString(w) < minDefaultLength {
			t.Fatalf("unexpected short word %q", w)
		}
	}
	if !found {
		t.Fatalf("expected 'alpha' in results; got %v", res)
	}
}

func TestDictionaryCacheBehavior(t *testing.T) {
	// create temp file
	tmp := "tmp_cache_test.txt"
	defer func() { _ = os.Remove(tmp) }()
	_ = os.WriteFile(tmp, []byte("one\ntwo\nthree\n"), 0644)

	// ensure cache is clear
	ClearDictionaryCache(tmp)
	w1, err := loadDictionary(tmp)
	if err != nil {
		t.Fatalf("failed to load dictionary: %v", err)
	}
	if len(w1) == 0 {
		t.Fatalf("expected words in dictionary")
	}

	// modify underlying file
	_ = os.WriteFile(tmp, []byte("alpha\nbeta\n"), 0644)

	// without clearing cache, loadDictionary should return old contents
	w2, err := loadDictionary(tmp)
	if err != nil {
		t.Fatalf("failed to load dictionary second time: %v", err)
	}
	if len(w2) != len(w1) {
		t.Fatalf("expected cached results to match initial load; got %v vs %v", len(w2), len(w1))
	}

	// after clearing cache, it should pick up new contents
	ClearDictionaryCache(tmp)
	w3, err := loadDictionary(tmp)
	if err != nil {
		t.Fatalf("failed to load dictionary after clear: %v", err)
	}
	if len(w3) == len(w1) {
		t.Fatalf("expected updated contents after clearing cache")
	}
}

func TestParallelMatchesSerial(t *testing.T) {
	// build synthetic data
	words := makeWords(10000)

	sres, err := searchFromWords(words, "a", "lph", 0)
	if err != nil {
		t.Fatalf("serial search failed: %v", err)
	}
	pres, err := searchFromWordsParallel(words, "a", "lph", 0)
	if err != nil {
		t.Fatalf("parallel search failed: %v", err)
	}
	if len(sres) != len(pres) {
		t.Fatalf("mismatched result lengths serial=%d parallel=%d", len(sres), len(pres))
	}
	for i := range sres {
		if sres[i] != pres[i] {
			t.Fatalf("mismatch at index %d: %q != %q", i, sres[i], pres[i])
		}
	}
}
