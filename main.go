package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
)

// https://pkg.go.dev/github.com/eiannone/keyboard@v0.0.0-20220611211555-0d226195f203
// https://pkg.go.dev/golang.org/x/crypto/ssh/terminal

// Maze dimensions
const rows, cols = 31, 41

// Cell dimensions (including spaces)
const cellWidth = 2

// Emoji maze cell types
const wall = "🧱"
const empty = "  "
const start = "🏁"
const exit = "🚪"
const player = "🧑"

// Maze represents the game's maze
type Maze struct {
	grid [][]string
	row  int
	col  int
}

// Node represents a node in the A* algorithm
type Node struct {
	row, col int
	parent   *Node
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")

	// <?> overkill?
	// cmd := exec.Command("clear") // Use "cls" on Windows
	// cmd.Stdout = os.Stdout
	// cmd.Run()
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
	}
}

// SetStartExit sets the start and exit positions
func (m *Maze) SetStartExit() {
	m.grid[0][1] = start
	m.grid[rows-1][cols-2] = exit
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

// IsGameOver checks if the game is over
func (m *Maze) IsGameOver() bool {
	return m.grid[m.row][m.col] == exit
}

// Print displays the maze
func (m *Maze) Print() {
	clearScreen()

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

// AutoSolve finds a path to the exit using the A* algorithm
func (m *Maze) AutoSolve() {
	var startNode *Node
	var exitNode *Node

	// Find the start and exit positions
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if m.grid[i][j] == start {
				startNode = &Node{row: i, col: j}
			} else if m.grid[i][j] == exit {
				exitNode = &Node{row: i, col: j}
			}
		}
	}

	if startNode == nil || exitNode == nil {
		fmt.Println("Error: Start or exit not found.")
		return
	}

	// Initialize the open and closed sets
	openSet := []*Node{startNode}
	closedSet := make(map[int]map[int]bool)

	for len(openSet) > 0 {
		// Find the node with the lowest total cost in the open set
		current := openSet[0]
		currentIndex := 0
		for i, node := range openSet {
			if f := m.g(current) + m.h(node, exitNode); f < m.g(current)+m.h(openSet[currentIndex], exitNode) {
				current = node
				currentIndex = i
			}
		}

		// Remove the current node from the open set and add it to the closed set
		openSet = append(openSet[:currentIndex], openSet[currentIndex+1:]...)
		if closedSet[current.row] == nil {
			closedSet[current.row] = make(map[int]bool)
		}
		closedSet[current.row][current.col] = true

		// Check if we have reached the exit
		if current.row == exitNode.row && current.col == exitNode.col {
			// Reconstruct the path
			path := []*Node{current}
			for current.parent != nil {
				current = current.parent
				path = append([]*Node{current}, path...)
			}

			// Move the player along the path
			for i := 1; i < len(path); i++ {
				next := path[i]
				if next.row > current.row {
					m.MovePlayer('s')
				} else if next.row < current.row {
					m.MovePlayer('w')
				} else if next.col > current.col {
					m.MovePlayer('d')
				} else if next.col < current.col {
					m.MovePlayer('a')
				}
				current = next
				m.Print()
				time.Sleep(80 * time.Millisecond) // Adjust the speed of auto-solve
			}
			break
		}

		// Explore the neighbors of the current node
		directions := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, dir := range directions {
			row, col := current.row+dir[0], current.col+dir[1]
			// Check if the neighbor is a valid cell and not in the closed set
			if row >= 0 && row < rows && col >= 0 && col < cols && m.grid[row][col] != wall && closedSet[row][col] != true {
				// Create a neighbor node
				neighbor := &Node{row: row, col: col, parent: current}

				// Check if the neighbor is not in the open set or has a lower cost
				neighborInOpenSet := false
				for _, openNode := range openSet {
					if openNode.row == neighbor.row && openNode.col == neighbor.col {
						neighborInOpenSet = true
						break
					}
				}

				if !neighborInOpenSet || m.g(neighbor) < m.g(current) {
					// Add the neighbor to the open set
					if !neighborInOpenSet {
						openSet = append(openSet, neighbor)
					}
				}
			}
		}
	}
}

// g calculates the cost from the start node to a node
func (m *Maze) g(n *Node) int {
	// In this simple version, we assume that each step has a cost of 1
	// You can customize this function to consider different costs for each step
	return 1
}

// h estimates the cost from a node to the goal (heuristic)
func (m *Maze) h(n *Node, goal *Node) int {
	// Use Manhattan distance as a heuristic (distance between two points)
	return abs(n.row-goal.row) + abs(n.col-goal.col)
}

// abs returns the absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// HandleKeyPress handles keyboard input for the game
func HandleKeyPress(m *Maze, playAgain *bool) {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
			*playAgain = false
			// done <- false
			break
		}

		if char == 'w' || char == 'a' || char == 's' || char == 'd' {
			m.MovePlayer(char)
		}
	}
}

// Function to display a simple text animation
func animatedtext() {
	screenWidth, screenHeight, err := term.GetSize(0)
	if err != nil {
		fmt.Println("Error here.")
	}
	text := `
	Solve
	the
	maze`

	textLines := strings.Split(text, "\n")
	numLines := len(textLines)
	originX := screenWidth/2 - len(textLines[0])/2
	originY := screenHeight/2 - numLines/2

	for i := 0; i < numLines; i++ {
		clearScreen()
		for j := 0; j <= i; j++ {
			fmt.Print("\033[H") // Move cursor to the top-left corner
			for k := 0; k < originY; k++ {
				fmt.Println() // Add padding lines
			}
			for k, line := range textLines {
				if k <= i {
					fmt.Printf("%s\n", strings.Repeat(" ", originX)+line)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	autoSolve := false
	if len(os.Args) > 1 && os.Args[1] == "auto" {
		autoSolve = true
	}

	playAgain := true

	// Handle Ctrl+C to exit gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		playAgain = false
	}()

	for playAgain {
		animatedtext()

		maze := NewMaze()
		maze.Generate()
		maze.SetStartExit()

		// Initialize the keyboard
		err := keyboard.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer keyboard.Close()

		if autoSolve {
			maze.AutoSolve()
		} else {
			for !maze.IsGameOver() {
				maze.Print()

				// Read a single key press
				char, key, err := keyboard.GetKey()
				if err != nil {
					log.Fatal(err)
				}

				if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
					playAgain = false
					break
				}

				if char == 'w' || char == 'a' || char == 's' || char == 'd' {
					maze.MovePlayer(char)
				}
			}

			if playAgain {
				fmt.Println("Congratulations! You've reached the exit (🚪).")
				// Ask the player if they want to play another maze
				fmt.Print("Do you want to play another maze? (y/n): ")
				var response string
				fmt.Scan(&response)
				playAgain = strings.ToLower(response) == "y"
			}
		}
	}
}
