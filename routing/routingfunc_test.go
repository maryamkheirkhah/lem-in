package routing

import (
	"lem-in/sys"
	"reflect"
	"testing"
)

func TestFindByName(t *testing.T) {
	// Establish test variables / struct
	sys.Network = []sys.Room{{Name: "room1"}, {Name: "room2"}, {Name: "room3"},
		{Name: "room4"}, {Name: "room5"}, {Name: "room6"}}

	// Test valid input
	result := findByName("room4")
	if result == nil {
		t.Errorf("\nfunction findByName returning nil for valid input")
	} else if result.Name != "room4" {
		t.Errorf("\nfunction findByName returning incorrect room for valid input"+
			"\ngot: %v \nexpected: %v ", result.Name, "room4")
	}

	// Test invalid input
	result = findByName("room7")
	if result != nil {
		t.Errorf("\nfunction findByName not returning nil for invalid input")
	}
}

func TestCheckRouteConflict(t *testing.T) {
	// Establish test variables / structs
	sys.Network = []sys.Room{{Name: "1"}, {Name: "2"}, {Name: "3"},
		{Name: "4"}, {Name: "5"}, {Name: "6"}}
	route1 := []*sys.Room{&sys.Network[0], &sys.Network[1], &sys.Network[2], &sys.Network[5]}
	route2 := []*sys.Room{&sys.Network[0], &sys.Network[3], &sys.Network[4], &sys.Network[5]}
	route3 := []*sys.Room{&sys.Network[0], &sys.Network[2], &sys.Network[5]}
	route4 := []*sys.Room{&sys.Network[0], &sys.Network[4], &sys.Network[5]}
	route5 := []*sys.Room{&sys.Network[0], &sys.Network[5]}
	route6 := []*sys.Room{}

	conflict1, errTrue1 := checkRouteConflict(route1, route3)
	conflict2, errTrue2 := checkRouteConflict(route2, route4)
	noConflict1, errFalse1 := checkRouteConflict(route1, route2)
	noConflict2, errFalse2 := checkRouteConflict(route3, route5)
	_, errBadInput := checkRouteConflict(route1, route6)

	// Run tests / comparisons of received vs. expected
	if !conflict1 {
		t.Errorf("\nfunction checkRouteconflict not returning expected values for route conflict (1):"+
			"\nroute1: %v \nroute2: %v \nconflict: %v \nerror: %v", route1, route3, conflict1, errTrue1)
	} else if !conflict2 {
		t.Errorf("\nfunction checkRouteconflict not returning expected values for route conflict (2):"+
			"\nroute1: %v \nroute2: %v \nconflict: %v \nerror: %v", route2, route4, conflict2, errTrue2)
	} else if noConflict1 {
		t.Errorf("\nfunction checkRouteconflict not returning expected values for NO route conflict (1):"+
			"\nroute1: %v \nroute2: %v \nconflict: %v \nerror: %v", route1, route2, noConflict1, errFalse1)
	} else if noConflict2 {
		t.Errorf("\nfunction checkRouteconflict not returning expected values for NO route conflict (2):"+
			"\nroute1: %v \nroute2: %v \nconflict: %v \nerror: %v", route3, route5, noConflict2, errFalse2)
	} else if errBadInput == nil {
		t.Errorf("\nfunction checkRouteconflict not detecting invalid input:"+
			"\nroute1: %v \nroute2: %v \nerror: %v", route1, route6, errBadInput)
	}
}

func TestMoveAnt(t *testing.T) {
	// Establish test variables
	sys.Network = []sys.Room{{Name: "1", Class: "start", AntID: 0}, {Name: "2", Class: "intermediate", AntID: 4},
		{Name: "3", Class: "intermediate", AntID: 0}, {Name: "4", Class: "intermediate", AntID: 0},
		{Name: "5", Class: "intermediate", AntID: 1}, {Name: "6", Class: "end", AntID: 0}}
	Routes = [][]*sys.Room{{&sys.Network[0], &sys.Network[1], &sys.Network[2], &sys.Network[3], &sys.Network[4], &sys.Network[5]}}
	errMoveAntFalse1 := moveAnt(Routes[0], 2)
	errMoveAntFalse2 := moveAnt(Routes[0], 3)
	errMoveAntValid1 := moveAnt(Routes[0], 1)
	errMoveAntValid2 := moveAnt(Routes[0], 4) // Test moving to end room

	// Perform tests / comparisons of received vs. expected
	if errMoveAntValid1 != nil || sys.Network[2].AntID != 4 || sys.Network[1].AntID != 0 {
		t.Errorf("\nfunction moveAnt not producing expected results (1)"+
			"\ngot error: %v got room AntID (to): %v \ngot room AntID (from): %v",
			errMoveAntValid1, sys.Network[2].AntID, sys.Network[1].AntID)
	} else if errMoveAntValid2 != nil || sys.Network[4].AntID != 0 || sys.Network[5].AntID != 0 {
		t.Errorf("\nfunction moveAnt not producing expected results (2)"+
			"\ngot error: %v got room AntID (to): %v \ngot room AntID (from): %v",
			errMoveAntValid2, sys.Network[5].AntID, sys.Network[4].AntID)
	} else if errMoveAntFalse1 == nil {
		t.Errorf("\nfunction moveAnt not producing error for invalid input (1)")
	} else if errMoveAntFalse2 == nil {
		t.Errorf("\nfunction moveAnt not producing error for invalid input (2)")
	}
}

func TestMoveNew(t *testing.T) {
	// Establish test variables
	AntID = 6
	sys.TotalAntNbr = 100
	sys.Network = []sys.Room{{Name: "1", Class: "start", AntID: 0}, {Name: "2", Class: "intermediate", AntID: 4},
		{Name: "3", Class: "intermediate", AntID: 0}, {Name: "4", Class: "intermediate", AntID: 1},
		{Name: "5", Class: "intermediate", AntID: 5}, {Name: "6", Class: "end", AntID: 0}}

	// Valid input
	Routes = [][]*sys.Room{
		{&sys.Network[0], &sys.Network[1], &sys.Network[5]},
		{&sys.Network[0], &sys.Network[2], &sys.Network[5]}}
	AntGrouping = []int{0, 2}
	errMoveNewValid := moveNewAnts()

	// Invalid input
	Routes = [][]*sys.Room{
		{&sys.Network[0], &sys.Network[3], &sys.Network[5]},
		{&sys.Network[0], &sys.Network[4], &sys.Network[5]}}
	AntGrouping = []int{2, 1}
	errMoveNewFalse := moveNewAnts()

	// Perform tests / comparisons of received vs. expected
	if errMoveNewValid != nil {
		t.Errorf("\nfunction moveNewAnts producing unexpected error for valid input"+
			"\ngot: %v", errMoveNewValid)
	} else if errMoveNewFalse == nil {
		t.Errorf("\nfunction moveNewAnts not producing error for invalid input")
	}
}

func TestMoveExistingAnts(t *testing.T) {
	// Establish test variables
	sys.Network = []sys.Room{{Name: "1", Class: "start", AntID: 0}, {Name: "2", Class: "intermediate", AntID: 4},
		{Name: "3", Class: "intermediate", AntID: 3}, {Name: "4", Class: "intermediate", AntID: 2},
		{Name: "5", Class: "intermediate", AntID: 1}, {Name: "6", Class: "end", AntID: 0}}

	// Valid input
	Routes = [][]*sys.Room{
		{&sys.Network[0], &sys.Network[1], &sys.Network[3], &sys.Network[5]},
		{&sys.Network[0], &sys.Network[2], &sys.Network[4], &sys.Network[5]}}
	errMoveExistingAnts := moveExistingAnts()

	gotAntID := []int{sys.Network[0].AntID, sys.Network[0].AntID, sys.Network[0].AntID,
		sys.Network[0].AntID, sys.Network[0].AntID, sys.Network[0].AntID}
	correctAntID := []int{0, 0, 4, 3, 2, 0}

	// Perform tests / comparisons of received vs. expected
	if errMoveExistingAnts != nil {
		t.Errorf("\nfunction moveExistingAnts producing unexpected error for valid input"+
			"\ngot: %v", errMoveExistingAnts)
	} else if reflect.DeepEqual(gotAntID, correctAntID) {
		t.Errorf("\nfunction moveNewAnts not correctly changing AntID values for rooms"+
			"\ngot: %v \nexpected: %v", gotAntID, correctAntID)
	}
}

func TestExecuteMoves(t *testing.T) {
	// Establish test variables
	AntID = 1
	sys.TotalAntNbr = 99
	CurrentTurnStr = ""
	sys.Network = []sys.Room{{Name: "1", Class: "start", AntID: 0}, {Name: "2", Class: "intermediate", AntID: 0},
		{Name: "3", Class: "intermediate", AntID: 0}, {Name: "4", Class: "intermediate", AntID: 0},
		{Name: "5", Class: "intermediate", AntID: 0}, {Name: "6", Class: "end", AntID: 0}}

	// Valid input
	Routes = [][]*sys.Room{
		{&sys.Network[0], &sys.Network[1], &sys.Network[2], &sys.Network[5]},
		{&sys.Network[0], &sys.Network[3], &sys.Network[5]},
		{&sys.Network[0], &sys.Network[5]}}
	AntGrouping = []int{32, 33, 34}

	errExecuteMoves := executeMoves()

	// Perform tests / comparisons of received vs. expected
	if errExecuteMoves != nil {
		t.Errorf("\nfunction executeMoves returning unexpected error for valid input"+
			"\ngot: %v", errExecuteMoves)
	} else if AntID != sys.TotalAntNbr {
		t.Errorf("\nfunction executeMoves not altering counters corrently"+
			"\ngot TotalAntNbr: %v \ngot AntID: %v", sys.TotalAntNbr, AntID)
	}
}
