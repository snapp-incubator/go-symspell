package internal

import (
	"bufio"
	"fmt"
	verbositypkg "github.com/snapp-incubator/symspell/internal/verbosity"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
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

func (s *SymSpell) LookupCompound(phrase string, maxEditDistance int) []SuggestItem {
	terms1 := parseWords(phrase, true, true)
	var suggestions []SuggestItem
	var suggestionParts []SuggestItem
	isLastCombi := false
	for i, term := range terms1 {
		suggestions, _ = s.Lookup(terms1[i], verbositypkg.Top, maxEditDistance)
		// Combine adjacent terms
		if i > 0 && !isLastCombi {
			suggestionsCombi, _ := s.Lookup(terms1[i-1]+terms1[i], verbositypkg.Top, maxEditDistance)
			if len(suggestionsCombi) > 0 {
				best1 := suggestionParts[len(suggestionParts)-1]
				var best2 SuggestItem
				if len(suggestions) > 0 {
					best2 = suggestions[0]
				} else {
					best2 = createWithProbability(terms1[i], maxEditDistance+1)
				}
				distance1 := best1.Distance + best2.Distance

				if distance1 >= 0 && suggestionsCombi[0].Distance+1 < distance1 ||
					(suggestionsCombi[0].Distance+1 == distance1 &&
						float64(suggestionsCombi[0].Count) > (float64(best1.Count)/s.N)*float64(best2.Count)) {
					suggestionsCombi[0].Distance += 1
					suggestionParts[len(suggestionParts)-1] = suggestionsCombi[0]
					s.ReplacedWords[terms1[i-1]] = suggestionsCombi[0]
					isLastCombi = true
					continue
				}
			}
		}
		isLastCombi = false

		// Handle terms with no perfect suggestion
		if len(suggestions) > 0 && (suggestions[0].Distance == 0 || len(terms1[i]) == 1) {
			suggestionParts = append(suggestionParts, suggestions[0])
		} else {
			var suggestionSplitBest *SuggestItem
			if len(suggestions) > 0 {
				suggestionSplitBest = &suggestions[0]
			}
			if len(terms1[i]) > 1 {
				runes := []rune(term)
				for j := 1; j < len(runes); j++ {
					part1 := string(runes[:j])
					part2 := string(runes[j:])
					suggestions1, _ := s.Lookup(part1, verbositypkg.Top, maxEditDistance)
					suggestions2, _ := s.Lookup(part2, verbositypkg.Top, maxEditDistance)
					if len(suggestions1) == 0 || len(suggestions2) == 0 {
						continue
					}
					tmpTerm := suggestions1[0].Term + " " + suggestions2[0].Term
					tmpDistance := s.distanceCompare(term, tmpTerm, maxEditDistance)
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
					var tmpCount int
					if count, exists := s.Bigrams[tmpTerm]; exists {
						tmpCount = count

						// Update count if split corrections match
						if len(suggestions) > 0 {
							bestSI := suggestions[0]
							if suggestions1[0].Term+suggestions2[0].Term == term {
								tmpCount = int(math.Max(float64(tmpCount), float64(bestSI.Count+2)))
							} else if bestSI.Term == suggestions1[0].Term || bestSI.Term == suggestions2[0].Term {
								tmpCount = int(math.Max(float64(tmpCount), float64(bestSI.Count+1)))
							}
						} else if suggestions1[0].Term+suggestions2[0].Term == term {
							tmpCount = int(math.Max(
								float64(tmpCount),
								math.Max(
									float64(suggestions1[0].Count),
									float64(suggestions2[0].Count),
								)+2,
							))
						}
					} else {
						// Calculate Naive Bayes probability-based count
						tmpCount = int(math.Min(
							float64(s.BigramCountMin),
							float64(suggestions1[0].Count)/s.N*float64(suggestions2[0].Count),
						))
					}

					splitSuggestion := SuggestItem{Term: tmpTerm, Distance: tmpDistance, Count: tmpCount}
					if suggestionSplitBest == nil || splitSuggestion.Count > suggestionSplitBest.Count {
						suggestionSplitBest = &splitSuggestion
					}
				}
			}
			if suggestionSplitBest != nil {
				suggestionParts = append(suggestionParts, *suggestionSplitBest)
				s.ReplacedWords[terms1[i]] = *suggestionSplitBest
			} else {
				item := createWithProbability(terms1[i], maxEditDistance+1)
				suggestionParts = append(suggestionParts, item)
				s.ReplacedWords[terms1[i]] = item
			}
		}
	}

	// Combine final suggestions
	joinedTerm := ""
	joinedCount := s.N
	for _, item := range suggestionParts {
		joinedTerm += item.Term + " "
		joinedCount *= float64(item.Count) / s.N
	}
	joinedTerm = strings.TrimSpace(joinedTerm)
	finalSuggestion := SuggestItem{
		Term:     joinedTerm,
		Distance: s.distanceCompare(phrase, joinedTerm, math.MaxInt32),
		Count:    int(joinedCount),
	}

	return []SuggestItem{finalSuggestion}
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
		fmt.Println("Error parsing integer: ", err)
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
