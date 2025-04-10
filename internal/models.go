package internal

var (
	ants int // number of ants

	expectingStartRoom bool
	startRoomFound     = false

	expectingEndRoom bool
	endRoomFound     = false

	startRoom string // saves start room name
	endRoom   string // saves quick room name
)

// tunnels created after reading the file
var tunnels = make(map[string][]string)

// used for quick room lookup
var rooms = make(map[string]Room)

// Rooms created after reading the file
type Room struct {
	Name string
	X    int
	Y    int
}

// contains all the paths from DFS
var allPaths [][]string

var (
	bestStepPath          []string   // best path from step calculator
	bestStepDisjointPaths [][]string // includes bestStepPath and others
)
