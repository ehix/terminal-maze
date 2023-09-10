package amaze

type Difficulty struct {
	min          int
	max          int
	stages       int
	perStage     int
	currentStage []int
	sections     [][]int
}

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

func (d *Difficulty) SetCurrentStage(played int) {
	if played < (d.stages * d.perStage) {
		if played%d.perStage == 0 {
			d.currentStage = append(d.currentStage, d.sections[played/d.perStage]...)
		}
	}
}

func (d *Difficulty) GetDimensions() (int, int) {
	var dims []int
	for i := 0; i < 2; i++ {
		dims = append(dims, getRandom(d.currentStage))
	}
	return dims[0], dims[1]
}

func NewDifficulty() Difficulty {
	// n = 34 = (17 * 2) = len(0-17) i.e. (2*((3*6)-1) - (41-7)) = 0
	// or (2*((2*9)-1), etc..
	d := Difficulty{
		min:      7,
		max:      41,
		stages:   9, //6
		perStage: 2, //3
	}
	d.setSections()
	return d
}
