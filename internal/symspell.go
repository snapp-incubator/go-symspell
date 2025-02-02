package internal

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/snapp-incubator/go-symspell/pkg/edit_distance"
	"github.com/snapp-incubator/go-symspell/pkg/options"
)

// SymSpell represents the Symmetric Delete spelling correction algorithm.
type SymSpell struct {
	MaxDictionaryEditDistance int
	PrefixLength              int
	CountThreshold            int
	SplitThreshold            int
	PreserveCase              bool
	SplitWordBySpace          bool
	Words                     map[string]int
	BelowThresholdWords       map[string]int
	Deletes                   map[string][]string
	maxLength                 int
	distanceComparer          edit_distance.IEditDistance
	// lookup compound
	N              float64
	Bigrams        map[string]int
	BigramCountMin int
}

// NewSymSpell is the constructor for the SymSpell struct.
func NewSymSpell(opt ...options.Options) (*SymSpell, error) {
	opts := options.DefaultOptions
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

	return &SymSpell{
		MaxDictionaryEditDistance: opts.MaxDictionaryEditDistance,
		PrefixLength:              opts.PrefixLength,
		CountThreshold:            opts.CountThreshold,
		SplitThreshold:            opts.SplitItemThreshold,
		PreserveCase:              opts.PreserveCase,
		SplitWordBySpace:          opts.SplitWordBySpace,
		Words:                     make(map[string]int),
		BelowThresholdWords:       make(map[string]int),
		Deletes:                   make(map[string][]string),
		distanceComparer:          edit_distance.NewEditDistance(edit_distance.DamerauLevenshtein), // todo add more edit distance algorithms
		maxLength:                 0,
		Bigrams:                   make(map[string]int),
		N:                         1024908267229,
		BigramCountMin:            math.MaxInt,
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
	} else if countPrev, found := s.Words[key]; found {
		// Increment the count
		s.Words[key] = incrementCount(count, countPrev)
		return false
	}
	if count < s.CountThreshold {
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
	for deleteWord := range edits {
		s.Deletes[deleteWord] = append(s.Deletes[deleteWord], key)
	}

	return true
}

func (s *SymSpell) edits(word string, editDistance int, deleteWords map[string]bool, currentDistance int) {
	editDistance++
	runes := []rune(word)
	if len(runes) == 0 {
		if utf8.RuneCountInString(word) <= s.MaxDictionaryEditDistance {
			deleteWords[""] = true
		}
		return
	}
	for i := currentDistance; i < len(runes); i++ {
		deleteWord := string(runes[:i]) + string(runes[i+1:])
		if !deleteWords[deleteWord] {
			deleteWords[deleteWord] = true
		}
		if editDistance < s.MaxDictionaryEditDistance {
			s.edits(deleteWord, editDistance, deleteWords, i)
		}
	}
}

// editsPrefix function corresponds to _edits_prefix in Python, handling Unicode characters correctly
func (s *SymSpell) editsPrefix(key string) map[string]bool {
	hashSet := make(map[string]bool)
	if utf8.RuneCountInString(key) <= s.MaxDictionaryEditDistance {
		hashSet[""] = true
	}
	runes := []rune(key)
	if len(runes) > s.PrefixLength {
		key = string(runes[:s.PrefixLength])
	}
	hashSet[key] = true
	s.edits(key, 0, hashSet, 0)
	return hashSet
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
		s.createDictionaryEntry(term, count)
	}

	if err = scanner.Err(); err != nil {
		return false, err
	}

	return true, nil
}

func incrementCount(count, countPrevious int) int {
	// Ensure the count does not exceed the maximum value for int64
	if math.MaxInt64-countPrevious > count {
		return countPrevious + count
	}
	return math.MaxInt64
}
