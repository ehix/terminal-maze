package amaze

import (
	"fmt"
	"math/rand"
)

// <!> Need to decide on how these should be handled.
var tiles TileSet
var Name string

type Maze struct {
	maxRows int
	maxCols int
	grid    [][]string               // The entire maze grid
	row     int                      // Current row
	col     int                      // Current column
	start   []int                    // The starting index
	end     []int                    // The ending
	trail   []struct{ row, col int } // Store the player's trail
}

// Creates a new maze.
func NewMaze(maxRows int, maxCols int) *Maze {
	Name, tiles = SetRandomTiles()
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

// Sets the start and exit positions.
func (m *Maze) SetStartExit() {
	var randomCol []int
	// Get two random values for the start and end column indices
	for i := 0; i < 2; i++ {
		index := getRandom(getOddSeries(1, m.maxCols-2))
		randomCol = append(randomCol, index)
	}
	startCol, exitCol := randomCol[0], randomCol[1]
	// Add start and exit emoji
	m.grid[0][startCol] = tiles.start
	m.grid[m.maxRows-1][exitCol] = tiles.exit

	// Add to Maze for Generate function
	m.row = 0
	m.col = startCol
	m.start = []int{0, startCol}
	m.end = []int{m.maxRows - 1, exitCol}
}

// Asserts whether the game is over.
func (m *Maze) IsGameOver() bool {
	return m.grid[m.row][m.col] == tiles.exit
}

// Generates the maze using a Recursive Backtracking Algorithm.
func (m *Maze) Generate() {
	// Changes m.grid in place, by removing walls and replacing with spaces.
	stack := []struct{ row, col int }{{1, m.col}}
	// Initally, row: 1 (row: 0 will be walls), col: current/randomly decided.
	for len(stack) > 0 {
		current := stack[len(stack)-1] // Pop last item off the stack
		row, col := current.row, current.col
		m.grid[row][col] = tiles.empty
		var neighbors []struct{ row, col int }
		// Get all directions two spaces N, S, E, W
		directions := generateDirections(2)
		for _, d := range directions {
			nextRow, nextCol := row+d.row, col+d.col
			// If it's in boundaries and a wall...
			if nextRow > 0 && nextRow < (m.maxRows-1) && nextCol > 0 && nextCol < (m.maxCols-1) && m.grid[nextRow][nextCol] == tiles.wall {
				neighbors = append(neighbors, struct{ row, col int }{nextRow, nextCol})
			}
		}

		// If there's valid neighbors...
		if len(neighbors) > 0 {
			// Choose one at random and get it's indices
			randomIndex := rand.Intn(len(neighbors))
			randomNeighbor := neighbors[randomIndex]
			nRow, nCol := randomNeighbor.row, randomNeighbor.col
			// Make the tile empty, and also fill any adjecent cell on the path taken
			m.grid[nRow][nCol] = tiles.empty
			m.grid[row+(nRow-row)/2][col+(nCol-col)/2] = tiles.empty
			// Append the route taken as the start of the next pass
			stack = append(stack, struct{ row, col int }{nRow, nCol})
		} else {
			stack = stack[:len(stack)-1] // Drop from the stack
		}
	}
	// Add more empty cells with a certain probability
	m.addRandomEmpty(0.1)
}

// Sets the cell in the direction given to an empty tile.
func (m *Maze) addEmptyInDirection(d struct{ row, col int }) bool {
	complete := false
	nextRow, nextCol := m.row+d.row, m.col+d.col
	if nextRow > 0 && nextRow < (m.maxRows-1) && nextCol > 0 && nextCol < (m.maxCols-1) && m.grid[nextRow][nextCol] == tiles.wall {
		m.grid[nextRow][nextCol] = tiles.empty
		complete = true
	}
	return complete
}

// Removes a random number of wall tiles from the grid given some probability.
func (m *Maze) addRandomEmpty(rate float64) int {
	var removed int
	for i, row := range m.grid {
		if i != 0 && i != m.maxRows-1 {
			for j, tile := range row {
				if j != 0 && j != m.maxCols-1 {
					if tile == tiles.wall && rand.Float64() < rate {
						m.grid[i][j] = tiles.empty
						removed++
					}
				}
			}
		}
	}
	return removed
}

// Wrapper to addEmptyInDirection to capture key input.
func (m *Maze) MakePath(key rune) int {
	var removed int
	var d struct{ row, col int }
	switch key {
	case 'w':
		d.row = -1
	case 'a':
		d.col = -1
	case 's':
		d.row = 1
	case 'd':
		d.col = 1
	}
	if m.addEmptyInDirection(d) {
		removed++
	}
	return removed
}

// Wrapper for addRandomEmpty.
func (m *Maze) MakeEasy() int {
	return m.addRandomEmpty(0.5)
}

// Moves the player in the specified direction.
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

// Print displays the maze.
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

// Asserts if a given cell is in the trail passed.
func containsTrail(row, col int, trail []struct{ row, col int }) bool {
	for _, t := range trail {
		if t.row == row && t.col == col {
			return true
		}
	}
	return false
}
