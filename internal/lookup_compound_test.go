package internal

import (
	"encoding/json"
	"os"
	"testing"
)

type Entry struct {
	Typo    string `json:"typo"`
	Bigram  Result `json:"bigram"`
	Unigram Result `json:"unigram"`
}

type Result struct {
	NumResults int    `json:"num_results"`
	Term       string `json:"term"`
	Distance   int    `json:"distance"`
	Count      int    `json:"count"`
}

func TestLookupCompound(t *testing.T) {
	// Open and parse the JSON file
	file, err := os.Open("./tests/lookup_compound_data.json")
	if err != nil {
		t.Fatalf("Failed to open JSON file: %v", err)
	}
	defer file.Close()

	var data struct {
		Data []Entry `json:"data"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}

	// Create a SymSpell instance
	symSpell, _ := NewSymSpell(
		WithCountThreshold(1),
		WithMaxDictionaryEditDistance(3),
		WithPrefixLength(10))
	_, _ = symSpell.LoadDictionary("./tests/vocab.txt", 0, 1, " ")
	_, _ = symSpell.LoadBigramDictionary("./tests/vocab_bigram.txt", 0, 2, "")

	// Run test cases
	for _, entry := range data.Data {
		// Test for bigram
		results := symSpell.LookupCompound(entry.Typo, 3)
		if len(results) != entry.Bigram.NumResults {
			t.Errorf("Bigram: Expected %d results, got %d for typo '%s'", entry.Bigram.NumResults, len(results), entry.Typo)
		} else {
			if results[0].Term != entry.Bigram.Term {
				t.Errorf("Bigram: Expected term '%s', got '%s' for typo '%s'", entry.Bigram.Term, results[0].Term, entry.Typo)
			}
			if results[0].Distance != entry.Bigram.Distance {
				t.Errorf("Bigram: Expected distance %d, got %d for typo '%s'", entry.Bigram.Distance, results[0].Distance, entry.Typo)
			}
			if results[0].Count != entry.Bigram.Count {
				t.Errorf("Bigram: Expected count %d, got %d for typo '%s'", entry.Bigram.Count, results[0].Count, entry.Typo)
			}
		}

		// Test for unigram
		results = symSpell.LookupCompound(entry.Typo, 2)
		if len(results) != entry.Unigram.NumResults {
			t.Errorf("Unigram: Expected %d results, got %d for typo '%s'", entry.Unigram.NumResults, len(results), entry.Typo)
		} else {
			if results[0].Term != entry.Unigram.Term {
				t.Errorf("Unigram: Expected term '%s', got '%s' for typo '%s'", entry.Unigram.Term, results[0].Term, entry.Typo)
			}
			if results[0].Distance != entry.Unigram.Distance {
				t.Errorf("Unigram: Expected distance %d, got %d for typo '%s'", entry.Unigram.Distance, results[0].Distance, entry.Typo)
			}
			if results[0].Count != entry.Unigram.Count {
				t.Errorf("Unigram: Expected count %d, got %d for typo '%s'", entry.Unigram.Count, results[0].Count, entry.Typo)
			}
		}
	}
}
