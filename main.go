package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ehix/terminalmaze/amaze"
	"github.com/eiannone/keyboard"
)

// https://pkg.go.dev/github.com/eiannone/keyboard@v0.0.0-20220611211555-0d226195f203
// https://pkg.go.dev/golang.org/x/crypto/ssh/terminal

type session struct {
	// Number of consecutive games played
	played int
	// Keep playing if true
	cont bool
}

// Function to display a simple text animation
func animateText() {
	// screenWidth, screenHeight, err := term.GetSize(0)
	// if err != nil {
	// 	fmt.Println("Error here.")
	// }
	// text := `
	// Solve
	// the
	// maze`

	textLines := strings.Split(banner, "\n")
	numLines := len(textLines)
	// originX := screenWidth/2 - len(textLines[0])/2
	// originY := screenHeight/2 - numLines/2
	originX := 0
	originY := 0

	for i := 0; i < numLines; i++ {
		amaze.ClearScreen()
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
			time.Sleep(10 * time.Millisecond)
		}
	}
	time.Sleep(10 * time.Millisecond)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// animateText()
	autoSolve := false
	if len(os.Args) > 1 && os.Args[1] == "auto" {
		autoSolve = true
	}

	s := session{played: 0, cont: true}
	d := amaze.NewDifficulty()
	for s.cont {
		d.SetCurrentStage(s.played)
		row, col := d.GetDimensions()
		maze := amaze.NewMaze(row, col)
		maze.SetStartExit()
		maze.Generate()

		// Initialize the keyboard
		err := keyboard.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer keyboard.Close()

		if autoSolve {
			// Wait for a signal from a channel, otherwise run autosolve.
			// Running this in a go routine breaks it.
			maze.AutoSolve()
			break
		} else {
			for !maze.IsGameOver() {
				maze.Print()
				fmt.Printf("solved: %d\ndimensions: %dx%d", s.played, col, row)
				// Read a single key press
				char, key, err := keyboard.GetKey()
				if err != nil {
					log.Fatal(err)
				}

				if key == keyboard.KeyCtrlA {
					// Auto solve the current maze
					maze.AutoSolve()
				}

				if key == keyboard.KeyCtrlE {
					// Remove random wall tiles from maze
					maze.MakePath()
				}

				if key == keyboard.KeyCtrlR {
					// Restart a new maze
					break
				}

				if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
					// Terminate and close keyboard
					s.cont = false
					break
				}

				if char == 'w' || char == 'a' || char == 's' || char == 'd' {
					maze.MovePlayer(char)
				}
			}
		}
		s.played++
	}
	amaze.ClearScreen()
}

//https://patorjk.com/software/taag/#p=display&f=Graffiti&t=Type%20Something%20
// var banner = `
// ▄▄▄█████▓▓█████  ██▀███   ███▄ ▄███▓ ██▓ ███▄    █  ▄▄▄       ██▓        ███▄ ▄███▓ ▄▄▄      ▒███████▒▓█████     ██▓
// ▓  ██▒ ▓▒▓█   ▀ ▓██ ▒ ██▒▓██▒▀█▀ ██▒▓██▒ ██ ▀█   █ ▒████▄    ▓██▒       ▓██▒▀█▀ ██▒▒████▄    ▒ ▒ ▒ ▄▀░▓█   ▀    ▓██▒
// ▒ ▓██░ ▒░▒███   ▓██ ░▄█ ▒▓██    ▓██░▒██▒▓██  ▀█ ██▒▒██  ▀█▄  ▒██░       ▓██    ▓██░▒██  ▀█▄  ░ ▒ ▄▀▒░ ▒███      ▒██▒
// ░ ▓██▓ ░ ▒▓█  ▄ ▒██▀▀█▄  ▒██    ▒██ ░██░▓██▒  ▐▌██▒░██▄▄▄▄██ ▒██░       ▒██    ▒██ ░██▄▄▄▄██   ▄▀▒   ░▒▓█  ▄    ░██░
//   ▒██▒ ░ ░▒████▒░██▓ ▒██▒▒██▒   ░██▒░██░▒██░   ▓██░ ▓█   ▓██▒░██████▒   ▒██▒   ░██▒ ▓█   ▓██▒▒███████▒░▒████▒   ░██░
//   ▒ ░░   ░░ ▒░ ░░ ▒▓ ░▒▓░░ ▒░   ░  ░░▓  ░ ▒░   ▒ ▒  ▒▒   ▓▒█░░ ▒░▓  ░   ░ ▒░   ░  ░ ▒▒   ▓▒█░░▒▒ ▓░▒░▒░░ ▒░ ░   ░▓
//     ░     ░ ░  ░  ░▒ ░ ▒░░  ░      ░ ▒ ░░ ░░   ░ ▒░  ▒   ▒▒ ░░ ░ ▒  ░   ░  ░      ░  ▒   ▒▒ ░░░▒ ▒ ░ ▒ ░ ░  ░    ▒ ░
//   ░         ░     ░░   ░ ░      ░    ▒ ░   ░   ░ ░   ░   ▒     ░ ░      ░      ░     ░   ▒   ░ ░ ░ ░ ░   ░       ▒ ░
//             ░  ░   ░            ░    ░           ░       ░  ░    ░  ░          ░         ░  ░  ░ ░       ░  ░    ░
//                                                                                              ░
// `

var banner = `
▄▀▀▀█▀▀▄  ▄▀▀█▄▄▄▄  ▄▀▀▄▀▀▀▄  ▄▀▀▄ ▄▀▄  ▄▀▀█▀▄    ▄▀▀▄ ▀▄  ▄▀▀█▄   ▄▀▀▀▀▄     
█    █  ▐ ▐  ▄▀   ▐ █   █   █ █  █ ▀  █ █   █  █  █  █ █ █ ▐ ▄▀ ▀▄ █    █     
▐   █       █▄▄▄▄▄  ▐  █▀▀█▀  ▐  █    █ ▐   █  ▐  ▐  █  ▀█   █▄▄▄█ ▐    █     
   █        █    ▌   ▄▀    █    █    █      █       █   █   ▄▀   █     █      
 ▄▀        ▄▀▄▄▄▄   █     █   ▄▀   ▄▀    ▄▀▀▀▀▀▄  ▄▀   █   █   ▄▀    ▄▀▄▄▄▄▄▄▀
█          █    ▐   ▐     ▐   █    █    █       █ █    ▐   ▐   ▐     █        
▐          ▐                  ▐    ▐    ▐       ▐ ▐                  ▐        
▄▀▀▄ ▄▀▄  ▄▀▀█▄   ▄▀▀▀▀▄   ▄▀▀█▄▄▄▄                                           
█  █ ▀  █ ▐ ▄▀ ▀▄ █     ▄▀ ▐  ▄▀   ▐                                          
▐  █    █   █▄▄▄█ ▐ ▄▄▀▀     █▄▄▄▄▄                                           
  █    █   ▄▀   █   █        █    ▌                                           
▄▀   ▄▀   █   ▄▀     ▀▄▄▄▄▀ ▄▀▄▄▄▄                                            
█    █    ▐   ▐          ▐  █    ▐                                            
▐    ▐                      ▐                                                 
`
