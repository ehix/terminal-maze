package amaze

// Manages the maze dimensions as the play count increases.
type Difficulty struct {
	min          int
	max          int
	stages       int
	perStage     int
	currentStage []int
	sections     [][]int
}

// Creates a new Difficulty.
func NewDifficulty() Difficulty {
	// The general premise for these constants.
	// Given, the mazes are always made with dimensions that are odd:
	// If min-max = n, n/2 = odd integers. Therefore, len([n/2]) = 18 games.
	// Factors of 18 include, 1, 2, 3, 6, 9 and 18.
	// i.e. below we have, (2*((3*6)-1) - (41-7)) = 0.
	// So, 6 stages with 3 dimensions in each is permissable.
	d := Difficulty{
		min:      7,
		max:      41,
		stages:   9, // could be 6
		perStage: 2, // could be 3, etc...
	}
	d.setSections()
	return d
}

// Turns a odd series of integers into equal sections.
func (d *Difficulty) setSections() {
	s := getOddSeries(d.min, d.max)
	sectionSize := len(s) / d.stages
	// Divide the slice into sections
	for i := 0; i < d.stages; i++ {
		start := i * sectionSize
		end := (i + 1) * sectionSize
		// Check if it's the last section and adjust the end index if needed
		if i == d.stages-1 {
			end = len(s)
		}
		d.sections = append(d.sections, s[start:end])
	}
}

// Add a number of potential maze dimensions to the pool given the play count.
func (d *Difficulty) SetCurrentStage(played int) {
	if played < (d.stages * d.perStage) {
		if played%d.perStage == 0 {
			d.currentStage = append(d.currentStage, d.sections[played/d.perStage]...)
		}
	}
}

// Returns the calculated dimensions the maze should be.
func (d *Difficulty) GetDimensions() (int, int) {
	var dims []int
	for i := 0; i < 2; i++ {
		dims = append(dims, getRandom(d.currentStage))
	}
	return dims[0], dims[1]
}
