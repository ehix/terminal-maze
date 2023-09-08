package amaze

import (
	"fmt"
	"time"
)

// Node represents a node in the A* algorithm
type Node struct {
	row, col int
	parent   *Node
}

// AutoSolve finds a path to the exit using the A* algorithm
func (m *Maze) AutoSolve() {
	var startNode *Node
	var exitNode *Node
	startNode = &Node{row: m.start[0], col: m.start[1]}
	exitNode = &Node{row: m.end[0], col: m.end[1]}
	// // Find the start and exit positions
	// for i := 0; i < rows; i++ {
	// 	for j := 0; j < cols; j++ {
	// 		if m.grid[i][j] == tiles.start {
	// 			startNode = &Node{row: i, col: j}
	// 		} else if m.grid[i][j] == tiles.exit {
	// 			exitNode = &Node{row: i, col: j}
	// 		}
	// 	}
	// }

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
			if row >= 0 && row < rows && col >= 0 && col < cols && m.grid[row][col] != tiles.wall && closedSet[row][col] != true {
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
