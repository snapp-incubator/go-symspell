package symspell

import (
	"log"

	"github.com/snapp-incubator/go-symspell/internal"
	"github.com/snapp-incubator/go-symspell/internal/verbosity"
)

func NewSymSpell(opt ...internal.Options) SymSpell {
	symspell, err := internal.NewSymSpell(opt...)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	return symspell

}

// NewSymSpellWithLoadDictionary used when want Lookup only
func NewSymSpellWithLoadDictionary(dirPath string, termIndex, countIndex int, opt ...internal.Options) SymSpell {
	symspell := NewSymSpell(opt...)
	ok, err := symspell.LoadDictionary(dirPath, termIndex, countIndex, " ")
	if err != nil {
		log.Fatal("[Error] ", err)
	}
	if !ok {
		log.Fatal("[Error] loading dictionary has been failed")
	}
	return symspell
}

func NewSymSpellWithLoadBigramDictionary(vocabDirPath, bigramDirPath string, termIndex, countIndex int, opt ...internal.Options) SymSpell {
	symspell := NewSymSpell(opt...)
	ok, err := symspell.LoadDictionary(vocabDirPath, termIndex, countIndex, " ")
	if err != nil || !ok {
		log.Fatal("[Error] ", err)
	}
	ok, err = symspell.LoadBigramDictionary(bigramDirPath, termIndex, countIndex+1, "")
	if err != nil || !ok {
		if bigramDirPath != "" {
			log.Println("[Error] ", err)
		}
	}
	return symspell
}

func WithMaxDictionaryEditDistance(maxDictionaryEditDistance int) internal.Options {
	return internal.WithMaxDictionaryEditDistance(maxDictionaryEditDistance)
}

func WithPrefixLength(prefixLength int) internal.Options {
	return internal.WithPrefixLength(prefixLength)
}

func WithCountThreshold(countThreshold int) internal.Options {
	return internal.WithCountThreshold(countThreshold)
}

const (
	Top     = verbosity.Top
	Closest = verbosity.Closest
	All     = verbosity.All
)

type SymSpell interface {
	Lookup(phrase string, verbosity verbosity.Verbosity, maxEditDistance int) ([]internal.SuggestItem, error)
	LookupCompound(phrase string, maxEditDistance int) *internal.SuggestItem
	LoadBigramDictionary(corpusPath string, termIndex, countIndex int, separator string) (bool, error)
	LoadDictionary(corpusPath string, termIndex int, countIndex int, separator string) (bool, error)
}
