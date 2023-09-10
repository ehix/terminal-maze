package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/ehix/terminalmaze/pkg/amaze"
	"github.com/eiannone/keyboard"
)

type session struct {
	// Number of consecutive games played
	played int
	// Number of cheats
	cheated int
}

// Display a simple text animation
func animateText(text string, interval time.Duration, pause time.Duration) {
	textLines := strings.Split(text, "\n")
	numLines := len(textLines)
	originX := 0
	originY := 0
	for i := 0; i < numLines; i++ {
		amaze.ClearScreen()
		for j := 0; j <= i; j++ {
			cursorTopLeft()
			for k := 0; k < originY; k++ {
				fmt.Println() // Add padding lines
			}
			for k, line := range textLines {
				if k <= i {
					fmt.Printf("%s\n", strings.Repeat(" ", originX)+line)
				}
			}
			time.Sleep(interval)
		}
	}
	time.Sleep(pause)
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	printBanner()
	printControls()
	// Session implements some minor controls and metrics
	s := session{played: 0, cheated: 0}
	// Difficulty managages the maze size and the play count increases
	d := amaze.NewDifficulty()

gameplay:
	for {
		// Dictate the stage size
		d.SetCurrentStage(s.played)
		row, col := d.GetDimensions()
		// Create the maze
		maze := amaze.NewMaze(row, col)
		maze.SetStartExit()
		maze.Generate()

		// Increment level if complete without auto
		var auto bool
		// Print devil emoji in metrics
		var cheater bool

		// Initialize the keyboard
		err := keyboard.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer keyboard.Close()

		for !maze.IsGameOver() {
			// Display the maze and metrics
			maze.Print()
			printMetrics(amaze.Name, s.played, s.cheated, cheater, col, row)

			// Read a single key press
			char, key, err := keyboard.GetKey()
			if err != nil {
				log.Fatal(err)
			}

			switch key {
			// Auto solve the current maze
			case keyboard.KeyCtrlE:
				auto = true
				maze.AutoSolve()
			// Remove random wall tiles with some probability
			case keyboard.KeyCtrlQ:
				s.cheated += maze.MakeEasy()
				cheater = true
			// Remove a wall in a given direction
			case keyboard.KeyCtrlW:
				s.cheated += maze.MakePath('w')
				cheater = true
			case keyboard.KeyCtrlA:
				s.cheated += maze.MakePath('a')
				cheater = true
			case keyboard.KeyCtrlS:
				s.cheated += maze.MakePath('s')
				cheater = true
			case keyboard.KeyCtrlD:
				s.cheated += maze.MakePath('d')
				cheater = true
			// Regenerate a new maze
			case keyboard.KeyCtrlR:
				continue gameplay
			// Terminate and close keyboard
			case keyboard.KeyEsc, keyboard.KeyCtrlC:
				break gameplay
			}

			switch char {
			case 'w', 'a', 's', 'd':
				maze.MovePlayer(char)
			case 'c':
				printControls()
			}
		}

		if !auto {
			s.played++
		}
	}
	amaze.ClearScreen()
}

func printBanner() {
	interval := time.Duration(25) * time.Millisecond
	pause := time.Duration(50) * time.Millisecond
	animateText(bannerHollow, interval, pause)

	amaze.ClearScreen()
	fmt.Println(bannerFill)
	time.Sleep(time.Duration(2) * time.Second)
}

func printControls() {
	interval := time.Duration(20) * time.Millisecond
	pause := time.Duration(3) * time.Second
	animateText(controls, interval, pause)
}

func printMetrics(name string, played int, cheated int, cheater bool, width int, height int) {
	sformat := 11
	sname := "%-*s\n"
	if cheater {
		sname = strings.Replace(sname, "%-*s\n", "%-*s👹\n", 1)
	}
	fmt.Printf(sname, sformat, name)
	fmt.Printf("%-*s%-dx%d\n", sformat, "├─size:", width, height)
	fmt.Printf("%-*s%-2d\n", sformat, "├─solved:", played)
	fmt.Printf("%-*s%-2d", sformat, "├─cheated:", cheated)
	cursorTopLeft()
}

func cursorTopLeft() {
	fmt.Print("\033[H")
}

var bannerFill = `
                                  _                      
 _/_               ♥             //                      
 /  _  _   _ _ _  ,  _ _   __,  //    _ _ _   __,  __, _ 
(__(/_/ (_/ / / /_(_/ / /_(_/(_(/_   / / / /_(_/(_/_/_(/_
                                                  (/     
`
var bannerHollow = `
                                  _                      
 _/_               ♡             //                      
 /  _  _   _ _ _  ,  _ _   __,  //    _ _ _   __,  __, _ 
(__(/_/ (_/ / / /_(_/ / /_(_/(_(/_   / / / /_(_/(_/_/_(/_
                                                  (/     
`

var controls = `
controls/
├─ move/
│  ├─ wasd
├─ skip/
│  ├─ ctrl + r
├─ autosolve/
│  ├─ ctrl + e
├─ cheats/
│  ├─ make-easy/ 
│  │  ├─ ctrl + q
│  ├─ remove-block/ 
│  │  ├─ ctrl + wasd
├─ exit/
│  ├─ ctrl + c
│  ├─ esc
├─ credits/
│  ├─ ehix@1694300400
├─ ENJOY.md
`
