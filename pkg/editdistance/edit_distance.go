package editdistance

type IEditDistance interface {
	Distance(a, b string) int
}

func NewEditDistance(Type string) *EditDistance {
	return &EditDistance{Type: Type}
}

const (
	DamerauLevenshtein = "DamerauLevenshtein"
)

type EditDistance struct {
	Type string
}

func (d EditDistance) Distance(a, b string) int {
	switch d.Type {
	case DamerauLevenshtein:
		return damerauLevenshteinDistance(a, b)
	}
	return 0
}

func damerauLevenshteinDistance(a, b string) int {
	m := len(a)
	n := len(b)

	// If either string is empty, the distance is the length of the other string.
	if m == 0 {
		return n
	}
	if n == 0 {
		return m
	}

	// Create a distance matrix (m+1)x(n+1)
	distance := make([][]int, m+1)
	for i := range distance {
		distance[i] = make([]int, n+1)
	}

	// Initialize the distance matrix
	for i := 0; i <= m; i++ {
		distance[i][0] = i
	}
	for j := 0; j <= n; j++ {
		distance[0][j] = j
	}

	// Fill the distance matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}

			// Calculate distances for insertion, deletion, and substitution
			distance[i][j] = min(
				distance[i-1][j]+1, // Deletion
				min(
					distance[i][j-1]+1,      // Insertion
					distance[i-1][j-1]+cost, // Substitution
				),
			)

			// Check for transposition
			if i > 1 && j > 1 && a[i-1] == b[j-2] && a[i-2] == b[j-1] {
				distance[i][j] = min(
					distance[i][j],
					distance[i-2][j-2]+cost, // Transposition
				)
			}
		}
	}

	// Return the bottom-right value, which is the Damerau-Levenshtein distance
	return distance[m][n]
}
