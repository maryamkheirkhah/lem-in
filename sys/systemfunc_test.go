package sys

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

/*
Check error detection of:
	- one valid integer for number of ants --> too many (> 100 000) / too few ants (< 1)
	- one valid ##start room and one valid ##end room
	- no duplicate room, no duplicate coordinates
	- rooms with invalid coordinates
	- rooms that link to themselves
	- invalid or poorly-formatted input
	- MUST return an error message ERROR: invalid data format, <extra info>

Check integrity of output "InputInfo" struct:
	- total number of ants is recorded
	- rooms and their properties are recorded
		- name, coordinates, class and links
*/

func TestFindRoomIndex(t *testing.T) {
	roomIndex := -1
	var err error
	e := errors.New("error")
	Network = []Room{
		{Name: "Saturn", Coords: []int{1, 2}, Links: []*Room{}},
		{Name: "Neptune", Coords: []int{1, 3}, Links: []*Room{}}}
	testNames := []string{"Saturn", "Neptune", "Pluto", "Nibiru"}
	correctIndex := []int{0, 1, -1, -1}
	correctErr := []error{nil, nil, e, e}
	for i, name := range testNames {
		roomIndex, err = findRoomIndex(name)
		if roomIndex != correctIndex[i] {
			t.Errorf("\nfunction findRoomIndex returning incorrect index"+
				"\ninput: %v \ngot: %v \nexpected: %v", name, roomIndex, correctIndex[i])
		}
		if (err != nil && correctErr[i] == nil) || (err == nil && correctErr[i] != nil) {
			t.Errorf("\nfunction findRoomIndex not returning expected error value"+
				"\ninput: %v \ngot: %v \nexpected: %v", name, err, correctErr[i])
		}
	}
}

func TestInitialiseMapKey(t *testing.T) {

	// Set global variables
	TotalRoomNbr = 3
	NetworkMap = make(map[string][]*Room, TotalRoomNbr)
	Network = []Room{
		{Name: "1", Coords: []int{1, 2}, Links: []*Room{&Network[1]}},
		{Name: "2", Coords: []int{1, 3}, Links: []*Room{&Network[0]}},
		{Name: "3", Coords: []int{1, 4}, Links: []*Room{}}}

	NetworkMap[Network[0].Name] = append(NetworkMap[Network[0].Name], &Network[1])
	NetworkMap[Network[1].Name] = append(NetworkMap[Network[0].Name], &Network[0])

	newName := Network[2].Name       // Index 2 of Network has not yet been added to NetworkMap
	duplicateName := Network[0].Name // Index 0 of Network has not yet been added to NetworkMap
	errTrue := initialiseMapKey(newName)
	errFalse := initialiseMapKey(duplicateName)

	if errTrue != nil {
		t.Errorf("\nfunction initialiseMapKey returning error for non-existent map key"+
			"\ninput: %v \nerror: %v", newName, errTrue)
	} else if _, found := NetworkMap[newName]; !found {
		t.Errorf("\nfunction initialiseMapKey failing to write new valid map key"+
			"\ninput key: %v \nresulting NetworkMap: %v", newName, NetworkMap)
	} else if errFalse == nil {
		t.Errorf("\nfunction initialiseMapKey not returning error for pre-existing map key"+
			"\nexisting: %v \ninput key: %v", NetworkMap, duplicateName)
	}
}

func TestWriteLinks(t *testing.T) {
	// Reset / initialise global variables
	TotalRoomNbr = 4
	Network = []Room{
		{Name: "1", Coords: []int{1, 2}, Links: []*Room{&Network[1]}},
		{Name: "2", Coords: []int{1, 3}, Links: []*Room{&Network[0]}},
		{Name: "3", Coords: []int{1, 4}, Links: []*Room{}},
		{Name: "4", Coords: []int{1, 5}, Links: []*Room{}}}
	NetworkMap = make(map[string][]*Room, TotalRoomNbr)
	// NetworkMap["1"] = []*Room{&Network[1]}
	// NetworkMap["2"] = []*Room{&Network[0]}

	for _, room := range Network {
		errInitialiseMapKey := initialiseMapKey(room.Name)
		if errInitialiseMapKey != nil {
			t.Fatalf("\nfunction unable to initialise map keys for WriteLinks"+
				"\nerror: %v", errInitialiseMapKey)
		}
	}

	testLinksTrue1 := []string{"2", "3"}
	testLinksTrue2 := []string{"3", "4"}
	testLinksFalse1 := []string{"2", "1"}
	testLinksFalse2 := []string{"5", "1"}
	testLinksFalse3 := []string{"1", "1"}

	errReadLinksTrue1 := writeLinks(testLinksTrue1)
	errReadLinksTrue2 := writeLinks(testLinksTrue2)
	errReadLinksFalse1 := writeLinks(testLinksFalse1)
	errReadLinksFalse2 := writeLinks(testLinksFalse2)
	errReadLinksFalse3 := writeLinks(testLinksFalse3)

	// Check valid inputs
	if errReadLinksTrue1 != nil {
		t.Errorf("\nfunction writeLinks returning error for valid inputs (1)"+
			"\ninput: %v \nerror: %v", testLinksTrue1, errReadLinksTrue1)
	} else if errReadLinksTrue2 != nil {
		t.Errorf("\nfunction writeLinks returning error for valid inputs (2)"+
			"\ninput: %v \nerror: %v", testLinksTrue2, errReadLinksTrue2)
	}

	// Check invalid inputs
	if errReadLinksFalse1 == nil {
		t.Errorf("\nfunction writeLinks failing to return an error for invalid inputs (1)"+
			"\nexisting: %v \ninput: %v", NetworkMap, testLinksFalse1)
	} else if errReadLinksFalse2 == nil {
		t.Errorf("\nfunction writeLinks failing to return an error for invalid inputs (2)"+
			"\nexisting: %v \ninput: %v", NetworkMap, testLinksFalse2)
	} else if errReadLinksFalse3 == nil {
		t.Errorf("\nfunction writeLinks failing to return an error for invalid inputs (3)"+
			"\nexisting: %v \ninput: %v", NetworkMap, testLinksFalse3)
	}
}

func TestParseLinks(t *testing.T) {
	testLinkLineTrue1, testLinkLineTrue2, testLinkLineTrue3 := "1-2", "   Bob- Joe ", "Mary  -Sarah"
	testLinkLineFalse1, testLinkLineFalse2, testLinkLineFalse3 := "-1-2", "Bob Joe-", "Mary--Sarah"
	correctLinks1, correctLinks2, correctLinks3 := []string{"1", "2"}, []string{"Bob", "Joe"}, []string{"Mary", "Sarah"}

	links1, errTrue1 := parseLinks(testLinkLineTrue1)
	links2, errTrue2 := parseLinks(testLinkLineTrue2)
	links3, errTrue3 := parseLinks(testLinkLineTrue3)
	_, errFalse1 := parseLinks(testLinkLineFalse1)
	_, errFalse2 := parseLinks(testLinkLineFalse2)
	_, errFalse3 := parseLinks(testLinkLineFalse3)

	// Check valid inputs
	if !reflect.DeepEqual(links1, correctLinks1) || errTrue1 != nil {
		t.Errorf("\nfunction parseLinks returning unexpected results for valid inputs (1)"+
			"\ngot: %v \nerror: %v \nexpected: %v", links1, errTrue1, correctLinks1)
	} else if !reflect.DeepEqual(links2, correctLinks2) || errTrue2 != nil {
		t.Errorf("\nfunction parseLinks returning unexpected results for valid inputs (2)"+
			"\ngot: %v \nerror: %v \nexpected: %v", links2, errTrue2, correctLinks2)
	} else if !reflect.DeepEqual(links3, correctLinks3) || errTrue2 != nil {
		t.Errorf("\nfunction parseLinks returning unexpected results for valid inputs (3)"+
			"\ngot: %v \nerror: %v \nexpected: %v", links3, errTrue3, correctLinks3)
	}

	// Check invalid inputs
	if errFalse1 == nil {
		t.Errorf("\nfunction parseLinks failing to return an error with invalid inputs (1)"+
			"\ninput: %v", testLinkLineFalse1)
	} else if errFalse2 == nil {
		t.Errorf("\nfunction parseLinks failing to return an error with invalid inputs (2)"+
			"\ninput: %v", testLinkLineFalse2)
	} else if errFalse3 == nil {
		t.Errorf("\nfunction parseLinks failing to return an error with invalid inputs (3)"+
			"\ninput: %v", testLinkLineFalse3)
	}
}

func TestReadLinks(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	// Initialise variables
	var fileContents []string
	var errReadFile error
	var readLinksErr error
	e := errors.New("error")
	correctErr := []error{nil, e, nil, nil, nil, nil, nil, nil, nil, nil, e, e, e, e}

	// Perform ReadLinks validity checks
	for i, file := range testFiles {
		// Reset / empty global variables
		Network = []Room{}
		NetworkMap = make(map[string][]*Room)

		fileContents, errReadFile = readFile(file.Name(), "lem-in")
		if errReadFile != nil {
			t.Fatalf("\nerror in reading %s, error returned: \n%s", file.Name(), errReadFile)
		}

		_ = readRooms(fileContents)
		readLinksErr = readLinks(fileContents)

		if readLinksErr != nil && correctErr[i] == nil {
			t.Errorf("\nfunction readLinks returning unexpected error for test file index: %v"+
				"\ngot: %v", i, readLinksErr)
		} else if readLinksErr == nil && correctErr[i] != nil {
			t.Errorf("\nfunction readRooms not returning expected error for test file index: %v", i)
		}
	}
}

func TestGetRoomName(t *testing.T) {
	roomDetailsTrue := []string{"", "Correct", "12", "34"}
	roomDetailsFalse := []string{"", "", ""}
	roomNameTrue, indexTrue, errTrue := getRoomName(roomDetailsTrue)
	_, _, errFalse := getRoomName(roomDetailsFalse)

	if roomNameTrue != "Correct" || indexTrue != 2 || errTrue != nil {
		t.Errorf("\nfunction getRoomName not returning correct name / next index"+
			"\ninput: %v \nreceived name: %v, index: %v, error: %v",
			roomDetailsTrue, roomNameTrue, indexTrue, errTrue)
	} else if errFalse == nil {
		t.Errorf("\nfunction getRoomName not detecting invalid inputs"+
			"\ninput: %v", roomDetailsFalse)
	}
}

func TestGetRoomCoords(t *testing.T) {
	roomDetailsTrue := []string{"1", "", "12", "34"}
	roomDetailsFalse1 := []string{"2", "10000000000000000000000000000000", "3"}
	roomDetailsFalse2 := []string{"3", "67.3", "45.7"}
	coords, err := getRoomCoords(roomDetailsTrue, "1", 2)
	coordsF1, falseErr1 := getRoomCoords(roomDetailsFalse1, "2", 1)
	coordsF2, falseErr2 := getRoomCoords(roomDetailsFalse2, "3", 1)
	correctCoords := []int{12, 34}

	if !reflect.DeepEqual(coords, correctCoords) {
		t.Errorf("\nfunction getRoomCoords not returning correct coordinates"+
			"\ngot: %v \nexpecting: %v", coords, correctCoords)
	} else if err != nil {
		t.Errorf("\nfunction getRoomCoords not detecting invalid inputs"+
			"\ninput: %v \nreceived error: %v", roomDetailsTrue, err)
	} else if falseErr1 == nil {
		t.Errorf("\nfunction getRoomCoords not detecting invalid inputs (1)"+
			"\ninput: %v \ngot: %v", roomDetailsFalse1, coordsF1)
	} else if falseErr2 == nil {
		t.Errorf("\nfunction getRoomCoords not detecting invalid inputs (2)"+
			"\ninput: %v \ngot: %v", roomDetailsFalse2, coordsF2)
	}
}

func TestWriteRoom(t *testing.T) {
	// Reset global variables
	TotalRoomNbr = 3
	NetworkMap = make(map[string][]*Room)

	roomName1, class1, coords1 := "Big", "start", []int{23, 5}
	roomName2, class2, coords2 := "Medium", "intermediate", []int{-45, 6}
	roomName3, class3, coords3 := "Small", "end", []int{5, 10}
	roomFalse1, classFalse1, coordsFalse1 := "False", "stArt", []int{3, 0}

	room1, err1 := writeRoom(roomName1, class1, coords1)
	room2, err2 := writeRoom(roomName2, class2, coords2)
	room3, err3 := writeRoom(roomName3, class3, coords3)
	_, errFalse := writeRoom(roomFalse1, classFalse1, coordsFalse1)

	roomTrue1 := Room{Name: roomName1, Class: class1, Coords: coords1, Links: make([]*Room, 0, TotalRoomNbr)}
	roomTrue2 := Room{Name: roomName2, Class: class2, Coords: coords2, Links: make([]*Room, 0, TotalRoomNbr)}
	roomTrue3 := Room{Name: roomName3, Class: class3, Coords: coords3, Links: make([]*Room, 0, TotalRoomNbr)}
	_, found1 := NetworkMap[roomName1]
	_, found2 := NetworkMap[roomName2]
	_, found3 := NetworkMap[roomName3]
	_, found := NetworkMap[roomFalse1]

	// Check valid inputs
	if !reflect.DeepEqual(room1, roomTrue1) {
		t.Errorf("\nfunction writeRoom not returning correct room struct (1)"+
			"\ngot: %v, error: %v \nexpecting: %v", room1, err1, roomTrue1)
	} else if !reflect.DeepEqual(room2, roomTrue2) {
		t.Errorf("\nfunction writeRoom not returning correct room struct (2)"+
			"\ngot: %v, error: %v \nexpecting: %v", room2, err2, roomTrue2)
	} else if !reflect.DeepEqual(room3, roomTrue3) {
		t.Errorf("\nfunction writeRoom not returning correct room struct (3)"+
			"\ngot: %v, error: %v \nexpecting: %v", room3, err3, roomTrue3)
	} else if !found1 || !found2 || !found3 || found {
		t.Errorf("\nfunction writeRoom not resulting in filling global NetworkMap variable"+
			"\ngot: %v", NetworkMap)
	}

	// Check invalid input
	if errFalse == nil {
		t.Errorf("\nfunction writeRoom not detecting invalid inputs"+
			"\ninput: %v, %v, %v", roomFalse1, classFalse1, coordsFalse1)
	}
}

func TestParseRoom(t *testing.T) {
	TotalRoomNbr = 3
	correctRoom1, correctRoom2, correctRoom3 := "vgh 23 0", "   BigRoom -55 100", " SmallRoom   10     10000"
	correctClass1, correctClass2, correctClass3 := "start", "intermediate", "end"
	falseRoom1, falseRoom2 := "g 0 100000000000000000000", "   BigRoom 20 100"
	falseClass := "End"
	expectedRoom1 := Room{Name: "vgh", Class: "start", Coords: []int{23, 0},
		Links: make([]*Room, 0, TotalRoomNbr)}
	expectedRoom2 := Room{Name: "BigRoom", Class: "intermediate", Coords: []int{-55, 100},
		Links: make([]*Room, 0, TotalRoomNbr)}
	expectedRoom3 := Room{Name: "SmallRoom", Class: "end", Coords: []int{10, 10000},
		Links: make([]*Room, 0, TotalRoomNbr)}

	corrRoom1, err1 := parseRoom(correctRoom1, correctClass1)
	corrRoom2, err2 := parseRoom(correctRoom2, correctClass2)
	corrRoom3, err3 := parseRoom(correctRoom3, correctClass3)
	_, errFalse1 := parseRoom(falseRoom1, correctClass1)
	_, errFalse2 := parseRoom(falseRoom2, falseClass)

	// Check return of correct inputs
	if !reflect.DeepEqual(corrRoom1, expectedRoom1) || err1 != nil {
		t.Errorf("\nfunction parseRoom not producing expected results (1):"+
			"\ninput: %v, %v"+"\ngot: %v, %v"+"\nexpected: %v, %v",
			correctRoom1, correctClass1, corrRoom1, err1, expectedRoom1, nil)
	} else if !reflect.DeepEqual(corrRoom2, expectedRoom2) || err2 != nil {
		t.Errorf("\nfunction parseRoom not producing expected results (2):"+
			"\ninput: %v, %v"+"\ngot: %v, %v"+"\nexpected: %v, %v",
			correctRoom2, correctClass2, corrRoom2, err2, expectedRoom2, nil)
	} else if !reflect.DeepEqual(corrRoom3, expectedRoom3) || err3 != nil {
		t.Errorf("\nfunction parseRoom not producing expected results (3):"+
			"\ninput: %v, %v"+"\ngot: %v, %v"+"\nexpected: %v, %v",
			correctRoom3, correctClass3, corrRoom3, err3, expectedRoom3, nil)
	}

	// Check return of false inputs
	if errFalse1 == nil {
		t.Errorf("\nfunction parseRoom not detecting invalid inputs (1)"+
			"\ninput: %v", falseRoom1)
	} else if errFalse2 == nil {
		t.Errorf("\nfunction parseRoom not detecting invalid inputs (2)"+
			"\ninput: %v, class: %v", falseRoom2, falseClass)
	}
}

func TestCheckRoomDuplicates(t *testing.T) {
	Network = []Room{
		{Name: "Joe", Coords: []int{1, 2}},
		{Name: "Bob", Coords: []int{1, 3}}}
	uniqueRoom := Room{Name: "Barry", Coords: []int{1, 4}}
	doubleRoomName := Room{Name: "Joe", Coords: []int{1, 5}}
	doubleRoomCoords := Room{Name: "Frank", Coords: []int{1, 2}}

	// Perform validity checks
	trueErr := checkRoomDuplicates(uniqueRoom)
	doubleNameErr := checkRoomDuplicates(doubleRoomName)
	doubleCoordsErr := checkRoomDuplicates(doubleRoomCoords)

	if trueErr != nil {
		t.Errorf("\nfunction checkRoomDuplicates returning error with valid input"+
			"\nexisting: %v \ninput: %v \n got: %v", Network, uniqueRoom, trueErr)
	} else if doubleNameErr == nil {
		t.Errorf("\nfunction checkRoomDuplicates not detecting duplicate room names")
	} else if doubleCoordsErr == nil {
		t.Errorf("\nfunction checkRoomDuplicates not detecting duplicate room coordinates")
	}
}

func TestReadRooms(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	// Initialise variables
	var fileContents []string
	var errReadFile error
	var readRoomsErr error
	e := errors.New("error")
	correctNetworkLength := []int{6, 17, 4, 14, 4, 6, 6, 27, 6, 6, 0, 0, 0, 0}
	correctErr := []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, e, e, e, e}

	// Perform ReadRooms validity checks
	for i, file := range testFiles {
		// Reset / empty global variables
		Network = []Room{}
		NetworkMap = make(map[string][]*Room)

		fileContents, errReadFile = readFile(file.Name(), "lem-in")
		if errReadFile != nil {
			t.Fatalf("\nerror in reading %s, error returned: \n%s", file.Name(), errReadFile)
		}

		readRoomsErr = readRooms(fileContents)

		if readRoomsErr != nil && correctErr[i] == nil {
			t.Errorf("\nfunction readRooms returning unexpected error for test file index: %v"+
				"\ngot: %v", i, readRoomsErr)
		} else if readRoomsErr == nil && correctErr[i] != nil {
			t.Errorf("\nfunction readRooms not returning expected error for test file index: %v", i)
		} else if len(Network) != correctNetworkLength[i] {
			t.Errorf("\nfunction readRooms not populating Network with expected number of rooms for test file index: %v"+
				"\n# rooms populated: %v \n# rooms expected: %v", i, len(Network), correctNetworkLength[i])
		}
	}
}

func TestCountRooms(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	// Initialise variables
	var fileContents []string
	var errReadFile error
	var errCountRooms error
	e := errors.New("error")
	correctTotalRoomNbr := []int{6, 17, 4, 14, 4, 6, 6, 27, 6, 6, 17, 0, 0, 0}
	correctErr := []error{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, e, e, e}

	// Perform CountRooms validity checks on example files
	for i, file := range testFiles {
		// Reset global variables
		Network = make([]Room, 0)
		TotalRoomNbr = 0

		fileContents, errReadFile = readFile(file.Name(), "lem-in")
		if errReadFile != nil {
			t.Fatalf("\nerror in reading %s, error returned: \n%s", file.Name(), errReadFile)
		}
		errCountRooms = countRooms(fileContents)
		if TotalRoomNbr != correctTotalRoomNbr[i] {
			t.Errorf("\nfunction countRooms not returning correct room count"+
				"\ngot: %v, expected, %v", TotalRoomNbr, correctTotalRoomNbr[i])
		} else if TotalRoomNbr != cap(Network) {
			t.Errorf("\nfunction countRooms not populating global Network variable with the correct room capacity"+
				"\nNumber of Rooms: %v, Capacity of Network, %v", TotalRoomNbr, cap(Network))
		} else if errCountRooms != nil && correctErr[i] == nil {
			t.Errorf("\nfunction countRooms producing unexpected error:"+
				"\ngot: %v", errCountRooms)
		} else if errCountRooms == nil && correctErr[i] != nil {
			t.Errorf("\nfunction countRooms not detecting errors in input:"+
				"\n%v", fileContents)
		}
	}
}

func TestReadAnts(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	// Initialise variables
	var fileContents []string
	var errReadFile error
	var errReadAnts error
	e := errors.New("error")
	correctAnts := []int{0, 20, 4, 10, 20, 4, 9, 9, 100, 1000, 0, 1000005, 0, -250}
	correctErr := []error{e, nil, nil, nil, nil, nil, nil, nil, nil, nil,
		e, e, e, e}

	// Perform ReadAnt validity checks on example files
	for i, file := range testFiles {
		TotalAntNbr = 0
		fileContents, errReadFile = readFile(file.Name(), "lem-in")
		if errReadFile != nil {
			t.Fatalf("\nerror in reading %s, error returned: \n%s", file.Name(), errReadFile)
		}
		errReadAnts = readAnts(fileContents)
		if TotalAntNbr != correctAnts[i] {
			t.Errorf("\nfunction readAnts not reading correct ant number"+
				"\ngot: %v, expected, %v", TotalAntNbr, correctAnts[i])
		} else if errReadAnts != nil && correctErr[i] == nil {
			t.Errorf("\nfunction readAnts producing unexpected error:"+
				"\ngot: %v", errReadAnts)
		} else if errReadAnts == nil && correctErr[i] != nil {
			t.Errorf("\nfunction readAnts not detecting errors in input:"+
				"\n%v", fileContents)
		}
	}
}

func TestFindDirPath(t *testing.T) {
	testDirPath1, testRootDir1 := "Users/test1/test2", "test1"
	testDirPath2, testRootDir2 := "Users/test1/test2", "test3"
	testDirPath3, testRootDir3 := "Users/test1/test2", "loooooooooooooonnnnggggteeeeeestttttt"
	newPath1, err1 := findDirPath(testDirPath1, testRootDir1)
	_, err2 := findDirPath(testDirPath2, testRootDir2)
	_, err3 := findDirPath(testDirPath3, testRootDir3)
	if newPath1 != "Users/test1" || err1 != nil {
		t.Errorf("\nfunction findDirPath failing on trimming to valid root directory")
	} else if err2 == nil {
		t.Errorf("\nfunction findDirPath failing to detect non valid root directory path (not present)")
	} else if err3 == nil {
		t.Errorf("\nfunction findDirPath failing to detect non valid root directory path " +
			"(root directory length exceeds directory path length)")
	}
}

func TestReadFile(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	// First test lengths of resulting string slices
	testSlice := make([][]string, len(testFiles))
	correctLengths := []int{16, 40, 10, 34, 11, 16, 17, 66, 17, 17, 41, 39, 38, 24}
	var errReadFile error
	for i, file := range testFiles {
		testSlice[i], errReadFile = readFile(file.Name(), "lem-in")
		if errReadFile != nil {
			t.Fatalf("\nerror in reading %s, error returned: \n%s", file.Name(), errReadFile)
		} else if len(testSlice[i]) != correctLengths[i] {
			t.Fatalf("\nerror in reading %s, incorrect length returned \ngot: %v, expected: %v",
				file.Name(), len(testSlice[i]), correctLengths[i])
		}
	}

	// Random checks of slice element values
	randomValue1, randomValue7 := testSlice[1][3], testSlice[7][22]
	correctValue1, correctValue7 := "1 7 0", "H3 5 2"
	if !reflect.DeepEqual(randomValue1, correctValue1) {
		t.Errorf("\nincorrect values returned from 4th line of %s \ngot: %s, expected: %s",
			testFiles[1].Name(), randomValue1, correctValue1)
	} else if !reflect.DeepEqual(randomValue7, correctValue7) {
		t.Errorf("\nincorrect values returned from 23rd line of %s \ngot: %s, expected: %s",
			testFiles[7].Name(), randomValue7, correctValue7)
	}
}

func TestCheckValidLines(t *testing.T) {
	e := errors.New("ERROR")
	var err error
	lines := [][]string{{"goodName 1 2"}, {"good_Name 1 2"}, {"1-2"}, {"1=2"}}
	correct := []error{nil, e, nil, e}

	// Perform checks
	for i, line := range lines {
		err = checkValidLines(line)
		if err == nil && correct[i] != nil {
			t.Errorf("\nnot returning error when error expected for input: %v", line)
		} else if err != nil && correct[i] == nil {
			t.Errorf("\nreturning error when no error expected for input: %v \ngot: %v", line, err)
		}
	}
}

func TestSetup(t *testing.T) {
	testFiles, errReadDir := os.ReadDir("./examples")
	if errReadDir != nil {
		t.Errorf("error in reading from examples directory: %v", errReadDir)
	}

	e := errors.New("error")
	correctErr := []error{e, e, nil, nil, nil, nil, nil, nil, nil, nil, e, e, e, e}

	for i, file := range testFiles {
		errSetup := Setup(file.Name())
		if (errSetup != nil && correctErr[i] == nil) ||
			(errSetup == nil && correctErr[i] != nil) {
			t.Errorf("\nfunction Setup not working as expected for file: %s", file.Name())
		}
	}
}
