package internal

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"

	verbositypkg "github.com/snapp-incubator/symspell/internal/verbosity"
)

func parseWords(phrase string, preserveCase bool, splitBySpace bool) []string {
	if splitBySpace {
		if preserveCase {
			return strings.Split(phrase, " ")
		}
		return strings.Split(strings.ToLower(phrase), " ")
	}

	// Regex pattern to match words, including handling apostrophes
	var pattern string
	if preserveCase {
		pattern = `(\p{L}+['’]*\p{L}*)`
	} else {
		phrase = strings.ToLower(phrase)
		pattern = `(\p{L}+['’]*\p{L}*)`
	}

	re := regexp.MustCompile(pattern)
	return re.FindAllString(phrase, -1)
}

func (s *SymSpell) LookupCompound(phrase string, maxEditDistance int) *SuggestItem {
	terms1 := parseWords(phrase, false, false)
	cp := compoundProcessor{
		suggestions:     make([]SuggestItem, 0),
		suggestionParts: make([]SuggestItem, 0),
		replacedWords:   make(map[string]SuggestItem),
		isLastCombi:     false,
	}
	for i, _ := range terms1 {
		cp.terms1 = terms1[i]
		cp.suggestions, _ = s.Lookup(cp.terms1, verbositypkg.Top, maxEditDistance)
		// Combine adjacent terms
		if i > 0 && !cp.isLastCombi {
			cp.terms2 = terms1[i-1]
			suggestionsCombi, _ := s.Lookup(cp.terms2+cp.terms1, verbositypkg.Top, maxEditDistance)
			if len(suggestionsCombi) > 0 {
				best1 := cp.suggestionParts[len(cp.suggestionParts)-1]
				best2 := s.getBestSuggestion2(cp, maxEditDistance)
				skip := s.validateCombinationDistance(best1, best2, suggestionsCombi[0], &cp)
				if skip {
					continue
				}
			}
		}
		cp.isLastCombi = false

		// Handle terms with no perfect suggestion
		if len(cp.suggestions) > 0 && (cp.suggestions[0].Distance == 0 || len(cp.terms1) == 1) {
			cp.suggestionParts = append(cp.suggestionParts, cp.suggestions[0])
		} else {
			var suggestionSplitBest *SuggestItem
			if len(cp.suggestions) > 0 {
				suggestionSplitBest = &cp.suggestions[0]
			}
			if len(cp.terms1) > 1 {
				runes := []rune(cp.terms1)
				for j := 1; j < len(runes); j++ {
					suggestions1, suggestions2, isValid := s.getSuggestions(runes, j, maxEditDistance)
					if !isValid {
						continue
					}
					cp.suggestion1, cp.suggestion2 = *suggestions1, *suggestions2
					tmpDistance := s.distanceCompare(cp.terms1, cp.tempTerm(), maxEditDistance)
					if tmpDistance < 0 {
						tmpDistance = maxEditDistance + 1
					}
					// Update suggestionSplitBest based on distance
					if suggestionSplitBest != nil {
						if tmpDistance > suggestionSplitBest.Distance {
							continue
						}
						if tmpDistance < suggestionSplitBest.Distance {
							suggestionSplitBest = nil
						}
					}

					// Check for bigrams
					tmpCount := s.checkForBigram(&cp)

					splitSuggestion := SuggestItem{Term: cp.tempTerm(), Distance: tmpDistance, Count: tmpCount}
					if suggestionSplitBest == nil || splitSuggestion.Count > suggestionSplitBest.Count {
						suggestionSplitBest = &splitSuggestion
					}
				}
			}
			if suggestionSplitBest != nil {
				cp.updateReplaceWord(terms1[i], *suggestionSplitBest)
			} else {
				item := createWithProbability(terms1[i], maxEditDistance+1)
				cp.updateReplaceWord(terms1[i], item)
			}
		}
	}

	return s.finalizeAnswer(phrase, cp.suggestionParts)
}

func (s *SymSpell) getBestSuggestion2(cp compoundProcessor, maxEditDistance int) SuggestItem {
	var best2 SuggestItem
	if len(cp.suggestions) > 0 {
		best2 = cp.suggestions[0]
	} else {
		best2 = createWithProbability(cp.terms1, maxEditDistance+1)
	}
	return best2
}

func (s *SymSpell) validateCombinationDistance(best1 SuggestItem, best2 SuggestItem, suggestionsCombine SuggestItem, cp *compoundProcessor) bool {
	distance1 := best1.Distance + best2.Distance

	if distance1 >= 0 && suggestionsCombine.Distance+1 < distance1 ||
		(suggestionsCombine.Distance+1 == distance1 &&
			float64(suggestionsCombine.Count) > (float64(best1.Count)/s.N)*float64(best2.Count)) {
		suggestionsCombine.Distance += 1
		cp.suggestionParts[len(cp.suggestionParts)-1] = suggestionsCombine
		cp.replacedWords[cp.terms2] = suggestionsCombine
		cp.isLastCombi = true
		return true
	}
	return false
}

func (s *SymSpell) getSuggestions(runes []rune, split int, maxEditDistance int) (*SuggestItem, *SuggestItem, bool) {
	part1 := string(runes[:split])
	part2 := string(runes[split:])
	suggestions1, _ := s.Lookup(part1, verbositypkg.Top, maxEditDistance)
	suggestions2, _ := s.Lookup(part2, verbositypkg.Top, maxEditDistance)
	if len(suggestions1) == 0 || len(suggestions2) == 0 {
		return nil, nil, false
	}
	return &suggestions1[0], &suggestions2[0], true
}

func (s *SymSpell) checkForBigram(cp *compoundProcessor) int {
	var tmpCount int
	if count, exists := s.Bigrams[cp.tempTerm()]; exists {
		tmpCount = count

		// Update count if split corrections match
		if len(cp.suggestions) > 0 {
			bestSI := cp.suggestions[0]
			if cp.suggestion1.Term+cp.suggestion2.Term == cp.terms1 {
				tmpCount = int(math.Max(float64(tmpCount), float64(bestSI.Count+2)))
			} else if bestSI.Term == cp.suggestion1.Term || bestSI.Term == cp.suggestion2.Term {
				tmpCount = int(math.Max(float64(tmpCount), float64(bestSI.Count+1)))
			}
		} else if cp.suggestion1.Term+cp.suggestion2.Term == cp.terms1 {
			tmpCount = int(math.Max(
				float64(tmpCount),
				math.Max(
					float64(cp.suggestion1.Count),
					float64(cp.suggestion2.Count),
				)+2,
			))
		}
	} else {
		// Calculate Naive Bayes probability-based count
		tmpCount = int(math.Min(
			float64(s.BigramCountMin),
			float64(cp.suggestion1.Count)/s.N*float64(cp.suggestion2.Count),
		))
	}
	return tmpCount
}

func (c *compoundProcessor) updateReplaceWord(terms1 string, item SuggestItem) {
	c.suggestionParts = append(c.suggestionParts, item)
	c.replacedWords[terms1] = item
}

func (s *SymSpell) finalizeAnswer(phrase string, suggestionParts []SuggestItem) *SuggestItem {
	joinedTerm := ""
	joinedCount := s.N
	for _, item := range suggestionParts {
		joinedTerm += item.Term + " "
		joinedCount *= float64(item.Count) / s.N
	}
	joinedTerm = strings.TrimSpace(joinedTerm)

	return &SuggestItem{
		Term:     joinedTerm,
		Distance: s.distanceCompare(phrase, joinedTerm, math.MaxInt32),
		Count:    int(joinedCount),
	}
}

func createWithProbability(term string, distance int) SuggestItem {
	// Calculate Naive Bayes probability as the count
	probabilityCount := int(10 / math.Pow(10, float64(len(term))))

	return SuggestItem{
		Term:     term,
		Distance: distance,
		Count:    probabilityCount,
	}
}

// Helper function to safely parse integers
func tryParseInt64(value string) (int, bool) {
	parsed, err := strconv.Atoi(value)
	if err != nil {
		fmt.Println("[ERROR] parsing integer: ", err)
		return 0, false
	}
	return parsed, true
}

// Load bigram dictionary from a stream
func (s *SymSpell) LoadBigramDictionaryStream(corpusStream *os.File, termIndex, countIndex int, separator string) bool {
	scanner := bufio.NewScanner(corpusStream)

	// Define minimum parts depending on the separator
	minParts := 3
	if separator != "" {
		minParts = 2
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Split line by the separator
		var parts []string
		if separator == "" {
			parts = strings.Fields(line)
		} else {
			parts = strings.Split(line, separator)
		}

		if len(parts) < minParts {
			continue
		}

		// Parse count
		count, ok := tryParseInt64(parts[countIndex])
		if !ok {
			continue
		}

		// Create the key
		var key string
		if separator == "" {
			key = parts[termIndex] + " " + parts[termIndex+1]
		} else {
			key = parts[termIndex]
		}
		// Add to bigram dictionary
		s.Bigrams[key] = count

		// Update the minimum bigram count
		if count < s.BigramCountMin {
			s.BigramCountMin = count
		}
	}

	return true
}

func (s *SymSpell) LoadBigramDictionary(
	corpusPath string,
	termIndex, countIndex int,
	separator string,
) (bool, error) {
	if corpusPath == "" {
		return false, fmt.Errorf("corpus path cannot be empty")
	}
	// Check if the file exists
	file, err := os.Open(corpusPath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Use the stream-based loading function
	return s.LoadBigramDictionaryStream(file, termIndex, countIndex, separator), nil
}

type compoundProcessor struct {
	suggestions     []SuggestItem
	suggestionParts []SuggestItem
	replacedWords   map[string]SuggestItem
	terms1          string
	terms2          string
	suggestion1     SuggestItem
	suggestion2     SuggestItem
	isLastCombi     bool
}

func (c *compoundProcessor) tempTerm() string {
	return fmt.Sprintf("%s %s", c.suggestion1.Term, c.suggestion2.Term)
}
