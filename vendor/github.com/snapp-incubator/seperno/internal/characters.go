package internal

// Persian and Arabic Characters
const (
	// Persian Numbers
	faD0 = '\u06F0'
	faD1 = '\u06F1'
	faD2 = '\u06F2'
	faD3 = '\u06F3'
	faD4 = '\u06F4'
	faD5 = '\u06F5'
	faD6 = '\u06F6'
	faD7 = '\u06F7'
	faD8 = '\u06F8'
	faD9 = '\u06F9'

	// Arabic Numbers
	arD0 = '\u0660'
	arD1 = '\u0661'
	arD2 = '\u0662'
	arD3 = '\u0663'
	arD4 = '\u0664'
	arD5 = '\u0665'
	arD6 = '\u0666'
	arD7 = '\u0667'
	arD8 = '\u0668'
	arD9 = '\u0669'

	// English Numbers
	enD0 = '\u0030'
	enD1 = '\u0031'
	enD2 = '\u0032'
	enD3 = '\u0033'
	enD4 = '\u0034'
	enD5 = '\u0035'
	enD6 = '\u0036'
	enD7 = '\u0037'
	enD8 = '\u0038'
	enD9 = '\u0039'

	//combining characters
	tilde                      rune = 771  // COMBINING tilde
	leftHalfRingBelow          rune = 796  // COMBINING left half ring below
	dotBelow                   rune = 803  // COMBINING dot below
	diaeresisBelow             rune = 804  // COMBINING diaeresis below
	ringBelow                  rune = 805  // COMBINING ring below
	ogonek                     rune = 808  // COMBINING ogonek
	verticalLineBelow          rune = 809  // COMBINING vertical line below
	breveBelow                 rune = 814  // COMBINING breve below
	invertedBreveBelow         rune = 815  // COMBINING inverted breve below
	longStrokeOverlay          rune = 822  // COMBINING long stroke overlay
	fermata                    rune = 850  // COMBINING fermata
	doubleBreveBelow           rune = 860  // COMBINING double breve below
	doubleRightwardsArrowBelow rune = 866  // COMBINING double rightwards arrow bellow
	plusSignBelow              rune = 799  // COMBINING plus sign below
	lowLine                    rune = 818  // COMBINING LOW LINE
	DIAERESIS                  rune = 776  // COMBINING DIAERESIS
	seagullBelow               rune = 828  // COMBINING SEAGULL BELOW
	leftAngleAbove             rune = 794  // COMBINING LEFT ANGLE ABOVE
	acuteAccent                rune = 769  // COMBINING ACUTE ACCENT
	upTackBelow                rune = 797  // COMBINING UP TACK BELOW
	candrabindu                rune = 784  // COMBINING candrabindu
	caronBelow                 rune = 812  // COMBINING caron BELOW
	snakeBelow                 rune = 7618 // COMBINING SNAKE BELOW
	dotAbove                   rune = 775  // COMBINING DOT ABOVE
	doubleMacronBelow          rune = 863  // COMBINING DOUBLE macron BELOW
	zigzagAbove                rune = 859  // COMBINING ZIGZAG ABOVE
	graveAccent                rune = 768  // COMBINING GRAVE ACCENT
	upwardsArrowBelow          rune = 846  // COMBINING UPWARDS ARROW BELOW
	tildeBelow                 rune = 816  // COMBINING tilde BELOW
	turnedCommaAbove           rune = 786  // COMBINING TURNED COMMA ABOVE
	caron                      rune = 780  // COMBINING caron
	overline                   rune = 773  // COMBINING overline
	graphemeJoiner             rune = 847  // COMBINING GRAPHEME JOINER
	macron                     rune = 772  // COMBINING macron
	enclosingDiamond           rune = 8415 // COMBINING ENCLOSING DIAMOND
	circumflexAccentBelow      rune = 813  // COMBINING CIRCUMFLEX ACCENT BELOW
	asteriskBelow              rune = 857  // COMBINING ASTERISK BELOW
	doubleBreve                rune = 861  // COMBINING DOUBLE BREVE
	palatalizedHookBelow       rune = 801  // COMBINING PALATALIZED HOOK BELOW
	longSolidusOverlay         rune = 824  // COMBINING LONG SOLIDUS OVERLAY
	doubleMacron               rune = 862  // COMBINING DOUBLE macron
	macronAcute                rune = 7620 // COMBINING macron-ACUTE
	graveMacron                rune = 7621 // COMBINING GRAVE-macron
	verticalTilde              rune = 830  // COMBINING VERTICAL tilde
	equalSignBelow             rune = 839  // COMBINING EQUALS SIGN BELOW
	latinSmallLetterX          rune = 879  // COMBINING LATIN SMALL LETTER X
	clockwiseRingOverlay       rune = 8409 // COMBINING CLOCKWISE RING OVERLAY
	anticlockwiseRingOverlay   rune = 8410 // COMBINING ANTICLOCKWISE RING OVERLAY
	verticalLineAbove          rune = 781  // COMBINING VERTICAL LINE ABOVE
	invertedBreve              rune = 785  // COMBINING INVERTED BREVE
	commaBelow                 rune = 806  // COMBINING COMMA BELOW
	cedilla                    rune = 807  // COMBINING cedilla
	invertedBridgeBelow        rune = 826  // COMBINING INVERTED BRIDGE BELOW

	// Mosavet Characters
	fatheh              rune = 1614 // Arabic Fatha
	tashdid             rune = 1617 // Arabic Shadda
	sokun               rune = 1618 // Arabic Sokun
	zameh               rune = 1615 // Arabic Damma
	kasreh              rune = 1616 // Arabic Kasra
	tanvinZameh         rune = 1612 // Arabic Dammatan
	tanvinFatheh        rune = 1611 // Arabic Fathatan
	tanvinKasreh        rune = 1613 // Arabic Kasratan
	alefLittle          rune = 1648 // Arabic Alef With Hamza Above (Small Alef)
	arabicTatweel       rune = 'ـ'  // Arabic Tatweel
	arabicMaddahAbove   rune = 1619 // Arabic Maddah Above
	arabicSmallYeh      rune = 1766
	arabicSubscriptAlef rune = 1622 // Arabic Subscript Alef

	//Hamzeh
	persianHamza rune = 1569 //ARABIC LET­TER HAMZA

	//Arabic Yeh
	arabicYehWithThreeDotsAbove = '\u063F'
	arabicYehWithInvertedV      = '\u063D'
	rtl                         = '\u200F'
	ltr                         = '\u200E'
	cc                          = '\uFFFF' //Convert Characters: replace all non-ASCII characters in the file with an escape sequence in the form of \uFFFF

	// Spaces and Special Characters
	spaceZeroWidthNonJoiner = '\u200c'
	space                   = '\u0020'
	noBreakSpace            = '\u00A0'
	zeroWidthNoBreakSpace   = '\uFEFF'
	zeroWidthJoiner         = '\u200D'
	zeroWidthSpace          = '\u200B'

	// Example Additional Constants
	nullChar = '\u0000'

	emptyString = ""       // Empty string
	nullString  = "\u0000" // Null character
	dot         = "."      // dot
	NewLine     = "\n"     // Newline

	vs0  rune = '\uFE00' // Variation Selector 0
	vs1  rune = '\uFE01' // Variation Selector 1
	vs2  rune = '\uFE02' // Variation Selector 2
	vs3  rune = '\uFE03' // Variation Selector 3
	vs4  rune = '\uFE04' // Variation Selector 4
	vs5  rune = '\uFE05' // Variation Selector 5
	vs6  rune = '\uFE06' // Variation Selector 6
	vs7  rune = '\uFE07' // Variation Selector 7
	vs8  rune = '\uFE08' // Variation Selector 8
	vs9  rune = '\uFE09' // Variation Selector 9
	vs10 rune = '\uFE0A' // Variation Selector 10
	vs11 rune = '\uFE0B' // Variation Selector 11
	vs12 rune = '\uFE0C' // Variation Selector 12
	vs13 rune = '\uFE0D' // Variation Selector 13
	vs14 rune = '\uFE0E' // Variation Selector 14
	vs15 rune = '\uFE0F' // Variation Selector 15
)

// punctuations slice
var punctuations = []rune{
	'.', '·', '؟', '?', '!', '‼', '⁉',
	'(', ')', '،', '٬', '٫', ':', ';',
	'{', '}', '[', ']', '"', '\'', '/', '\\',
	'+', '=', '%', '>', '<', '-', '_', '~',
}

// EndOfLineCharacters slice
var endOfLinesChar = []rune{
	'.', '؟', '?', '!', '‼', '⁉',
}

// basicCharacters represents an array of important Persian and Arabic characters for normalization
var basicCharacters = []rune{
	' ',  // Placeholder for space normalization
	'ی',  // Persian Yeh
	'ا',  // Alef
	'و',  // Waw
	'ک',  // Keheh
	'ه',  // Heh
	'پ',  // Peh
	'چ',  // Cheh
	'ژ',  // Jeh
	'گ',  // Gaf
	'ب',  // Beh
	'ت',  // Teh
	'ث',  // Theh
	'ج',  // Jeem
	'ح',  // Hah
	'خ',  // Khah
	'ذ',  // Thal
	'د',  // Dal
	'ر',  // Reh
	'ز',  // Zain
	'س',  // Seen
	'ش',  // Sheen
	'ص',  // Sad
	'ض',  // Dad
	'ط',  // Tah
	'ظ',  // Zah
	'ع',  // Ain
	'غ',  // Ghain
	'ق',  // Qaf
	'ل',  // Lam
	'م',  // Meem
	'ن',  // Noon
	'ف',  // Feh
	faD0, // Persian Number 0
	faD1, // Persian Number 1
	faD2, // Persian Number 2
	faD3, // Persian Number 3
	faD4, // Persian Number 4
	faD5, // Persian Number 5
	faD6, // Persian Number 6
	faD7, // Persian Number 7
	faD8, // Persian Number 8
	faD9, // Persian Number 9
}
