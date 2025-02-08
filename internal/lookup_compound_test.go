package internal

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/snapp-incubator/go-symspell/pkg/options"
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
		options.WithCountThreshold(1),
		options.WithMaxDictionaryEditDistance(3),
		options.WithPrefixLength(10),
		options.WithSplitItemThreshold(1))
	_, _ = symSpell.LoadDictionary("./tests/vocab.txt", 0, 1, " ")

	// Run test cases
	for _, entry := range data.Data {
		// Test for unigram
		results := symSpell.LookupCompound(entry.Typo, 2)
		if results == nil {
			t.Errorf("Unigram: Expected %d results, got %v for typo '%s'", entry.Unigram.NumResults, results, entry.Typo)
		} else {
			if results.Term != entry.Unigram.Term {
				t.Errorf("Unigram: Expected term '%s', got '%s' for typo '%s'", entry.Unigram.Term, results.Term, entry.Typo)
			}
			if results.Distance != entry.Unigram.Distance {
				t.Errorf("Unigram: Expected distance %d, got %d for typo '%s'", entry.Unigram.Distance, results.Distance, entry.Typo)
			}
			if results.Count != entry.Unigram.Count {
				t.Errorf("Unigram: Expected count %d, got %d for typo '%s'", entry.Unigram.Count, results.Count, entry.Typo)
			}
		}
	} // Run test cases

	_, _ = symSpell.LoadBigramDictionary("./tests/vocab_bigram.txt", 0, 2, "")
	for _, entry := range data.Data {
		// Test for bigram
		results := symSpell.LookupCompound(entry.Typo, 3)
		if results == nil {
			t.Errorf("Bigram: Expected %d results, got %v for typo '%s'", entry.Bigram.NumResults, results, entry.Typo)
		} else {
			if results.Term != entry.Bigram.Term {
				t.Errorf("Bigram: Expected term '%s', got '%s' for typo '%s'", entry.Bigram.Term, results.Term, entry.Typo)
			}
			if results.Distance != entry.Bigram.Distance {
				t.Errorf("Bigram: Expected distance %d, got %d for typo '%s'", entry.Bigram.Distance, results.Distance, entry.Typo)
			}
			if results.Count != entry.Bigram.Count {
				t.Errorf("Bigram: Expected count %d, got %d for typo '%s'", entry.Bigram.Count, results.Count, entry.Typo)
			}
		}
	}
}

func Test_separateNumbers(t *testing.T) {
	type args struct {
		inputs string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "first number",
			args: args{
				inputs: "15خرداد",
			},
			want: []string{"15", "خرداد"},
		},
		{
			name: "first word",
			args: args{
				inputs: "خرداد15",
			},
			want: []string{"خرداد", "15"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitWordAndNumber(tt.args.inputs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("separateNumbers() = %v, want %v", got, tt.want)
			}
		})
	}
}
