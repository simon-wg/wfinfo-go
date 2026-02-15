package internal

const (
	matchScore    = 2
	mismatchScore = -1
	gapPenalty    = -1
)

func smithWaterman(s string, ss []string) string {
	bestScore := -1
	bestMatch := ""

	// Convert query to runes for unicode safety
	queryRunes := []rune(s)

	for _, candidate := range ss {
		candidateRunes := []rune(candidate)
		score := calculateScore(queryRunes, candidateRunes)
		if score > bestScore {
			bestScore = score
			bestMatch = candidate
		}
	}

	return bestMatch
}

func calculateScore(s1, s2 []rune) int {
	rows := len(s1) + 1
	cols := len(s2) + 1

	// Initialize matrix with 0s
	// H[i][j] holds the score of the optimal local alignment ending at s1[i-1], s2[j-1]
	matrix := make([][]int, rows)
	for i := range matrix {
		matrix[i] = make([]int, cols)
	}

	maxScore := 0

	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			// Calculate match or mismatch
			scoreDir := mismatchScore
			if s1[i-1] == s2[j-1] {
				scoreDir = matchScore
			}

			// Calculate potential values
			diagonal := matrix[i-1][j-1] + scoreDir
			up := matrix[i-1][j] + gapPenalty
			left := matrix[i][j-1] + gapPenalty

			// Smith-Waterman rule: value is max(0, diag, up, left)
			value := 0
			if diagonal > value {
				value = diagonal
			}
			if up > value {
				value = up
			}
			if left > value {
				value = left
			}

			matrix[i][j] = value

			// Track the maximum score found anywhere in the matrix
			if value > maxScore {
				maxScore = value
			}
		}
	}

	return maxScore
}
