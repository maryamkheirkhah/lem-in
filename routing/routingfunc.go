package routing

import (
	"errors"
	"fmt"
	"lem-in/sys"
	"strconv"
	"strings"
)

var (
	CurrentTurnStr    string // For moving ants
	AntID             = 1    // For moving ants
	TotalAntsFinished = 0    // For moving ants
	Routes            [][]*sys.Room
	AntGrouping       []int
)

// PRINTING FUNCTIONS FOR DE-BUGGING
/*
func printRoute(route []*sys.Room) {
	fmt.Print("ROUTE len: " + strconv.Itoa(len((route))) + "  . Route: ")
	for _, room := range route {
		fmt.Print(room.Name + " ")
	}
	fmt.Println()
}

func printRouteSet(setRoutes [][]*sys.Room) {
	for _, route := range setRoutes {
		printRoute(route)
	}
}
*/

/*
maxInt is a function that takes two integers and returns the value of the largest positive integer, along
with an error value. If both integers are smaller than or equal to zero, a non-nil error is returned.
*/
func maxInt(int1, int2 int) (int, error) {
	if int1 <= 0 && int2 <= 0 {
		return 0, errors.New("\nERROR: internal malfunction, function \" maxInt \" " +
			"called with both inputs less than or equal to zero")
	}

	if int1 > 0 && int2 > 0 {
		if int1 > int2 {
			return int1, nil
		}
		return int2, nil
	}

	// Case if one value is less than or equal to zero
	if int1 > 0 {
		return int1, nil
	}
	return int2, nil
}

/*
findByName takes an input room name in the form of a string, and searches the global sys.Network
variable for the room with a matching name. It then returns a pointer to this room within the global
Network struct.
*/
func findByName(name string) *sys.Room {
	for _, room := range sys.Network {
		if room.Name == name {
			return &room
		}

	}
	return nil
}

/*
sortRoutes takes an input slice of routes ([][]*sys.Room) and applies a BUBBLE SORT algorithm to sort them
in ascending order of route length. The sorted slice is then returned, along with a non-nil error if the input
slice of routes has a length of zero.
*/
func sortRoutes(setRoutes [][]*sys.Room) ([][]*sys.Room, error) {
	if len(setRoutes) == 0 {
		return setRoutes, errors.New("\nERROR: internal malfunction, the function \" sortRoutes \" " +
			"given an input slice of routes with length zero")
	}
	changeCounter := 1
	for changeCounter > 0 {
		changeCounter = 0
		for i := 1; i < len(setRoutes); i++ {
			if len(setRoutes[i]) < len(setRoutes[i-1]) {
				setRoutes[i], setRoutes[i-1] = setRoutes[i-1], setRoutes[i]
				changeCounter++
			}
		}
	}
	return setRoutes, nil
}

/*
convertStringToSlice takes an input string of paths, with individual routes assumed to be separated
by a ",", and individual rooms within the paths separated by a " ". The function parses this input
string to produce a slice of strings, with each string representing an individual route. It then
further separates these strings according to room names, calling the local findByName function to
write to the eventual output variable, a slice of routes, with each route itself a slice of pointers
to rooms in the global Network variable (*sys.Room).
*/
func convertStringToSlice(strAllPaths string) ([][]*sys.Room, error) {
	strArrAllPath := strings.Split(strAllPaths, ",")
	arrAllPaths := make([][]*sys.Room, len(strArrAllPath))
	for j, path := range strArrAllPath {
		arrPath := strings.Split(path, " ")
		for _, name := range arrPath {
			arrAllPaths[j] = append(arrAllPaths[j], findByName(name))
		}
	}
	return arrAllPaths, nil
}

/*
dfsString is a function that performs a depth-first search on a given room in a system of interconnected rooms.
It returns a string with all paths from the starting room to the end room, separated by commas.
The function takes a pointer to the current room being searched (currentRoom *sys.Room), a string with the
current path taken (strPath string), and a string with all paths found so far (strAllPaths string).
If the current room has already been visited, the function returns strAllPaths. Otherwise, it marks the
current room as visited and adds the room's name to the current path. If the current room is the end room,
the function adds the current path to strAllPaths and returns the resulting string, with all paths separated
by commas. If the current room is not the end room, the function iterates over all linked rooms and calls itself
recursively on each one. Finally, the function marks the current room as unvisited and returns a string with
all the compiled paths.
*/
func dfsString(currentRoom *sys.Room, strPath, strAllPaths string) string {
	if currentRoom.Visited {
		return strAllPaths
	}
	currentRoom.Visited = true
	strPath += " " + currentRoom.Name
	if currentRoom.Name == sys.End.Name {
		currentRoom.Visited = false
		strAllPaths += "," + strPath[1:]
		return strAllPaths
	} else {
		for _, next := range currentRoom.Links {
			strAllPaths = dfsString(next, strPath, strAllPaths)
		}
	}
	currentRoom.Visited = false
	return strAllPaths
}

/*
runningDFS is a function that calls a depth-first search on a system of interconnected rooms to find all
routes from the start room to the end room. It returns a slice of slices of pointers to Room objects,
representing the routes ordered in terms of ascending length, and an error value. If an error is returned
at any point, or if no valid routes are found, the function returns an empty slice of slices and the non-nil
error value.
*/
func runningDFS() ([][]*sys.Room, error) {
	// Initialise working variables
	var strAllPaths string
	var err error
	var allRoutes [][]*sys.Room

	if sys.Start == nil || sys.End == nil {
		return allRoutes, errors.New("\nERROR: internal malfunction, the function \" runningDFS \" " +
			"called while the sys.Start and/or sys.End rooms are empty")
	}

	// Perform call depth-first-search algorithms, and convert to relevant outputs.
	strAllPaths = dfsString(sys.Start, "", strAllPaths)
	if err != nil {
		return allRoutes, err
	}

	// Return error if no valid routes found
	if len(strAllPaths) == 0 {
		return allRoutes, errors.New("\nERROR: invalid data format, no valid routes between " +
			"start and end rooms could be found")
	}

	allRoutes, err = convertStringToSlice(strAllPaths[1:])
	if err != nil {
		return allRoutes, err
	}

	// Return error if no valid routes found
	if len(allRoutes) == 0 {
		return allRoutes, errors.New("\nERROR: invalid data format, no valid routes between " +
			"start and end rooms could be found")
	}

	// Sort routes in ascending order of length before returning
	allRoutes, err = sortRoutes(allRoutes)
	if err != nil {
		return allRoutes, err
	}

	return allRoutes, nil
}

/*
checkRouteConflict takes two input routes ([]*Room) and scans their respective nodes for conflicts.
If any of the intermediate nodes are found in both routes, a "true" boolean is returned. If no
duplicate intermediate nodes are found, a "false" boolean is returned. An error value is also returned
in all cases, and is non-nil if input errors are detected.
*/
func checkRouteConflict(route1, route2 []*sys.Room) (bool, error) {
	if len(route1) == 0 || len(route2) == 0 {
		return false, errors.New("\nERROR: internal malfunction, route of zero length used as " +
			"input to \" checkRouteConflict \" function")
	}
	for i, room1 := range route1 {
		if i == 0 || i == len(route1)-1 {
			continue // Skip start / end room comparisons
		}
		for j, room2 := range route2 {
			if j == 0 || j == len(route2)-1 {
				continue // Skip start / end room comparisons
			}
			if room1.Name == room2.Name {
				return true, nil
			}
		}
	}
	return false, nil
}

/*
createConflictMap takes an input slice of routes, and parses each one to construct a map where each key
is the integer index of a particular route within the input slice of routes, and the corresponding values are
the integer indices of all routes within the input slice that conflict with the route specified in the key.
The map is then returned, along with an error value, which is non-nil in the event that any local function
calls return an error, or if the input slice of routes has a length of zero.
*/
func createConflictMap(allRoutes [][]*sys.Room) (map[int][]int, error) {
	// Initialize output variables
	conflictMap := make(map[int][]int, len(allRoutes))
	var hasConflict bool
	var err error

	if len(allRoutes) == 0 {
		return conflictMap, errors.New("\nERROR: internal malfunction, allRoutes variable with zero length " +
			"used as input to \" createConflictMap \" function")
	}

	// Iterate over all routes in the input slice
	for i, route1 := range allRoutes {
		// Initialize a slice (with max. possible capacity) to hold the conflicting routes for the current route
		conflicts := make([]int, 0, len(allRoutes))

		// Iterate over all routes again to check for conflicts
		for j, route2 := range allRoutes {
			// Skip the current route if it is the same as the route being checked
			if i == j {
				continue
			}

			// Check if the current route conflicts with the route being checked
			hasConflict, err = checkRouteConflict(route1, route2)
			if err != nil {
				return conflictMap, err
			} else if hasConflict {
				conflicts = append(conflicts, j)
			}
		}

		// Add the conflicting routes for the current route to the conflict map
		conflictMap[i] = conflicts
	}

	return conflictMap, nil
}

/*
findBestPath is a function that takes a slice of slices of pointers to Room objects representing routes, and
a slice of integers representing the number of ants assigned to each route. For each route, it calculates the
total length by adding the number of ants already assigned to it to the actual length of the route. If this
total length is less than the current minimum length, the function updates the minimum length and index
variables. It returns the index of the route with the minimum length, taking into account the number of ants
already assigned to it.
*/
func findBestPath(routeCombo [][]*sys.Room, antGrouping []int) int {
	minLength := antGrouping[0] + len(routeCombo[0])
	minIndex := 0
	for i, route := range routeCombo {
		if (antGrouping[i] + len(route)) < minLength {
			minLength = antGrouping[i] + len(route)
			minIndex = i
		}
	}
	return minIndex
}

/*
calcAntGrouping is a function that takes a slice of slices of pointers to Room objects representing routes,
and iterates over the total number of ants (referenced to by the global sys.TotalAntNbr variable) and assigns
each one to the best route using the findBestPath function. Finally, it returns the ant grouping slice and an
error value, which is non-nil if the input slice has a length of zero.
*/
func calcAntGrouping(routeCombo [][]*sys.Room) ([]int, error) {
	if len(routeCombo) == 0 {
		return []int{}, errors.New("\nERROR: internal malfunction, \" calcAntGrouping \" function called " +
			"with an input slice of routes with a length of zero")
	}

	antGrouping := make([]int, len(routeCombo))
	index := 0
	for i := 0; i < sys.TotalAntNbr; i++ {
		index = findBestPath(routeCombo, antGrouping)
		antGrouping[index]++
	}
	return antGrouping, nil
}

/*
intersection is a function that takes two input slice of integers and returns a slice of integers containing
the intersecting values of the two slices
*/
func intersection(a, b []int) []int {
	var result []int
	for _, aVal := range a {
		for _, bVal := range b {
			if aVal == bVal {
				result = append(result, aVal)
			}
		}
	}
	return result
}

/*
contains is a function that iterates over an input slice of integers and returns true if the integer element
e is found in the slice, and false otherwise. It is used in the getNonConflictingCombinations function to
check if the current route is already in the combination.
*/
func contains(intSlice []int, e int) bool {
	for _, nbr := range intSlice {
		if nbr == e {
			return true
		}
	}
	return false
}

/*
getNonConflictingCombinations is a recursive function that takes a map of with routeKey (int) keys and a slice
of the corresponding conflicting routeKeys ([]int) for that key, as well as a specified routeKey integer. It
returns all combinations of routes that include the specified route (found with routeteKey input) and do not
include any conflicting routes.
*/
func getNonConflictingCombinations(routeConflictMap map[int][]int, combination []int) [][]int {
	// Initialize a slice to hold the combinations
	var combinations [][]int
	// Add the current combination as a valid combination
	combinations = append(combinations, combination)
	// Iterate over all routes in the map
	for routeKey, conflicts := range routeConflictMap {
		// Check if the current route is not already in the combination and has no conflicts with the routes in the combination
		if !contains(combination, routeKey) && len(intersection(conflicts, combination)) == 0 {
			// If the current route is not in the combination and has no conflicts, we can add it to the combinations
			newCombination := append(combination, routeKey)
			combinations = append(combinations, newCombination)
			// Recursively generate combinations for the new combination
			subCombinations := getNonConflictingCombinations(routeConflictMap, newCombination)
			combinations = append(combinations, subCombinations...)
		}
	}
	return combinations
}

/*
compileRoute takes an input slice of slices of pointers to Room objects representing all routes routes found, as
well as a slice of integers represemting the indices of chosen routes. The function then compiles a route combination
corresponding to the input route indices, and returns the resulting route combination along with an error value. A
non-nil error is returned if any of the input indices are outside the range of the slice of all routes, or if either
of the inputs have a length of zero.
*/
func compileRoute(allRoutes [][]*sys.Room, routeIndices []int) ([][]*sys.Room, error) {
	// Check if either input slice has a length of zero
	if len(allRoutes) == 0 || len(routeIndices) == 0 {
		return nil, errors.New("\nERROR: internal malfunction, the function \" compileRoute \" given " +
			"inputs with zero length")
	}

	routeCombo := make([][]*sys.Room, 0, len(routeIndices))
	var err error

	// Iterate through the slice of route indices
	for _, index := range routeIndices {
		// Check if the index is within the range of the slice of all routes
		if index < 0 || index >= len(allRoutes) {
			return nil, errors.New("\nERROR: internal malfunction, the function \" compareRatings \" given " +
				"an input slice containing index integers outside the range of the input allRoutes variable" +
				"allRoutes length: " + strconv.Itoa(len(allRoutes)) + ", input index: " + strconv.Itoa(index))
		}
		// Append the route corresponding to the index to the compiled route
		routeCombo = append(routeCombo, allRoutes[index])
	}

	// Sort output before returning
	routeCombo, err = sortRoutes(routeCombo)
	if err != nil {
		return routeCombo, err
	}

	return routeCombo, nil
}

/*
calculateRating is a function that takes a slice of slices of pointers to Room objects representing all
routes, and a slice of integers representing indices of the routes to be considered. It returns a slice of
integers representing the number of turns and number of ant moves for the selected route combination, as well
as an error value. This error value is non-nil if any local function calls produce an error (e.g. sortRoutes,
calcAntGrouping or maxInt) or if the input slice of route indices has a length of zero.
*/
func calculateRating(allRoutes [][]*sys.Room, routeIndices []int) ([]int, error) {
	if len(routeIndices) == 0 {
		return []int{}, errors.New("\nERROR: internal malfunction, the function \" routeComboRating \" " +
			"given zero-length slice of routes as input")
	}

	// Establish working and output variables
	output := make([]int, 2)
	var routeTurn int

	// Compile routeCombo with index-specified routes
	routeCombo, err := compileRoute(allRoutes, routeIndices)
	if err != nil {
		return output, err
	}

	// Assign ants to input route
	antGrouping, err := calcAntGrouping(routeCombo)
	if err != nil {
		return []int{}, err
	}

	// Calculate ratings for the route combination (routeCombo)
	routeTurn, err = maxInt(len(routeCombo[0])-2, 1)
	if err != nil {
		return []int{}, err
	}
	nbrTurns := antGrouping[0] + routeTurn
	nbrAntMoves := 0
	for i, route := range routeCombo {
		nbrAntMoves += antGrouping[i] * (len(route) - 1)
	}
	output[0], output[1] = nbrTurns, nbrAntMoves
	return output, nil
}

/*
compareRatings compares two slices of integers and returns a boolean value and an error. The input slices,
ratingToBeTested and ratingTestedAgainst, are expected to have a length of 2. The function compares the
first element in each slice (the total number of turns) and if the first element in ratingToBeTested is
less than the first element in ratingTestedAgainst, the function returns true and a nil error. If the first
elements in both slices are equal, the function compares the second elements (the total number of ant moves)
and ff the second element in ratingToBeTested is less than the second element in ratingTestedAgainst,
the function returns true and a nil error. If none of these conditions are met, the function returns false
and a nil error. If the length of ratingTestedAgainst is 0 and the length of ratingToBeTested is 2, this is
assumed to be a "startup" condition and the function returns true and a nil error. Otherwise, if the length
of either ratingTestedAgainst or ratingToBeTested is not 2, the function returns false and a non-nil error.
*/
func compareRatings(ratingToBeTested, ratingTestedAgainst []int) (bool, error) {
	// Account for startup, where there is not yet a "best rating"
	if len(ratingTestedAgainst) == 0 && len(ratingToBeTested) == 2 {
		return true, nil
		// If invalid input, except for startup, both inputs must have a length of 2
	} else if len(ratingTestedAgainst) != 2 || len(ratingToBeTested) != 2 {
		return false, errors.New("\nERROR: internal malfunction, the function \" compareRatings \" given " +
			"input rating slices which don't both have a length of 2 (in non-startup conditions)")
	}

	if (ratingToBeTested[0] < ratingTestedAgainst[0]) || // Compare total number of turns first
		((ratingToBeTested[0] == ratingTestedAgainst[0]) && // If total number of turns are equal,
			(ratingToBeTested[1] < ratingTestedAgainst[1])) { // compare total number of ant moves
		return true, nil
	}
	return false, nil
}

/*
findBestRouteCombo is a piece of RECURSIVE BEAUTY and takes a map of routes as input, where the key is a route ID and the
value is a slice of all routes that conflict with the route specified in the key. It also takes a function calculateRating
that takes a slice of integers representing a combination of routes as input and returns a slice of integer ratings for
that combination. The function recursively iterates over all routes in the map, generating all valid routes combinations
(not conflicting) and calculates the rating for each combination using the calculateRating function, maintaining variables
for the best rated combination and best rating thus far. Finally, the function returns the best rated combination along with
an error value, which is non-nil if any local function calls (calculateRating, compareRating and compileRoute) return an error.
*/
func findBestRouteCombo(allRoutes [][]*sys.Room, routeConflictMap map[int][]int,
	calculateRating func([][]*sys.Room, []int) ([]int, error)) ([][]*sys.Room, error) {
	// Initialise working and output variables
	var err error
	var bestRouteCombo [][]*sys.Room
	var bestCombination []int
	var bestRating []int
	var rating []int
	var better bool

	// Iterate over routes in conflict map
	for routeKey := range routeConflictMap {
		// Check if route has conflicts
		// If the current route has conflicts, check all combinations that do not include any conflicting routes
		combinations := getNonConflictingCombinations(routeConflictMap, []int{routeKey})
		for _, combination := range combinations {
			// Check the rating of the current combination
			rating, err = calculateRating(allRoutes, combination)
			if err != nil {
				return bestRouteCombo, err
			}

			better, err = compareRatings(rating, bestRating)
			if err != nil {
				return bestRouteCombo, err
			}
			// If the current combination is found to be better, update the best rating and best rated combination variables
			if better {
				bestRating = rating
				bestCombination = combination
			}
		}
	}
	bestRouteCombo, err = compileRoute(allRoutes, bestCombination)
	if err != nil {
		return bestRouteCombo, err
	}
	return bestRouteCombo, nil
}

/*
fillNextValues loops through the global Routes variable, which is assumed to have already been filtered
and ordered. Each room in a path (aside from the start and end room) is then assigned its Next element
value, which is a pointer to the next room on the route. A non-nil error is returned in the event that
the room (not start or end) already has a Next value, meaning it already features in another route (conflict).
*/
func fillNextValues() error {
	for i := 0; i < len(Routes); i++ {
		for j := 1; j < len(Routes[i])-1; j++ { // Skip start and end rooms of each route
			if Routes[i][j].Next != nil {
				return errors.New("\nERROR: internal malfunction, the function \" fillNextValues \" trying to " +
					"fill Next value for a Room which has an existing Next value (exists on multiple routes)")
			}
			Routes[i][j].Next = Routes[i][j+1]
		}
	}
	return nil
}

/*
filterRoutes acts on the Routes global variable ([][]*Room), removing the longer route of
conflicting route pairs, and ordering the global Routes variable in ascending order of
length. A non-nil error is returned if an internal error is encountered in any of the
above operations.
*/
func filterRoutes(allRoutes [][]*sys.Room) ([][]*sys.Room, error) {
	if len(allRoutes) == 0 {
		return allRoutes, errors.New("\nERROR: internal malfuntion, the \" filterRoutes \" function called" +
			"while no routes recorded in global Routes variable")
	}

	// Find optimal combination of valid, non-duplicate routes
	routeConflictMap, err := createConflictMap(allRoutes)
	if err != nil {
		return allRoutes, err
	}
	Routes, err = findBestRouteCombo(allRoutes, routeConflictMap, calculateRating)
	if err != nil {
		return allRoutes, err
	}

	// Assign Next values (*.sys.Room) for all rooms on valid routes
	err = fillNextValues()
	if err != nil {
		return allRoutes, err
	}
	return Routes, nil
}

/*
moveANT takes an input route ([]*sys.Room) as well as the index of a room on the route. An ant is then
moved from this room to the next room on the route. A non-nil error is returned if an ant is not present
in the specified room, or an ant is already present in the next room in the chain.
*/
func moveAnt(route []*sys.Room, index int) error {
	// Check if room has ant
	if route[index].AntID == 0 {
		return errors.New("\nERROR: internal malfunction, input to function \" moveAnt \" is invalid" +
			"\nno ant present in specified room, with name: " + route[index].Name)
	} else if route[index+1].AntID != 0 {
		return errors.New("\nERROR: internal malfunction, input to function \" moveAnt \" is invalid" +
			"\nant already present in next room for route, with name: " + route[index+1].Name)
	}

	// Write to global CurrentTurnStr variable (string to be printed out)
	if len(CurrentTurnStr) == 0 { // If first entry, don't begin with space
		CurrentTurnStr = CurrentTurnStr + "L" + strconv.Itoa(route[index].AntID) + "-" + route[index+1].Name
	} else {
		CurrentTurnStr = CurrentTurnStr + " L" + strconv.Itoa(route[index].AntID) + "-" + route[index+1].Name
	}

	// Move ant to / from rooms
	if route[index+1].Class == "end" {
		TotalAntsFinished++
		route[index].AntID = 0
	} else {
		route[index+1].AntID, route[index].AntID = route[index].AntID, 0
	}
	return nil
}

/*
moveNew takes no input, scans the global Routes variable and places an ant on the first room after the start room
for every route where the corresponding ant counter > 0 (from global AntGrouping variable). A non-nil error
is returned if the first room of any respective route still has an ant in it (conflict).
*/
func moveNewAnts() error {
	for i, route := range Routes {
		if AntGrouping[i] != 0 && AntID <= sys.TotalAntNbr {
			if route[1].AntID != 0 {
				return errors.New("\nERROR: internal malfunction, input to function \" moveNew \" is invalid" +
					"\nant already present in route's first room, with name: " + route[1].Name)
			}

			// Write to global CurrentTurnStr variable (string to be printed out)
			if len(CurrentTurnStr) == 0 { // If first entry, don't begin with space
				CurrentTurnStr = CurrentTurnStr + "L" + strconv.Itoa(AntID) + "-" + route[1].Name
			} else {
				CurrentTurnStr = CurrentTurnStr + " L" + strconv.Itoa(AntID) + "-" + route[1].Name
			}

			// Place ant in 1st room of route
			if Routes[i][1].Class == "end" {
				TotalAntsFinished++ // If start and end room directly connected
			} else {
				Routes[i][1].AntID = AntID
			}

			AntGrouping[i]--
			// Update global counters
			if AntID == sys.TotalAntNbr {
				break
			}
			AntID++
		}
	}
	return nil
}

/*
moveExisting takes no input, but scans the global Routes variable for those rooms where ants are already
placed. These ants are then moved to the next room on their respective routes, and these moves are
recorded to the global currentTurnStr variable. A non-nil error is returned if any internal process
encounters an error during the function's execution.
*/
func moveExistingAnts() error {
	CurrentTurnStr = ""
	for i, route := range Routes {
		// Scan each respective route backwards to ensure that space is opened for forward movement of ants
		for j := len(route) - 2; j >= 1; j-- {
			if Routes[i][j].AntID != 0 {
				errMoveAnt := moveAnt(Routes[i], j)
				if errMoveAnt != nil {
					return errMoveAnt
				}
			}
		}
	}
	return nil
}

/*
executeMoves takes no input and operates on the global Routes variable, writing to individual rooms
when ants have moved in / out. It calls the local functions moveExisting and moveNew which write the
ant movements to the global CurrentTurnStr variable. This string variable is printed out to the terminal
after each successive loop within the function. A non-nil error is returned if any of the local function
calls result in an error.
*/
func executeMoves() error {
	var errExecuteMoves error
	fmt.Println()
	for TotalAntsFinished < sys.TotalAntNbr {
		errExecuteMoves = moveExistingAnts()
		if errExecuteMoves != nil {
			return errExecuteMoves
		}
		errExecuteMoves = moveNewAnts()
		if errExecuteMoves != nil {
			return errExecuteMoves
		}
		fmt.Println(CurrentTurnStr)
	}
	fmt.Println()
	return nil
}

/*
Run is a global function within the lem-in/routing package which calls several local functions
to perform a network route analysis, filtering, and ant-routeing task. It operates on the global
variables "Routes" ([][]*sys.Room) and "AntGrouping" ([]int), printing out the results of each turn
(relative ant movements) to the terminal, until completion where all ants have been successfully
routed from the start room to end room. A non-nil error is returned if any of the local functions
encounter an error during their execution.
*/
func Run() error {
	allRoutes, err := runningDFS()
	if err != nil {
		return err
	}

	Routes, err = filterRoutes(allRoutes)
	if err != nil {
		return err
	}

	AntGrouping, err = calcAntGrouping(Routes)
	if err != nil {
		return err
	}

	err = executeMoves()
	if err != nil {
		return err
	}
	return nil
}
