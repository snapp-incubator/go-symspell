package internal

import (
	"errors"
	verbositypkg "github.com/snapp-incubator/symspell/internal/verbosity"
	"slices"
	"sort"
)

func (s *SymSpell) Lookup(
	phrase string,
	verbosity verbositypkg.Verbosity,
	maxEditDistance int,
) ([]SuggestItem, error) {
	if maxEditDistance > s.MaxDictionaryEditDistance {
		return nil, errors.New("distance too large")
	}
	cp := newCandidateProcessor(maxEditDistance, verbosity, phrase)
	// Early exit - word too big to match any words
	if cp.phraseLen-maxEditDistance > s.maxLength {
		return cp.suggestions, nil
	}

	// Quick look for exact match
	shouldEnd := s.checkExactMatch(phrase, verbosity, &cp)
	if shouldEnd {
		return cp.suggestions, nil
	}

	if maxEditDistance == 0 {
		return cp.suggestions, nil
	}
	cp.consideredSuggestions[phrase] = true
	// Add original prefix
	phrasePrefixRunes := s.getOriginPrefix(&cp)
	cp.candidates = append(cp.candidates, string(phrasePrefixRunes))
	// Process candidates
	s.processCandidate(phrase, maxEditDistance, &cp)

	cp.sortCandidate()

	return cp.suggestions, nil
}

func (s *SymSpell) getOriginPrefix(cp *candidateProcessor) []rune {
	phrasePrefixRunes := cp.phraseRunes
	if cp.phraseLen > s.PrefixLength {
		phrasePrefixRunes = cp.phraseRunes[:s.PrefixLength]
	}
	return phrasePrefixRunes
}

func (s *SymSpell) checkExactMatch(phrase string, verbosity verbositypkg.Verbosity, cp *candidateProcessor) bool {
	if count, found := s.Words[phrase]; found {
		cp.suggestions = append(cp.suggestions, SuggestItem{Term: phrase, Distance: 0, Count: count})
		if verbosity != verbositypkg.All {
			return true
		}
	}
	return false
}

func (s *SymSpell) processCandidate(phrase string, maxEditDistance int, cp *candidateProcessor) {
	for cp.candidatePointer < len(cp.candidates) {
		candidate, candidateRunes := s.preProcessCandidate(cp)

		if cp.lenDiff > cp.maxEditDistance2 {
			if cp.verbosity == verbositypkg.All {
				continue
			}
			break
		}

		// Check suggestions for the candidate
		if dictSuggestions, found := s.Deletes[candidate]; found {
			for _, suggestion := range dictSuggestions {
				if suggestion == phrase {
					continue
				}
				cp.updateSuggestion(suggestion)
				skip := s.checkSuggestionToSkip(cp, suggestion, candidate)
				if skip {
					continue
				}
				cp.resetDistance()
				if cp.candidateLen == 0 {
					cp.distance = max(cp.phraseLen, cp.suggestionLen)
					if cp.distance > cp.maxEditDistance2 || cp.consideredSuggestions[suggestion] {
						continue
					}
				} else if cp.suggestionLen == 1 {
					skip = s.checkFirstRuneDistance(phrase, cp.suggestionRunes, cp, suggestion)
					if skip {
						continue
					}
				} else {
					s.updateMinDistance(maxEditDistance, cp)
					skip = s.checkDistanceToSkip(phrase, maxEditDistance, cp, suggestion)
					if skip {
						continue
					}
				}
				if cp.distance <= cp.maxEditDistance2 {
					s.updateSuggestions(suggestion, cp)
				}
			}
		}
		if cp.lenDiff <= maxEditDistance && cp.candidateLen <= s.PrefixLength {
			if cp.verbosity != verbositypkg.All && cp.lenDiff >= cp.maxEditDistance2 {
				continue
			}
			s.addEditDistance(candidateRunes, cp)
		}
	}
}

func (s *SymSpell) preProcessCandidate(cp *candidateProcessor) (string, []rune) {
	candidate := cp.candidates[cp.candidatePointer]
	cp.candidatePointer++
	candidateRunes := []rune(candidate)
	cp.candidateLen = len(candidateRunes)
	cp.lenDiff = cp.phraseLen - cp.candidateLen
	return candidate, candidateRunes
}

func (s *SymSpell) checkSuggestionToSkip(cp *candidateProcessor, suggestion string, candidate string) bool {
	if abs(cp.suggestionLen-cp.phraseLen) > cp.maxEditDistance2 || cp.suggestionLen < cp.candidateLen ||
		(cp.suggestionLen == cp.candidateLen && suggestion != candidate) {
		return true
	}
	suggestionPrefixLen := min(cp.suggestionLen, s.PrefixLength)
	if suggestionPrefixLen > cp.phraseLen && suggestionPrefixLen-cp.candidateLen > cp.maxEditDistance2 {
		return true
	}
	return false
}

func (s *SymSpell) checkDistanceToSkip(phrase string, maxEditDistance int, cp *candidateProcessor, suggestion string) bool {
	if s.PrefixLength-maxEditDistance == cp.candidateLen {
		skip := s.checkProcessShouldSkip(cp)
		if skip {
			return true
		}
	}
	// delete in suggestion prefix is somewhat expensive, and
	// only pays off when verbosity is TOP or CLOSEST
	if cp.consideredSuggestions[suggestion] {
		return true
	}
	cp.consideredSuggestions[suggestion] = true
	cp.distance = s.distanceCompare(phrase, suggestion, cp.maxEditDistance2)
	if cp.distance < 0 {
		return true
	}
	return false
}

func (s *SymSpell) updateMinDistance(maxEditDistance int, cp *candidateProcessor) {
	if s.PrefixLength-maxEditDistance == cp.candidateLen {
		cp.minDistance = min(cp.phraseLen, cp.suggestionLen) - s.PrefixLength
	} else {
		cp.minDistance = 0
	}
}

func (s *SymSpell) checkFirstRuneDistance(phrase string, suggestionRunes []rune, cp *candidateProcessor, suggestion string) bool {
	var distanceCalc = func() int {
		phraseRunesList := []rune(phrase)
		// Check if the first rune of suggestion exists in phrase
		if slices.Contains(phraseRunesList, suggestionRunes[0]) {
			return cp.phraseLen - 1
		} else {
			return cp.phraseLen
		}
	}
	cp.distance = distanceCalc()
	if cp.distance > cp.maxEditDistance2 || cp.consideredSuggestions[suggestion] {
		return true
	}
	return false
}

func (s *SymSpell) checkProcessShouldSkip(cp *candidateProcessor) bool {
	if cp.minDistance > 1 &&
		string(cp.phraseRunes[cp.phraseLen+1-cp.minDistance:]) != string(cp.suggestionRunes[cp.suggestionLen+1-cp.minDistance:]) {
		return true
	}
	if cp.minDistance > 0 &&
		cp.phraseRunes[cp.phraseLen-cp.minDistance] != cp.suggestionRunes[cp.suggestionLen-cp.minDistance] {
		if cp.phraseRunes[cp.phraseLen-cp.minDistance-1] != cp.suggestionRunes[cp.suggestionLen-cp.minDistance] ||
			cp.phraseRunes[cp.phraseLen-cp.minDistance] != cp.suggestionRunes[cp.suggestionLen-cp.minDistance-1] {
			return true
		}
	}
	return false
}

func (s *SymSpell) updateSuggestions(suggestion string, cp *candidateProcessor) {
	suggestionCount := s.Words[suggestion]
	item := SuggestItem{Term: suggestion, Distance: cp.distance, Count: suggestionCount}

	if len(cp.suggestions) > 0 {
		if shouldContinue := s.updateBestSuggestion(cp, suggestionCount, item); shouldContinue {
			return
		}
	}
	// Update maxEditDistance2 if verbosity is not ALL
	if cp.verbosity != verbositypkg.All {
		cp.maxEditDistance2 = cp.distance
	}
	cp.suggestions = append(cp.suggestions, SuggestItem{Term: suggestion, Distance: cp.distance, Count: s.Words[suggestion]})
	return
}

func (s *SymSpell) updateBestSuggestion(cp *candidateProcessor, suggestionCount int, item SuggestItem) bool {
	if cp.verbosity == verbositypkg.Closest {
		// Keep only the closest suggestions
		if cp.distance < cp.maxEditDistance2 {
			cp.suggestions = []SuggestItem{}
		}
	} else if cp.verbosity == verbositypkg.Top {
		// Keep the top suggestion based on count or distance
		if cp.distance < cp.maxEditDistance2 || suggestionCount > cp.suggestions[0].Count {
			cp.maxEditDistance2 = cp.distance
			cp.suggestions[0] = item
		}
		return true
	}
	return false
}

func (s *SymSpell) addEditDistance(candidateRunes []rune, cp *candidateProcessor) {
	for i := 0; i < len(candidateRunes); i++ {
		deleteItem := string(candidateRunes[:i]) + string(candidateRunes[i+1:])
		if !cp.consideredDeletes[deleteItem] {
			cp.consideredDeletes[deleteItem] = true
			cp.candidates = append(cp.candidates, deleteItem)
		}
	}
}

func (s *SymSpell) distanceCompare(a, b string, maxDistance int) int {
	distance := s.distanceComparer.Distance(a, b)

	// Check if the distance exceeds the maxDistance
	if distance > maxDistance {
		return -1
	}
	return distance
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

type candidateProcessor struct {
	candidates            []string
	consideredDeletes     map[string]bool
	consideredSuggestions map[string]bool
	maxEditDistance2      int
	candidatePointer      int
	verbosity             verbositypkg.Verbosity
	phraseLen             int
	phraseRunes           []rune
	candidateLen          int
	distance              int
	minDistance           int
	suggestions           []SuggestItem
	suggestionRunes       []rune
	suggestionLen         int
	lenDiff               int
}

func newCandidateProcessor(maxEditDistance int, verbosity verbositypkg.Verbosity, phrase string) candidateProcessor {
	return candidateProcessor{
		candidates:            make([]string, 0),
		consideredDeletes:     make(map[string]bool),
		consideredSuggestions: make(map[string]bool),
		maxEditDistance2:      maxEditDistance,
		candidatePointer:      0,
		verbosity:             verbosity,
		phraseLen:             len([]rune(phrase)),
		phraseRunes:           []rune(phrase),
		candidateLen:          0,
		distance:              0,
		minDistance:           0,
		suggestions:           []SuggestItem{},
		lenDiff:               0,
	}
}

func (c *candidateProcessor) resetDistance() {
	c.distance, c.minDistance = 0, 0
}

func (c *candidateProcessor) updateSuggestion(suggestion string) {
	c.suggestionRunes = []rune(suggestion)
	c.suggestionLen = len(c.suggestionRunes)
}

func (c *candidateProcessor) sortCandidate() {
	if len(c.suggestions) > 1 {
		sort.Slice(c.suggestions, func(i, j int) bool {
			if c.suggestions[i].Distance == c.suggestions[j].Distance {
				return c.suggestions[i].Count > c.suggestions[j].Count
			}
			return c.suggestions[i].Distance < c.suggestions[j].Distance
		})
	}
}
