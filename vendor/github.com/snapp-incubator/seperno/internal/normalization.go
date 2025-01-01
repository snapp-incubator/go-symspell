package internal

import (
	"regexp"
	"strings"
)

type Normalize struct {
	convertHalfSpaceToSpace bool
	urlRemover              bool
	outerSpaceRemover       bool
	spaceCombiner           bool
	normalizePunctuations   bool
	endsWithEndOfLineChar   bool
}

func NewNormalizer(conf NormalizerOptions) *Normalize {
	return &Normalize{
		convertHalfSpaceToSpace: conf.ConvertHalfSpaceToSpace,
		urlRemover:              conf.URLRemover,
		outerSpaceRemover:       conf.OuterSpaceRemover,
		spaceCombiner:           conf.SpaceCombiner,
		normalizePunctuations:   conf.NormalizePunctuations,
		endsWithEndOfLineChar:   conf.EndsWithEndOfLineChar,
	}
}

// FindHalfSpace replaces a specific Unicode half-space with the given string representation
func (n Normalize) FindHalfSpace(input, halfSpace string) string {
	return strings.ReplaceAll(input, "\u200c", halfSpace)
}

// SpaceNormalizer normalizes spaces in the given input string based on the provided flag
func (n Normalize) spaceNormalizer(input string) string {
	// Replace specific HTML representation
	input = strings.ReplaceAll(input, "&zwnj;", " ")

	// Convert input to a rune slice
	inputRunes := []rune(input)

	// Iterate through the runes
	for i := 0; i < len(inputRunes); i++ {
		if n.convertHalfSpaceToSpace {
			switch inputRunes[i] {
			case spaceZeroWidthNonJoiner,
				space, noBreakSpace,
				zeroWidthNoBreakSpace,
				zeroWidthSpace:
				inputRunes[i] = ' ' // Replace with a space
				break
			case zeroWidthJoiner:
				inputRunes[i] = nullChar // Replace with nullChar
			}
		} else {
			switch inputRunes[i] {
			case space, noBreakSpace,
				zeroWidthNoBreakSpace,
				zeroWidthSpace:
				inputRunes[i] = ' ' // Replace with a space
				break

			case zeroWidthJoiner:
				inputRunes[i] = nullChar // Replace with nullChar
				break
			}
		}
	}

	// Create a new string, trim it, and replace nullChar
	output := strings.TrimSpace(string(inputRunes))
	output = strings.ReplaceAll(output, string(nullChar), "")
	output = strings.ToLower(output)

	return output
}

// BasicNormalizer normalizes a Persian input string.
// If input is nil, it returns nil. Applies specific transformations to the input string.
func (n Normalize) BasicNormalizer(input string) string {
	if input == "" {
		return ""
	}

	// Call SpecialYehNormalizer
	input = n.specialYehNormalizer(input)

	// Apply SpaceNormalizer
	input = n.spaceNormalizer(input)

	// Convert input to a rune slice for character-by-character processing
	inputRunes := []rune(input)

	for i := 0; i < len(inputRunes); i++ {
		switch inputRunes[i] {

		// "الف" group replacements break
		case 'أ', 'ﺍ', 'إ', 'ﺁ', 'ا', 'آ', 'ٵ', 'ﴽ', 'ﺂ', 'ﺄ', 'ﺈ', 'ﺎ', 'Ĩ', 'ٱ', 'ٲ', 'ﭐ', 'ﭑ', 'ﺇ':
			inputRunes[i] = basicCharacters[2]
			break

		// "ب" group replacements break
		case 'ﺒ', 'ٮ', 'ݕ', 'ﺏ', 'ﺐ', 'ﺑ':
			inputRunes[i] = basicCharacters[10]
			break

		// "ژ" group replacements break
		case 'ژ', 'ﮊ':
			inputRunes[i] = basicCharacters[8]
			break

		// "ی" group replacements break
		case 'ى', 'ئ', 'ي', 'ﯧ', 'ﯿ', 'ﻴ', 'ۍ', 'ې', 'ۑ', 'ﯤ', 'ﯼ', 'ﯽ', 'ﯾ',
			'ﺉ', 'ﺊ', 'ﻯ', 'ﻰ', 'ﻱ', 'ﻲ', 'ﻳ', 'ے':
			inputRunes[i] = basicCharacters[1]
			break
		// group replacements break
		case 'ك', 'ڪ', 'ﮑ', 'ﻜ', 'ػ', 'ګ', 'ڬ', 'ڭ', 'ڮ',
			'ݢ', 'ݣ', 'ݤ', 'ﻙ', 'ﻛ', 'ﮏ', 'ﮐ':
			inputRunes[i] = basicCharacters[4]
			break

		// "ه" group replacements break
		case 'ە', 'ہ', 'ﮭ', 'ھ', 'ۿ', 'ﮪ', 'ﮫ', 'ﮬ',
			'ﻪ', 'ﻫ', 'ﻬ', 'ﻩ':
			inputRunes[i] = basicCharacters[5]
			break

		// "م" group replacements break
		case 'ﻤ', '۾', 'ݥ', 'ﻡ', 'ﻢ', 'ﻣ':
			inputRunes[i] = basicCharacters[30]
			break

		// "ن" group replacements break
		case 'ﻨ', 'ڹ', 'ں', 'ڻ', 'ݧ', 'טּ', 'ﮟ', 'ﻥ', 'ﻦ', 'ﻧ':
			inputRunes[i] = basicCharacters[31]
			break

		// "و" group replacements break
		case 'ﻮ', 'ؤ', 'ٷ', 'ﯣ', 'ﺆ', 'ٶ', 'ۄ', 'ۅ', 'ۆ', 'ۇ', 'ۈ', 'ۉ', 'ۊ', 'ۋ', 'ۏ',
			1928, 'ﯗ', 'ﯙ', 'ﯚ', 'ﯛ', 'ﯝ', 'ﯡ', 'ﯢ', 'ﺅ', 'ﻭ':
			inputRunes[i] = basicCharacters[3]
			break

		// "ی" additional group replacements break
		case 'ٸ', 'ﺌ':
			inputRunes[i] = basicCharacters[1]
			break

		// "ة به ه" group replacements break
		case 'ة', 'ۀ', 'ﺔ', 'ۂ', 'ۃ', 'ﺓ':
			inputRunes[i] = basicCharacters[5]
			break

		// "پ" group replacements break
		case 'ﭙ', 'ݐ', 'ݒ', 'ﭖ', 'ﭗ', 'ﭘ':
			inputRunes[i] = basicCharacters[6]
			break

		// "چ" group replacements break
		case 'ﭽ':
			inputRunes[i] = basicCharacters[7]
			break

		// "گ" group replacements break
		case 'ﮕ', 'ڰ', 'ڲ', 'ڳ', 'ڴ', 'ﮒ', 'ﮓ', 'ﮔ', 'ﮗ', 'ﮘ', 'ﮝ':
			inputRunes[i] = basicCharacters[9]
			break

		// "ت" group replacements break
		case 'ٹ', 'ٺ', 'ټ', 'ٿ', 'ݓ', 'ﺘ', 'ﺕ', 'ﺖ', 'ﺗ',
			'ﭞ', 'ﭟ', 'ﭠ', 'ﭡ', 'ﭥ', 'ﭦ':
			inputRunes[i] = basicCharacters[11]
			break

		// "ث" group replacements break
		case 'ﺜ', 'ٽ', 'ݑ', 'ﺙ', 'ﺚ':
			inputRunes[i] = basicCharacters[12]
			break

		// "چ" group replacements break
		case 'ڇ', 'ڿ', 'ݘ', 'ﭺ', 'ﭼ':
			inputRunes[i] = basicCharacters[7]
			break

		// "ج" group replacements break
		case 'ﺠ', 'ڃ', 'ﺝ', 'ﺞ', 'ﺟ':
			inputRunes[i] = basicCharacters[13]
			break

		// "ح" group replacements break
		case 'ﺤ', 'ځ', 'ﺡ', 'ﺢ', 'ﺣ':
			inputRunes[i] = basicCharacters[14]
			break

		// "خ" group replacements break
		case 'ﺨ', 'ڂ', 'ݗ', 'ﺥ', 'ﺦ', 'ﺧ':
			inputRunes[i] = basicCharacters[15]
			break

		// "د" group replacements break
		case 'ﺪ', 'ڈ', 'ډ', 'ڊ', 'ڋ', 'ڌ', 'ڍ', 'ڐ', 'ۮ', 'ﮈ', 'ﺩ':
			inputRunes[i] = basicCharacters[17]
			break

		// "ذ" group replacements break
		case 'ﺬ':
			inputRunes[i] = basicCharacters[16]
			break

		// "ر" group replacements break
		case 'ﺮ', 'ڑ', 'ڒ', 'ړ', 'ڔ', 'ڕ', 'ږ', 'ۯ', 'ݛ', 'ݬ', 'ﮍ', 'ﺭ':
			inputRunes[i] = basicCharacters[18]
			break

		// "ز" group replacements break
		case 'ﺰ', 'ڗ', 'ݫ', 'ﺯ':
			inputRunes[i] = basicCharacters[19]
			break

		// "س" group replacements break
		case 'ﺴ', 'ښ', 'ڛ', 'ݭ', 'ﺱ', 'ﺲ', 'ﺳ':
			inputRunes[i] = basicCharacters[20]
			break

		// "ش" group replacements break
		case 'ﺸ', 'ڜ', 'ۺ', 'ݜ', 'ﺵ', 'ﺶ', 'ﺷ':
			inputRunes[i] = basicCharacters[21]
			break

		// "ص" group replacements break
		case 'ﺼ', 'ڝ', 'ﺹ', 'ﺺ', 'ﺻ':
			inputRunes[i] = basicCharacters[22]
			break

		// "ض" group replacements break
		case 'ﻀ', 'ۻ', 'ﺽ', 'ﺾ', 'ﺿ':
			inputRunes[i] = basicCharacters[23]
			break

		// "ط" group replacements break
		case 'ﻄ', 'ﻁ', 'ﻂ', 'ﻃ':
			inputRunes[i] = basicCharacters[24]
			break

		// "ظ" group replacements break
		case 'ﻈ', 'ڟ', 'ﻅ', 'ﻆ', 'ﻇ':
			inputRunes[i] = basicCharacters[25]
			break

		// "ع" group replacements break
		case 'ﻌ', '؏', 'ڠ', 'ﻉ', 'ﻊ', 'ﻋ':
			inputRunes[i] = basicCharacters[26]
			break

		// "غ" group replacements break
		case 'ﻐ', 'ۼ', 'ݞ', 'ݟ', 'ﻍ', 'ﻎ', 'ﻏ':
			inputRunes[i] = basicCharacters[27]
			break

		// "ق" group replacements break
		case 'ﻘ', 'ڦ', 'ڧ', 'ڨ', 'ﻕ', 'ﻖ', 'ﻗ':
			inputRunes[i] = basicCharacters[28]
			break

		// "ف" group replacements break
		case '؋', 'ف', 'ڢ', 'ڣ', 'ڤ', 'ڥ', 5317, 'ﻔ', 'ﻓ', 'ﻑ', 'ﻒ':
			inputRunes[i] = basicCharacters[32]
			break

		// "ل" group replacements break
		case 'ﻠ', 'ڵ', 'ڶ', 'ڷ', 'ڸ', 'ݪ', 'ﻝ', 'ﻞ', 'ﻟ':
			inputRunes[i] = basicCharacters[29]
			break
			// Replace specific characters with null char break
		case sokun, 65150, 65151, fatheh, 65142,
			65143, zameh, 65144, 65145, kasreh, 65147, tashdid, 64607, 64608, 64609,
			tanvinFatheh, tanvinZameh, tanvinKasreh, alefLittle, persianHamza, 1620,
			1652, 1789, 64420, 64421, 65163, 65152, rtl, ltr, cc, arabicSubscriptAlef,
			tilde, leftHalfRingBelow, dotBelow, diaeresisBelow, ringBelow,
			ogonek, verticalLineBelow, breveBelow, invertedBreveBelow,
			longStrokeOverlay, fermata, doubleBreveBelow,
			doubleRightwardsArrowBelow, plusSignBelow, lowLine,
			DIAERESIS, seagullBelow, leftAngleAbove, acuteAccent,
			upTackBelow, candrabindu, caronBelow, snakeBelow, dotAbove,
			doubleMacronBelow, zigzagAbove, graveAccent, upwardsArrowBelow,
			tildeBelow, turnedCommaAbove, caron, overline, graphemeJoiner,
			macron, enclosingDiamond, circumflexAccentBelow, asteriskBelow,
			doubleBreve, palatalizedHookBelow, longSolidusOverlay, doubleMacron,
			graveMacron, verticalTilde, equalSignBelow, latinSmallLetterX,
			clockwiseRingOverlay, anticlockwiseRingOverlay, verticalLineAbove,
			invertedBreve, commaBelow, cedilla, invertedBridgeBelow,
			macronAcute, 64610, 65148, 65149:
			inputRunes[i] = nullChar
			break

		// "ی" group replacements break
		case arabicYehWithThreeDotsAbove, arabicYehWithInvertedV:
			inputRunes[i] = basicCharacters[1]
			break
			// Normalize Persian and other numeral variations to English numerals break
		case faD0, 3664: // Persian 0 and Thai 0 (๐)
			inputRunes[i] = enD0
			break
		case faD1, 3665, 185, 8321, 65297, 3793: // Persian 1, Thai 1 (๑), superscript 1 (¹), subscript 1 (₁), full-width 1, and others
			inputRunes[i] = enD1
			break
		case faD2, 8322, 178, 2536, 3666, 1399, 65298: // Persian 2, subscript 2 (₂), superscript 2 (²), full-width 2, and others
			inputRunes[i] = enD2

		// Normalize Persian and other numeral variations to English numerals break
		case faD3, 3667, 2537, 65299: // Persian 3, Thai 3 (๓), and others
			inputRunes[i] = enD3
			break
		case faD4, 3668, 8308, 3178, 65300: // Persian 4, Thai 4 (๔), superscript 4 (⁴), and others
			inputRunes[i] = enD4
			break
		case faD5, 65301: // Persian 5, full-width 5 (５)
			inputRunes[i] = enD5
			break
		case faD6, 3670: // Persian 6, Thai 6 (๖)
			inputRunes[i] = enD6
			// Normalize Persian and other numeral variations to English numerals break
		case faD7, 2925, 2797, 2669, 2413, 3671: // Persian 7, Devanagari 7 (७), Thai 7 (๗), and others
			inputRunes[i] = enD7
			break
		case faD8, 8328: // Persian 8, subscript 8 (₈)
			inputRunes[i] = enD8
			break
		case faD9, 3881, 2543, 6121: // Persian 9, and other numeral variations
			inputRunes[i] = enD9
			// Normalize Arabic numerals to English numerals break
		case arD0:
			inputRunes[i] = enD0
			break
		case arD1:
			inputRunes[i] = enD1
			break
		case arD2:
			inputRunes[i] = enD2
			break
		case arD3:
			inputRunes[i] = enD3
			break
		case arD4:
			inputRunes[i] = enD4
			break
		case arD5:
			inputRunes[i] = enD5
			break
		case arD6:
			inputRunes[i] = enD6
			break
		case arD7:
			inputRunes[i] = enD7
			break
		case arD8:
			inputRunes[i] = enD8
			break
		case arD9:
			inputRunes[i] = enD9

		// Replace punctuation marks with Persian equivalentsbreak
		case '?':
			inputRunes[i] = '؟'
			break
		case '%':
			inputRunes[i] = '٪'
			break
		case ';':
			inputRunes[i] = '؛'
			break
		case '：':
			inputRunes[i] = ':'
			break

		// Replace commas with Persian commasbreak
		case ',':
			inputRunes[i] = '،'
			break
		case '٬':
			inputRunes[i] = '،'
			break
		// Normalize various dash-like characters to underscorebreak
		case '–', '˗', '־', '­', '━', '—', '─', '_', '➖', '-', 'ـ':
			inputRunes[i] = '_'
			break
		// Replace various dash-like characters with ellipsisbreak
		case '┅', '┄', '┈':
			inputRunes[i] = '…'
			break
		}

	}
	// Convert the rune slice back to a string, trim spaces, replace null strings, and convert to lowercase
	s := strings.TrimSpace(string(inputRunes))
	s = strings.ReplaceAll(s, dot, nullString)
	s = strings.ReplaceAll(s, NewLine, nullString)
	s = strings.ReplaceAll(s, nullString, emptyString)
	s = strings.ToLower(s)
	// Convert runes back to a string
	stringInput := string(inputRunes)
	if n.spaceCombiner {
		stringInput = replaceMultiSpace(stringInput)
	}
	if n.urlRemover {
		stringInput = removeURLs(stringInput)
	}
	if n.normalizePunctuations {
		stringInput = normalizePunctuations(stringInput)
	}
	if n.endsWithEndOfLineChar {
		stringInput = normalizeEndsWithEndOfLineChar(stringInput)
	}
	if n.outerSpaceRemover { // should be last normalization step
		stringInput = removeOuterSpace(stringInput)
	}
	return stringInput
}

func (n Normalize) VariationSelectorsRemover(input []string) []string {
	output := make([]string, len(input))
	for i, str := range input {
		output[i] = removeVariationSelectors(str)
	}
	return output
}

func removeVariationSelectors(input string) string {
	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		switch runes[i] {
		case vs0, vs1, vs2, vs3, vs4, vs5, vs6, vs7, vs8, vs9,
			vs10, vs11, vs12, vs13, vs14, vs15:
			runes[i] = ' ' // Replace with space
		}
	}
	return string(runes)
}

func (n Normalize) specialYehNormalizer(input string) string {
	var builder strings.Builder
	for _, c := range input {
		switch c {
		case 'ے', 'ﮮ', 'ﮯ', 'ۓ', 'ﮱ': // Special "yeh" characters
			builder.WriteRune(basicCharacters[1]) // StandarD "ی"
			builder.WriteRune(basicCharacters[0]) // Spacebreak
			break
		case 'ﻩ', 'ﮦ': // Special "heh" characters
			builder.WriteRune(basicCharacters[5]) // StandarD "ه"
			builder.WriteRune(basicCharacters[0]) // space
			break
		default:
			builder.WriteRune(c)
		}
	}
	return builder.String()
}

// BasicNormalizerArray Normalize each string in an array with attention to Persian language.
func (n Normalize) BasicNormalizerArray(input []string) []string {
	for i := range input {
		input[i] = n.BasicNormalizer(input[i])
	}
	return input
}

// BasicNormalizerSlice Normalize each string in a slice (ArrayList equivalent in Go) with attention to Persian language.
func (n Normalize) BasicNormalizerSlice(input []string) []string {
	result := make([]string, len(input))
	for i, str := range input {
		result[i] = n.BasicNormalizer(str)
	}
	return result
}

func replaceMultiSpace(input string) string {
	re := regexp.MustCompile(`\s+`) // Matches one or more whitespace characters
	return re.ReplaceAllString(input, " ")
}

func removeOuterSpace(input string) string {
	// Create a regex pattern to remove spaces from start and enD
	re := regexp.MustCompile(`^\s+|\s+$`)
	// Remove leading and trailing spaces
	return re.ReplaceAllString(input, "")
}

func removeURLs(input string) string {
	// Regex to match URLs
	re := regexp.MustCompile(`https?://[^\s]+`)
	// Replace all URLs with an empty string
	return re.ReplaceAllString(input, "")
}

func normalizePunctuations(input string) string {
	var builder strings.Builder

	// Convert input string to runes for Unicode-safe processing
	for _, r := range input {
		if containsRune(punctuations, r) {
			// Replace punctuations with space (or remove by not appending)
			builder.WriteRune(' ') // Replace with space
		} else {
			builder.WriteRune(r) // Keep other characters
		}
	}

	return builder.String()
}

func normalizeEndsWithEndOfLineChar(input string) string {
	runes := []rune(input)
	if len(runes) == 0 {
		return input // Empty string remains unchanged
	}

	lastChar := runes[len(runes)-1]
	if containsRune(endOfLinesChar, lastChar) {
		// Remove the last character by slicing the rune slice
		return string(runes[:len(runes)-1])
	}
	return input // No end-of-line character to remove
}

func containsRune(slice []rune, r rune) bool {
	for _, s := range slice {
		if s == r {
			return true
		}
	}
	return false
}
