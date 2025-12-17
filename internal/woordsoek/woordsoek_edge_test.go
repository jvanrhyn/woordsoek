package woordsoek

import (
	"fmt"
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
