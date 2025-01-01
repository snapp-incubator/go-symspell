package verbosity

// Verbosity controls the quantity/closeness of returned suggestions.
type Verbosity int

const (
	// Top returns the closest match.
	Top Verbosity = iota
	// Closest returns all matches within max edit distance.
	Closest
	// All returns all possible suggestions.
	All
)
