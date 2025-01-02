package internal

import (
	"errors"
	verbositypkg "github.com/snapp-incubator/symspell/internal/verbosity"
	"math"
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

	var suggestions []SuggestItem
	phraseRunes := []rune(phrase)
	phraseLen := len(phraseRunes)

	// Early exit function
	earlyExit := func(sug []SuggestItem) []SuggestItem {
		return sug
	}

	// Early exit - word too big to match any words
	if phraseLen-maxEditDistance > s.maxLength {
		return earlyExit(suggestions), nil
	}

	// Quick look for exact match
	if count, found := s.Words[phrase]; found {
		suggestions = append(suggestions, SuggestItem{Term: phrase, Distance: 0, Count: count})
		if verbosity != verbositypkg.All {
			return earlyExit(suggestions), nil
		}
	}

	// Early termination for max edit distance == 0
	if maxEditDistance == 0 {
		return earlyExit(suggestions), nil
	}

	consideredDeletes := map[string]bool{}
	consideredSuggestions := map[string]bool{}
	consideredSuggestions[phrase] = true

	maxEditDistance2 := maxEditDistance
	candidatePointer := 0
	var candidates []string

	// Add original prefix
	phrasePrefixRunes := phraseRunes
	if phraseLen > s.PrefixLength {
		phrasePrefixRunes = phraseRunes[:s.PrefixLength]
	}
	phrasePrefix := string(phrasePrefixRunes)
	candidates = append(candidates, phrasePrefix)

	// Process candidates
	for candidatePointer < len(candidates) {
		candidate := candidates[candidatePointer]
		candidatePointer++
		candidateRunes := []rune(candidate)
		candidateLen := len(candidateRunes)
		lenDiff := phraseLen - candidateLen

		// Early termination: if candidate distance is already higher than
		// suggestion distance, then there are no better suggestions to be
		// expected
		if lenDiff > maxEditDistance2 {
			if verbosity == verbositypkg.All {
				// `max_edit_distance_2`` only updated when
				// verbosity != ALL. New candidates are generated from
				// deletes so it keeps getting shorter. This should never
				// be reached.
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
				suggestionRunes := []rune(suggestion)
				suggestionLen := len(suggestionRunes)
				if abs(suggestionLen-phraseLen) > maxEditDistance2 || suggestionLen < candidateLen ||
					(suggestionLen == candidateLen && suggestion != candidate) {
					continue
				}

				// True Damerau-Levenshtein Edit Distance: adjust distance,
				// if both distances>0. We allow simultaneous edits (deletes)
				// of max_edit_distance on on both the dictionary and the
				// phrase term. For replaces and adjacent transposes the
				// resulting edit distance stays <= max_edit_distance. For
				// inserts and deletes the resulting edit distance might
				// exceed max_edit_distance. To prevent suggestions of a
				// higher edit distance, we need to calculate the resulting
				// edit distance, if there are simultaneous edits on both
				// sides. Example: (bank==bnak and bank==bink, but bank!=kanb
				// and bank!=xban and bank!=baxn for max_edit_distance=1).
				// Two deletes on each side of a pair makes them all equal,
				// but the first two pairs have edit distance=1, the others
				// edit distance=2.
				distance := 0
				minDistance := 0
				if candidateLen == 0 {
					distance = max(phraseLen, suggestionLen)
					if distance > maxEditDistance2 || consideredSuggestions[suggestion] {
						continue
					}
				} else if suggestionLen == 1 {
					var distanceCalc = func() int {
						phraseRunesList := []rune(phrase)
						// Check if the first rune of suggestion exists in phrase
						if slices.Contains(phraseRunesList, suggestionRunes[0]) {
							return phraseLen - 1
						} else {
							return phraseLen
						}
					}
					distance = distanceCalc()
					if distance > maxEditDistance2 || consideredSuggestions[suggestion] {
						continue
					}
					//number of edits in prefix ==maxeditdistance AND no
					//identical suffix, then editdistance>max_edit_distance and
					//no need for Levenshtein calculation
					//(phraseLen >= prefixLength) &&
					//(suggestionLen >= prefixLength)
				} else {
					// handles the shortcircuit of min_distance assignment
					// when first boolean expression evaluates to False
					if s.PrefixLength-maxEditDistance == candidateLen {
						minDistance = min(phraseLen, suggestionLen) - s.PrefixLength
					} else {
						minDistance = 0
					}
					if s.PrefixLength-maxEditDistance == candidateLen {
						if minDistance > 1 &&
							string(phraseRunes[phraseLen+1-minDistance:]) != string(suggestionRunes[suggestionLen+1-minDistance:]) {
							continue
						}
						if minDistance > 0 &&
							phraseRunes[phraseLen-minDistance] != suggestionRunes[suggestionLen-minDistance] {
							if phraseRunes[phraseLen-minDistance-1] != suggestionRunes[suggestionLen-minDistance] ||
								phraseRunes[phraseLen-minDistance] != suggestionRunes[suggestionLen-minDistance-1] {
								continue
							}
						}
					}
					// delete_in_suggestion_prefix is somewhat expensive, and
					// only pays off when verbosity is TOP or CLOSEST
					if consideredSuggestions[suggestion] {
						continue
					}
					consideredSuggestions[suggestion] = true
					distance = s.distanceCompare(phrase, suggestion, maxEditDistance2)
					if distance < 0 {
						continue
					}
				}
				// do not process higher distances than those already found,
				// if verbosity<ALL (note: max_edit_distance_2 will always
				// equal max_edit_distance when Verbosity.ALL)
				if distance <= maxEditDistance2 {
					suggestionCount := s.Words[suggestion]
					item := SuggestItem{Term: suggestion, Distance: distance, Count: suggestionCount}

					if len(suggestions) > 0 {
						if verbosity == verbositypkg.Closest {
							// Keep only the closest suggestions
							if distance < maxEditDistance2 {
								suggestions = []SuggestItem{}
							}
						} else if verbosity == verbositypkg.Top {
							// Keep the top suggestion based on count or distance
							if distance < maxEditDistance2 || suggestionCount > suggestions[0].Count {
								maxEditDistance2 = distance
								suggestions[0] = item
							}
							continue
						}
					}
					// Update maxEditDistance2 if verbosity is not ALL
					if verbosity != verbositypkg.All {
						maxEditDistance2 = distance
					}
					suggestions = append(suggestions, SuggestItem{Term: suggestion, Distance: distance, Count: s.Words[suggestion]})
				}
			}
		}
		// Add Edits: derive edits (deletes) from candidate (phrase) and add
		// them to candidates list. this is a recursive process until the
		// maximum edit distance has been reached
		if lenDiff <= maxEditDistance && candidateLen <= s.PrefixLength {
			if verbosity != verbositypkg.All && lenDiff >= maxEditDistance2 {
				continue
			}
			for i := 0; i < len(candidateRunes); i++ {
				deleteItem := string(candidateRunes[:i]) + string(candidateRunes[i+1:])
				if !consideredDeletes[deleteItem] {
					consideredDeletes[deleteItem] = true
					candidates = append(candidates, deleteItem)
				}
			}
		}

	}

	// Sort suggestions
	if len(suggestions) > 1 {
		sort.Slice(suggestions, func(i, j int) bool {
			if suggestions[i].Distance == suggestions[j].Distance {
				return suggestions[i].Count > suggestions[j].Count
			}
			return suggestions[i].Distance < suggestions[j].Distance
		})
	}

	earlyExit(suggestions)
	return suggestions, nil
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

func incrementCount(count, countPrevious int) int {
	// Ensure the count does not exceed the maximum value for int64
	if math.MaxInt64-countPrevious > count {
		return countPrevious + count
	}
	return math.MaxInt64
}
