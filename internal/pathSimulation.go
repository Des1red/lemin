package internal

import (
	"fmt"
	"sort"
	"strings"
)

func Simulate() {
	// Determine the number of paths to use based on the number of ants
	numPaths := 1
	if ants > 3 && ants <= 6 {
		numPaths = 3
	} else {
		numPaths = len(bestStepDisjointPaths)
	}

	// Run the simulation using the selected paths
	simulateAnts(bestStepDisjointPaths[:numPaths])
}

// Ant represents the state of an ant in the simulation.
type Ant struct {
	ID    int // Ant identifier
	Path  int // Index into the paths slice (which path the ant is following)
	Index int // Current index (position) on that path
}

// moveAntsInTransit processes ants in transit so that their current room is freed
// immediately as they start to move. It returns the updated list of ants still in transit
// along with the movement commands for this turn.
func moveAntsInTransit(antsInTransit []Ant, paths [][]string, occupied map[string]bool) ([]Ant, []string) {
	output := []string{}
	newTransit := []Ant{}

	// Process ants in reverse order so that ants closer to the end are processed first.
	for i := len(antsInTransit) - 1; i >= 0; i-- {
		ant := antsInTransit[i]
		path := paths[ant.Path]
		currentRoom := path[ant.Index]
		nextIndex := ant.Index + 1

		// Free the current room immediately, since the ant is going to try to leave.
		delete(occupied, currentRoom)

		// Check if there is a valid next room and if it is not occupied.
		// Check if there is a next room and if it is not occupied.
		if nextIndex < len(path) && !occupied[path[nextIndex]] {
			move := fmt.Sprintf("L%d-%s", ant.ID, path[nextIndex])
			output = append(output, move)
			ant.Index = nextIndex

			// If this ant has not yet reached the final room, add it back to transit.
			if nextIndex < len(path)-1 {
				newTransit = append(newTransit, ant)
				occupied[path[nextIndex]] = true // Reserve the room.
			}
			// If the ant reaches the final room, do not add it back.
		} else {
			// If the ant couldn’t move, re-reserve its current room and keep it in transit.
			occupied[currentRoom] = true
			newTransit = append(newTransit, ant)
		}
	}

	return newTransit, output
}

// spawnNewAnts checks each available path and spawns a new ant if its first room is free.
// It returns the updated antsInTransit slice, the move commands (as strings) produced by spawning,
// and the updated nextAnt counter.
func spawnNewAnts(antsInTransit []Ant, paths [][]string, nextAnt int, totalAnts int, occupied map[string]bool) ([]Ant, []string, int) {
	output := []string{}

	// Try to spawn a new ant on each path (if available) for as long as there are ants to spawn.
	for pathIndex := 0; pathIndex < len(paths) && nextAnt <= totalAnts; pathIndex++ {
		// The spawn room is the first room after "start" (index 1).
		if !occupied[paths[pathIndex][1]] {
			antsInTransit = append(antsInTransit, Ant{ID: nextAnt, Path: pathIndex, Index: 1})
			move := fmt.Sprintf("L%d-%s", nextAnt, paths[pathIndex][1])
			output = append(output, move)
			occupied[paths[pathIndex][1]] = true // Reserve the spawn room.
			nextAnt++
		}
	}

	return antsInTransit, output, nextAnt
}

// simulateAnts is the main simulation function. It repeatedly:
//  1. Moves ants already in transit,
//  2. Spawns new ants,
//  3. Prints all moves for that turn,
//  4. And finally updates the turn count.
//
// The simulation stops when no moves occur on a turn.
func simulateAnts(paths [][]string) {
	if len(paths) == 0 {
		Log("no valid paths to simulate.", "error")
		return
	}

	// These values manage the ant simulation state.
	antsInTransit := []Ant{}
	nextAnt := 1
	turn := 1

	// occupied tracks the rooms in use during the current turn.
	occupied := make(map[string]bool)

	// The simulation loop runs until no moves are produced.
	for {
		turnOutput := []string{}

		// Phase 1: Move ants already in transit.
		var moves []string
		antsInTransit, moves = moveAntsInTransit(antsInTransit, paths, occupied)
		turnOutput = append(turnOutput, moves...)

		// Phase 2: Spawn new ants on the paths.
		antsInTransit, moves, nextAnt = spawnNewAnts(antsInTransit, paths, nextAnt, ants, occupied)
		turnOutput = append(turnOutput, moves...)

		// If no moves were made this turn, the simulation is complete.
		if len(turnOutput) == 0 {
			break
		}

		// Sort moves by ant ID for consistent ordering in output.
		sort.Slice(turnOutput, func(i, j int) bool {
			var id1, id2 int
			fmt.Sscanf(turnOutput[i], "L%d-", &id1)
			fmt.Sscanf(turnOutput[j], "L%d-", &id2)
			return id1 < id2
		})
		fmt.Println(strings.Join(turnOutput, " "))

		// If all ants have been spawned and none remain in transit, break early.
		if nextAnt > ants && len(antsInTransit) == 0 {
			break
		}

		turn++
		// Reset the occupied map for the next turn.
		occupied = make(map[string]bool)
	}

	// Log the total number of turns (only count turns in which moves were executed).
	Log(fmt.Sprintf("Total number of turns: %d\n", turn), "debug")
}
