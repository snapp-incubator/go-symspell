package internal

import (
	"testing"

	verbositypkg "github.com/snapp-incubator/go-symspell/internal/verbosity"
)

func TestDeletes(t *testing.T) {
	symSpell, _ := NewSymSpell(
		WithCountThreshold(1),
		WithMaxDictionaryEditDistance(2),
		WithPrefixLength(7),
	)
	symSpell.createDictionaryEntry("steama", 4)
	symSpell.createDictionaryEntry("steamb", 6)
	symSpell.createDictionaryEntry("steamc", 2)

	results, err := symSpell.Lookup("stream", verbositypkg.Top, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Term != "steamb" {
		t.Errorf("Expected term 'steamb', got '%s'", results[0].Term)
	}

	if results[0].Count != 6 {
		t.Errorf("Expected count 6, got %d", results[0].Count)
	}
}

func TestWordsWithSharedPrefixShouldRetainCounts(t *testing.T) {
	symSpell, err := NewSymSpell(
		WithCountThreshold(4),
		WithMaxDictionaryEditDistance(1),
		WithPrefixLength(3),
	)
	symSpell.createDictionaryEntry("pipe", 5)
	symSpell.createDictionaryEntry("pips", 10)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Test for "pipe"
	results, err := symSpell.Lookup("pipe", verbositypkg.All, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Term != "pipe" || results[0].Count != 5 {
		t.Errorf("Expected first result to be 'pipe' with count 5, got '%s' with count %d", results[0].Term, results[0].Count)
	}
	if results[1].Term != "pips" || results[1].Count != 10 {
		t.Errorf("Expected second result to be 'pips' with count 10, got '%s' with count %d", results[1].Term, results[1].Count)
	}

	// Test for "pips"
	results, err = symSpell.Lookup("pips", verbositypkg.All, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Term != "pips" || results[0].Count != 10 {
		t.Errorf("Expected first result to be 'pips' with count 10, got '%s' with count %d", results[0].Term, results[0].Count)
	}
	if results[1].Term != "pipe" || results[1].Count != 5 {
		t.Errorf("Expected second result to be 'pipe' with count 5, got '%s' with count %d", results[1].Term, results[1].Count)
	}

	// Test for "pip"
	results, err = symSpell.Lookup("pip", verbositypkg.All, 1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	if results[0].Term != "pips" || results[0].Count != 10 {
		t.Errorf("Expected first result to be 'pips' with count 10, got '%s' with count %d", results[0].Term, results[0].Count)
	}
	if results[1].Term != "pipe" || results[1].Count != 5 {
		t.Errorf("Expected second result to be 'pipe' with count 5, got '%s' with count %d", results[1].Term, results[1].Count)
	}
}

func TestVerbosityShouldControlLookupResults(t *testing.T) {
	symSpell, err := NewSymSpell(WithCountThreshold(0))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	symSpell.createDictionaryEntry("steam", 1)
	symSpell.createDictionaryEntry("steams", 2)
	symSpell.createDictionaryEntry("steem", 3)
	tests := []struct {
		verbosity  verbositypkg.Verbosity
		numResults int
	}{
		{verbositypkg.Top, 1},
		{verbositypkg.Closest, 2},
		{verbositypkg.All, 3},
	}

	for _, test := range tests {
		results, err := symSpell.Lookup("steems", test.verbosity, 2)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(results) != test.numResults {
			t.Errorf("Expected %d results for verbosity %d, got %d", test.numResults, test.verbosity, len(results))
		}
	}
}

func TestShouldReturnMostFrequent(t *testing.T) {
	symSpell, _ := NewSymSpell(WithCountThreshold(1), WithMaxDictionaryEditDistance(2))
	symSpell.createDictionaryEntry("steama", 4)
	symSpell.createDictionaryEntry("steamb", 6)
	symSpell.createDictionaryEntry("steamc", 2)

	results, err := symSpell.Lookup("stream", verbositypkg.Top, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Term != "steamb" {
		t.Errorf("Expected term 'steamb', got '%s'", results[0].Term)
	}

	if results[0].Count != 6 {
		t.Errorf("Expected count 6, got %d", results[0].Count)
	}
}

func TestShouldFindExactMatch(t *testing.T) {
	symSpell, err := NewSymSpell(WithCountThreshold(1))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	symSpell.createDictionaryEntry("steama", 4)
	symSpell.createDictionaryEntry("steamb", 6)
	symSpell.createDictionaryEntry("steamc", 2)

	results, err := symSpell.Lookup("streama", verbositypkg.Top, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Term != "steama" {
		t.Errorf("Expected term 'steama', got '%s'", results[0].Term)
	}
}

func TestShouldNotReturnNonWordDelete(t *testing.T) {
	symSpell, err := NewSymSpell(WithCountThreshold(10))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	symSpell.createDictionaryEntry("pawn", 10)

	terms := []string{"paw", "awn"}
	for _, term := range terms {
		results, err := symSpell.Lookup(term, verbositypkg.Top, 0)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("Expected no results for term '%s', but got %d results", term, len(results))
		}
	}
}

func TestShouldNotReturnLowCountWord(t *testing.T) {
	symSpell, _ := NewSymSpell(WithCountThreshold(10))
	symSpell.createDictionaryEntry("pawn", 1)

	results, err := symSpell.Lookup("pawn", verbositypkg.Closest, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected no results, got %d", len(results))
	}
}

func TestShouldNotReturnLowCountWordThatAreAlsoDeleteWord(t *testing.T) {
	symSpell, _ := NewSymSpell(WithCountThreshold(10))
	symSpell.createDictionaryEntry("flame", 20)
	symSpell.createDictionaryEntry("flam", 1)

	results, err := symSpell.Lookup("flam", verbositypkg.Top, 0)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Expected no results for term 'flam', but got %d results", len(results))
	}
}

func TestMaxEditDistanceTooLarge(t *testing.T) {
	symSpell, _ := NewSymSpell(WithCountThreshold(10))
	symSpell.createDictionaryEntry("flame", 20)
	symSpell.createDictionaryEntry("flam", 1)

	_, err := symSpell.Lookup("flam", verbositypkg.Top, 3)
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}
	if err.Error() != "distance too large" {
		t.Errorf("Expected error 'distance too large', got '%s'", err.Error())
	}
}

func TestTransferCasing(t *testing.T) {
	tests := []struct {
		entries []struct {
			Term  string
			Count int
		}
		typo       string
		correction string
	}{
		{
			entries: []struct {
				Term  string
				Count int
			}{{Term: "steam", Count: 4}},
			typo:       "Stream",
			correction: "steam",
		},
		{
			entries: []struct {
				Term  string
				Count int
			}{{Term: "steam", Count: 4}},
			typo:       "stream",
			correction: "steam",
		},
	}

	for _, test := range tests {
		symSpell, err := NewSymSpell()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		// Add entries to the dictionary
		for _, entry := range test.entries {
			symSpell.createDictionaryEntry(entry.Term, entry.Count)
		}

		// Perform lookup
		results, _ := symSpell.Lookup(test.typo, verbositypkg.Top, 2)
		if len(results) == 0 {
			t.Fatalf("Expected results for typo '%s', got none", test.typo)
		}

		// Check correction
		if results[0].Term != test.correction {
			t.Errorf("For typo '%s', expected correction '%s', got '%s'", test.typo, test.correction, results[0].Term)
		}
	}
}

func TestFarsi(t *testing.T) {
	symSpell, err := NewSymSpell(WithCountThreshold(1))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	symSpell.createDictionaryEntry("تجریش", 4)

	results, err := symSpell.Lookup("تحریش", verbositypkg.Top, 2)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Term != "تجریش" {
		t.Errorf("Expected term 'steama', got '%s'", results[0].Term)
	}
}
