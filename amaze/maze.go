package amaze

import (
	"fmt"
	"math/rand"
)

// Maze dimensions
const rows, cols = 31, 31

// Cell dimensions (including spaces)
const cellWidth = 2

// Emoji maze cell types
const wall = "🧱"
const empty = "  "
const start = "🏁"
const exit = "🚪"
const player = "🧑"
const trail = "👣"

type Maze struct {
	grid [][]string
	row  int
	col  int
}

type Trail struct {
	trailGrid [][]string               // Store the trail grid
	trail     []struct{ row, col int } // Store the player's trail
}

// NewMaze creates a new maze
func NewMaze() *Maze {
	m := &Maze{
		grid: make([][]string, rows),
		row:  0,
		col:  1,
	}

	for i := range m.grid {
		m.grid[i] = make([]string, cols)
		for j := range m.grid[i] {
			m.grid[i][j] = wall
		}
	}

	return m
}

// SetStartExit sets the start and exit positions
func (m *Maze) SetStartExit() {
	m.grid[0][1] = start
	m.grid[rows-1][cols-2] = exit
}

// IsGameOver checks if the game is over
func (m *Maze) IsGameOver() bool {
	return m.grid[m.row][m.col] == exit
}

// Generate generates the maze using Recursive Backtracking Algorithm
func (m *Maze) Generate() {
	stack := []struct{ row, col int }{{1, 1}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		row, col := current.row, current.col
		m.grid[row][col] = empty

		var neighbors []struct{ row, col int }
		directions := []struct{ row, col int }{{-2, 0}, {2, 0}, {0, -2}, {0, 2}}

		for _, d := range directions {
			nextRow, nextCol := row+d.row, col+d.col
			if nextRow >= 0 && nextRow < rows && nextCol >= 0 && nextCol < cols && m.grid[nextRow][nextCol] == wall {
				neighbors = append(neighbors, struct{ row, col int }{nextRow, nextCol})
			}
		}

		if len(neighbors) > 0 {
			randomIndex := rand.Intn(len(neighbors))
			randomNeighbor := neighbors[randomIndex]
			nRow, nCol := randomNeighbor.row, randomNeighbor.col

			m.grid[nRow][nCol] = empty
			m.grid[row+(nRow-row)/2][col+(nCol-col)/2] = empty

			stack = append(stack, struct{ row, col int }{nRow, nCol})
		} else {
			stack = stack[:len(stack)-1]
		}

		// Add more cells with a certain probability
		if rand.Float64() < 0.1 { // make this scale with every game?
			m.addRandomCell()
		}
	}
}

// Adds a random cell
func (m *Maze) addRandomCell() {
	// rand.Intn(max - min + 1) + min
	row := rand.Intn((rows-3)-3+1) + 3
	col := rand.Intn((cols-2)-3+1) + 3

	if row%2 == 0 {
		row++
	}
	if col%2 == 0 {
		col++
	}
	m.grid[row][col] = empty
}

// MovePlayer moves the player in the specified direction
func (m *Maze) MovePlayer(direction rune) {
	switch direction {
	case 'w':
		if m.row > 0 && m.grid[m.row-1][m.col] != wall {
			m.row--
		}
	case 'a':
		if m.col > 0 && m.grid[m.row][m.col-1] != wall {
			m.col--
		}
	case 's':
		if m.row < rows-1 && m.grid[m.row+1][m.col] != wall {
			m.row++
		}
	case 'd':
		if m.col < cols-1 && m.grid[m.row][m.col+1] != wall {
			m.col++
		}
	}
}

// Print displays the maze
func (m *Maze) Print() {
	ClearScreen()

	for i := 0; i < rows; i++ {
		rowStr := ""
		for j := 0; j < cols; j++ {
			cell := m.grid[i][j]
			if i == m.row && j == m.col {
				rowStr += player
			} else {
				rowStr += cell
			}
		}
		fmt.Println(rowStr)
	}
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")

	// <?> overkill?
	// cmd := exec.Command("clear") // Use "cls" on Windows
	// cmd.Stdout = os.Stdout
	// cmd.Run()
}
