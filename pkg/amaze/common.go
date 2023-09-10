package amaze

import (
	"fmt"
	"math/rand"
)

// Returns a struct of anonymous Direction types with row and col accessible.
func generateDirections(dist int) []struct{ row, col int } {
	return []struct{ row, col int }{{-dist, 0}, {dist, 0}, {0, -dist}, {0, dist}}
}

// Returns an odd series of integers within a given start and end point.
func getOddSeries(min int, max int) []int {
	var oddIndices []int
	for i := min; i <= max; i++ {
		if i%2 != 0 {
			oddIndices = append(oddIndices, i)
		}
	}
	return oddIndices
}

// Get a random value from a slice of int.
func getRandom(s []int) int {
	return s[rand.Intn(len(s))]
}

func CursorTopLeft() {
	fmt.Print("\033[H")
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
	// <?> overkill?
	// cmd := exec.Command("clear") // Use "cls" on Windows
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}
