package symspell

import (
	"log"

	"github.com/snapp-incubator/go-symspell/internal"
	"github.com/snapp-incubator/go-symspell/pkg/items"
	"github.com/snapp-incubator/go-symspell/pkg/options"
	"github.com/snapp-incubator/go-symspell/pkg/verbosity"
)

func NewSymSpell(opt ...options.Options) SymSpell {
	symspell, err := internal.NewSymSpell(opt...)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}
	return symspell
}

// NewSymSpellWithLoadDictionary used when want Lookup only
func NewSymSpellWithLoadDictionary(dirPath string, termIndex, countIndex int, opt ...options.Options) SymSpell {
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

func NewSymSpellWithLoadBigramDictionary(vocabDirPath, bigramDirPath string, termIndex, countIndex int, opt ...options.Options) SymSpell {
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

type SymSpell interface {
	Lookup(phrase string, verbosity verbosity.Verbosity, maxEditDistance int) ([]items.SuggestItem, error)
	LookupCompound(phrase string, maxEditDistance int) *items.SuggestItem
	LoadBigramDictionary(corpusPath string, termIndex, countIndex int, separator string) (bool, error)
	LoadDictionary(corpusPath string, termIndex int, countIndex int, separator string) (bool, error)
}
