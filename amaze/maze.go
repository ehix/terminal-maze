package amaze

import (
	"fmt"
	"math/rand"
)

var tiles TileSet

type Maze struct {
	maxRows int
	maxCols int
	grid    [][]string
	row     int
	col     int
	start   []int
	end     []int
	trail   []struct{ row, col int } // Store the player's trail
}

// NewMaze creates a new maze
func NewMaze(maxRows int, maxCols int) *Maze {
	tiles = SetRandomTiles()
	// tiles = tileSets["frog"]
	m := &Maze{
		maxRows: maxRows,
		maxCols: maxCols,
		grid:    make([][]string, maxRows),
		row:     0,
		col:     1,
	}

	// Turn the grid of empty slices into walls
	for i := range m.grid {
		m.grid[i] = make([]string, maxCols)
		for j := range m.grid[i] {
			m.grid[i][j] = tiles.wall
		}
	}
	return m
}

// SetStartExit sets the start and exit positions
func (m *Maze) SetStartExit() {
	var randomCol []int
	for i := 0; i < 2; i++ {
		index := getRandom(getOddSeries(1, m.maxCols-2))
		randomCol = append(randomCol, index)
	}
	startCol, exitCol := randomCol[0], randomCol[1]
	// Add start and exit emoji
	m.grid[0][startCol] = tiles.start
	m.grid[m.maxRows-1][exitCol] = tiles.exit

	// Add to for Generate function
	m.row = 0
	m.col = startCol
	m.start = []int{0, startCol}
	m.end = []int{m.maxRows - 1, exitCol}
}

// IsGameOver checks if the game is over
func (m *Maze) IsGameOver() bool {
	return m.grid[m.row][m.col] == tiles.exit
}

// Generate generates the maze using Recursive Backtracking Algorithm
func (m *Maze) Generate() {
	// SetStartExit is called after this, that's why row 0 and (I'm assuming) col 7 aren't changed in the loop?
	// Maybe set the start and exit before hand and then be able to use them here?

	// var endReached = false
	// Changes m.grid in place, by removing walls and replacing with spaces.
	stack := []struct{ row, col int }{{1, m.col}} // <?> could just find the start here, meaning we could remove m.start?
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		row, col := current.row, current.col
		m.grid[row][col] = tiles.empty
		var neighbors []struct{ row, col int }
		directions := generateDirections(2)

		for _, d := range directions {
			nextRow, nextCol := row+d.row, col+d.col
			// if it's in boundaries and a wall.
			if nextRow > 0 && nextRow < (m.maxRows-1) && nextCol > 0 && nextCol < (m.maxCols-1) && m.grid[nextRow][nextCol] == tiles.wall {
				neighbors = append(neighbors, struct{ row, col int }{nextRow, nextCol})
			}
		}

		if len(neighbors) > 0 {
			randomIndex := rand.Intn(len(neighbors)) // choose a random neighbor
			randomNeighbor := neighbors[randomIndex]
			nRow, nCol := randomNeighbor.row, randomNeighbor.col // get the indices

			m.grid[nRow][nCol] = tiles.empty                          // make the random cell empty
			m.grid[row+(nRow-row)/2][col+(nCol-col)/2] = tiles.empty  // fills in adjecent cell on the path taken
			stack = append(stack, struct{ row, col int }{nRow, nCol}) // append the route taken?
		} else {
			stack = stack[:len(stack)-1] // remove from the stack
		}
	}
	for i := 0; i < 3; i++ {
		// Add more empty cells with a certain probability
		if rand.Float64() < 0.7 { // make this scale with every game?
			m.addRandomEmpty()
		}
	}
}

func generateDirections(dist int) []struct{ row, col int } {
	return []struct{ row, col int }{{-dist, 0}, {dist, 0}, {0, -dist}, {0, dist}}
}

func getOddSeries(min int, max int) []int {
	var oddIndices []int
	for i := min; i <= max; i++ {
		if i%2 != 0 {
			oddIndices = append(oddIndices, i)
		}
	}
	return oddIndices
}

// Get a random value from a slice of int
func getRandom(s []int) int {
	return s[rand.Intn(len(s))]
}

// func getTwoRandom(s []int) (int, int) {
// 	var coords []int
// 	for i := 0; i < 2; i++ {
// 		index := s[rand.Intn(len(s))]
// 		coords = append(coords, index)
// 	}
// 	return coords[0], coords[1]
// }

func (m *Maze) addEmptyNextToPlayer() {
	directions := generateDirections(1)
	for _, d := range directions {
		nextRow, nextCol := m.row+d.row, m.col+d.col
		if nextRow > 0 && nextRow < (m.maxRows-1) && nextCol > 0 && nextCol < (m.maxCols-1) && m.grid[nextRow][nextCol] == tiles.wall {
			m.grid[nextRow][nextCol] = tiles.empty
			break
		}
	}
}

// Adds a random cell
func (m *Maze) addRandomEmpty() {
	row := getRandom(getOddSeries(1, m.maxRows-2))
	col := getRandom(getOddSeries(1, m.maxCols-2))
	if m.grid[row][col] == tiles.empty {
		directions := generateDirections(1)
		for _, d := range directions { // this doesn't guarantee a cell is removed.
			nextRow, nextCol := row+d.row, col+d.col
			if nextRow > 0 && nextRow < (m.maxRows-1) && nextCol > 0 && nextCol < (m.maxCols-1) && m.grid[nextRow][nextCol] == tiles.wall {
				m.grid[nextRow][nextCol] = tiles.empty
				break
			}
		}
	}
}

func (m *Maze) MakePath() {
	m.addEmptyNextToPlayer()
}

// MovePlayer moves the player in the specified direction
func (m *Maze) MovePlayer(direction rune) {
	switch direction {
	case 'w':
		if m.row > 0 && m.grid[m.row-1][m.col] != tiles.wall {
			m.row--
		}
	case 'a':
		if m.col > 0 && m.grid[m.row][m.col-1] != tiles.wall {
			m.col--
		}
	case 's':
		if m.row < m.maxRows-1 && m.grid[m.row+1][m.col] != tiles.wall {
			m.row++
		}
	case 'd':
		if m.col < m.maxCols-1 && m.grid[m.row][m.col+1] != tiles.wall {
			m.col++
		}
	}
	m.trail = append(m.trail, struct{ row, col int }{m.row, m.col})
}

// Print displays the maze
func (m *Maze) Print() {
	ClearScreen()
	for i := 0; i < m.maxRows; i++ {
		rowStr := ""
		for j := 0; j < m.maxCols; j++ {
			cell := m.grid[i][j]
			if i == m.row && j == m.col {
				rowStr += tiles.player
			} else if containsTrail(i, j, m.trail) {
				rowStr += tiles.trail
			} else {
				rowStr += cell
			}
		}
		fmt.Println(rowStr)
	}
}

func containsTrail(row, col int, trail []struct{ row, col int }) bool {
	for _, t := range trail {
		if t.row == row && t.col == col {
			return true
		}
	}
	return false
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")

	// <?> overkill?
	// cmd := exec.Command("clear") // Use "cls" on Windows
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}
