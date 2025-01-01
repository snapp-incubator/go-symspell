package internal

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	seperno "github.com/snapp-incubator/seperno"
	"github.com/snapp-incubator/symspell/internal/pkg/edit_distance"
)

// SymSpell represents the Symmetric Delete spelling correction algorithm.
type SymSpell struct {
	MaxDictionaryEditDistance int
	PrefixLength              int
	CountThreshold            int
	Words                     map[string]int
	BelowThresholdWords       map[string]int
	Deletes                   map[string][]string
	maxLength                 int
	distanceComparer          edit_distance.IEditDistance
	normalizer                func(string) string
	// lookup compound
	ReplacedWords  map[string]SuggestItem
	N              float64
	Bigrams        map[string]int
	BigramCountMin int
}

// SuggestItem represents a suggestion with distance and count (placeholder).
type SuggestItem struct {
	Term     string
	Distance int
	Count    int
}

// NewSymSpell is the constructor for the SymSpell struct.
func NewSymSpell(opt ...Options) (*SymSpell, error) {
	opts := DefaultOptions
	for _, config := range opt {
		config.Apply(&opts)
	}
	if opts.MaxDictionaryEditDistance < 0 {
		return nil, errors.New("maxDictionaryEditDistance cannot be negative")
	}
	if opts.PrefixLength < 1 {
		return nil, errors.New("prefixLength cannot be less than 1")
	}
	if opts.PrefixLength <= opts.MaxDictionaryEditDistance {
		return nil, errors.New("prefixLength must be greater than maxDictionaryEditDistance")
	}
	if opts.CountThreshold < 0 {
		return nil, errors.New("countThreshold cannot be negative")
	}
	normalizer := seperno.NewNormalize(
		seperno.WithOuterSpaceRemover(),
		seperno.WithOuterSpaceRemover(),
		seperno.WithURLRemover(),
		seperno.WithNormalizePunctuations(),
		seperno.WithEndsWithEndOfLineChar(),
		seperno.WithConvertHalfSpaceToSpace(),
	)

	return &SymSpell{
		MaxDictionaryEditDistance: opts.MaxDictionaryEditDistance,
		PrefixLength:              opts.PrefixLength,
		CountThreshold:            opts.CountThreshold,
		Words:                     make(map[string]int),
		BelowThresholdWords:       make(map[string]int),
		Deletes:                   make(map[string][]string),
		distanceComparer:          edit_distance.NewEditDistance(edit_distance.DamerauLevenshtein), // todo add more edit distance algorithms
		maxLength:                 0,
		Bigrams:                   make(map[string]int),
		ReplacedWords:             make(map[string]SuggestItem),
		N:                         math.MaxInt,
		BigramCountMin:            0,
		normalizer:                normalizer.BasicNormalizer,
	}, nil
}

// createDictionaryEntry creates or updates an entry in the dictionary.
func (s *SymSpell) createDictionaryEntry(key string, count int) bool {
	if count <= 0 {
		// Early return if count is zero or less
		if s.CountThreshold > 0 {
			return false
		}
		count = 0
	}

	// Check below-threshold words
	if s.CountThreshold > 1 {
		if countPrev, found := s.BelowThresholdWords[key]; found {
			// Increment the count
			count = incrementCount(count, countPrev)
			// Check if it reaches the threshold
			if count >= s.CountThreshold {
				delete(s.BelowThresholdWords, key)
			} else {
				s.BelowThresholdWords[key] = count
				return false
			}
		}
	}

	// Check existing words
	if countPrev, found := s.Words[key]; found {
		// Increment the count
		s.Words[key] = incrementCount(count, countPrev)
		return false
	} else if count < s.CountThreshold {
		// Add to below-threshold words
		s.BelowThresholdWords[key] = count
		return false
	}

	// Add a new word
	s.Words[key] = count

	// Update max length
	if len(key) > s.maxLength {
		s.maxLength = len(key)
	}

	// Create deletes
	edits := s.editsPrefix(key)
	for _, deleteWord := range edits {
		s.Deletes[deleteWord] = append(s.Deletes[deleteWord], key)
	}

	return true
}

// editsPrefix generates edits for a given word (placeholder implementation).
func (s *SymSpell) editsPrefix(word string) []string {
	// Placeholder logic for generating edits (deletes)
	var edits []string
	for i := 0; i < len(word); i++ {
		edits = append(edits, word[:i]+word[i+1:])
	}
	return edits
}

// LoadDictionary loads dictionary entries from a file.
func (s *SymSpell) LoadDictionary(corpusPath string, termIndex int, countIndex int, separator string) (bool, error) {
	if corpusPath == "" {
		return false, errors.New("corpus path cannot be empty")
	}

	// Check if the file exists
	if _, err := os.Stat(corpusPath); os.IsNotExist(err) {
		log.Printf("Dictionary file not found at %s.\n", corpusPath)
		return false, nil
	}

	// Open the file
	file, err := os.Open(corpusPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Load dictionary data from file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, separator)
		if len(fields) <= max(termIndex, countIndex) {
			continue // Skip invalid lines
		}

		term := fields[termIndex]
		count, err := strconv.Atoi(fields[countIndex])
		if err != nil {
			continue // Skip invalid counts
		}
		s.createDictionaryEntry(s.normalizer(term), count)
	}

	if err = scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}
