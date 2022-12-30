package sys

import (
	"errors"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Room struct {
	Name    string
	Class   string
	Coords  []int // [x-coord, y-coord]
	Links   []*Room
	AntID   int // Initialised / reset to 0, to signify unoccupied. Any positive integer = occupied
	Visited bool
	Next    *Room
}

var (
	RegexFileName = regexp.MustCompile(`^.+\.txt\z`)
	RegexComment  = regexp.MustCompile(`^#{1}[^#].*\z`)
	RegexAnts     = regexp.MustCompile(`^\s*-?\d+\s*\z`)
	RegexStart    = regexp.MustCompile(`^##start\s*\z`)
	RegexEnd      = regexp.MustCompile(`^##end\s*\z`)
	RegexRoom     = regexp.MustCompile(`^\s*[a-zA-Z0-9]+\s+-?\d+\s+-?\d+\s*\z`)
	RegexLink     = regexp.MustCompile(`^\s*[a-zA-Z0-9]+\s*-\s*[a-zA-Z0-9]+\s*\z`)
	RegexInt      = regexp.MustCompile(`^\s*-?\d+\s*\z`)
	RegexIntLim   = regexp.MustCompile(`^\s*-?\d{1,10}\s*\z`) // coordinate can be a maximum of 10 digits
	RegexString   = regexp.MustCompile(`^\S+\z`)
	RegexEmpty    = regexp.MustCompile(`^\s*\z`)
	RegexLinkChar = regexp.MustCompile(`^`)
	MaxAnts       = int(1000000)
	Network       = make([]Room, 0)          // FOR ROUTING FUNCTIONS
	NetworkMap    = make(map[string][]*Room) // FOR ROUTING FUNCTIONS
	TotalRoomNbr  = int(0)                   // FOR ROUTING FUNCTIONS
	TotalAntNbr   = int(0)                   // FOR ROUTING FUNCTIONS
	Start         *Room                      // FOR ROUTING FUNCTIONS
	End           *Room                      // FOR ROUTING FUNCTIONS
)

/*
findRoomIndex takes an input room name string, and searches the global Network variable ([]Room).
It returns an integer equivalent to the index of the room with the matching Name, along with an
error value which is non-nil in the event that the room could not be found.
*/
func findRoomIndex(roomName string) (int, error) {
	for i, roomInNetwork := range Network {
		if roomInNetwork.Name == roomName {
			return i, nil
		}
	}
	return -1, errors.New("\nERROR: invalid data format, room not found")
}

/*
InitialiseMapKey takes a key string value and checks if it exists in the global NetworkMap map
variable. If it exists, a non-nil error is returned. If it doesn't exist, the key is initialised
along with a []string value with capacity equal to the TotalRoomNbr integer variable.
*/
func initialiseMapKey(mapKey string) error {
	if _, found := NetworkMap[mapKey]; !found {
		NetworkMap[mapKey] = make([]*Room, 0, (TotalRoomNbr - 1))
	} else {
		return errors.New("\nERROR: invalid data format, found duplicate room name: " +
			mapKey)
	}
	return nil
}

/*
WriteLinks reads an input []string containing the names of two linked rooms. It then writes these
links to the global Network and NetworkMap variables for each respective room name / key. A non-nil
error is returned in the event an invalid / non-existent room name is given, if not exactly
two room names are provided in the input []string, if the two room names are the same (room links
to itself) or if the link already exists (duplicate).
*/
func writeLinks(roomLinks []string) error {
	// Initial input error checks
	if len(roomLinks) != 2 {
		return errors.New("\nERROR: invalid data format, too many / few links provided in link input: " +
			"\ninput: " + "[ " + strings.Join(roomLinks, " , ") + " ]")
	} else if roomLinks[0] == roomLinks[1] {
		return errors.New("\nERROR: invalid data format, room connects to itself in link input: " +
			"\ninput: " + "[ " + strings.Join(roomLinks, " , ") + " ]")
	} else if _, found := NetworkMap[roomLinks[0]]; !found {
		return errors.New("\nERROR: invalid data format, input contains non-existent room name: " +
			"\nroom name not found: " + "[ " + roomLinks[0] + " ]")
	} else if _, found := NetworkMap[roomLinks[1]]; !found {
		return errors.New("\nERROR: invalid data format, input contains non-existent room name: " +
			"\nroom name not found: " + "[ " + roomLinks[1] + " ]")
	}

	roomIndex := -1
	var errFindRoom error
	// Write links to global variables
	for i, roomInNetwork := range Network {
		if roomInNetwork.Name == roomLinks[0] {
			// WRITE LINKS FOR 1ST ROOM IN SLICE
			roomIndex, errFindRoom = findRoomIndex(roomLinks[1])
			if errFindRoom != nil {
				return errFindRoom
			}
			// Write to global Network variable
			for _, link := range roomInNetwork.Links {
				if link.Name == roomLinks[1] {
					return errors.New("\nERROR: invalid data format, link already exists in global Network variable (1):" +
						"\ninput link: " + "[ " + roomLinks[0] + " , " + roomLinks[1] + " ]")
				}
			}
			Network[i].Links = append(roomInNetwork.Links, &Network[roomIndex])
			// Write to global NetworkMap variable
			for _, link := range NetworkMap[roomLinks[0]] {
				if link.Name == roomLinks[1] {
					return errors.New("\nERROR: invalid data format, link already exists in global NetworkMap variable (1):" +
						"\ninput link: " + "[ " + roomLinks[0] + " , " + roomLinks[1] + " ]")
				}
			}
			NetworkMap[roomLinks[0]] = append(NetworkMap[roomLinks[0]], &Network[roomIndex])

		} else if roomInNetwork.Name == roomLinks[1] {
			// WRITE LINKS FOR 2ND ROOM IN SLICE
			roomIndex, errFindRoom = findRoomIndex(roomLinks[0])
			if errFindRoom != nil {
				return errFindRoom
			}
			// Write to global Newtork variable
			for _, link := range roomInNetwork.Links {
				if link.Name == roomLinks[0] {
					return errors.New("\nERROR: invalid data format, link already exists in global Network variable (2):" +
						"\ninput link: " + "[ " + roomLinks[0] + " , " + roomLinks[1] + " ]")
				}
			}
			Network[i].Links = append(roomInNetwork.Links, &Network[roomIndex])
			// Write to global NetworkMap variable
			for _, link := range NetworkMap[roomLinks[1]] {
				if link.Name == roomLinks[0] {
					return errors.New("\nERROR: invalid data format, link already exists in global NetworkMap variable (2):" +
						"\ninput link: " + "[ " + roomLinks[0] + " , " + roomLinks[1] + " ]")
				}
			}
			NetworkMap[roomLinks[1]] = append(NetworkMap[roomLinks[1]], &Network[roomIndex])
		}

	}
	return nil
}

/*
ParseLinks reads an input string and parses it for linked room values (sub-strings). If the
format does not correspond to expected values (exactly two distinct strings), then a non-nil
error is returned.
*/
func parseLinks(linkLine string) ([]string, error) {
	output := []string{}
	startIndex := 0
	start := false
	foundHyphen := false
	if RegexLink.MatchString(linkLine) {
		for i, letter := range linkLine {
			if foundHyphen && letter == '-' {
				return output, errors.New("\nERROR: invalid data format, more than one \" - \" discovered in input line" +
					"\nrooms may not have \" - \" in their name, nor may multiple \" - \" be used for link inputs")
			} else if !start && letter != '-' && letter != ' ' {
				// Detect link start
				start = true
				startIndex = i
				if i == len(linkLine)-1 {
					// End of input
					// Account for single character room names
					output = append(output, linkLine[startIndex:])
					break
				}
			} else if start && (letter == '-' || letter == ' ') {
				// Detect link end, record link
				output = append(output, linkLine[startIndex:i])
				start = false
				if letter == '-' {
					foundHyphen = true
				}
			} else if !start && letter == '-' {
				foundHyphen = true
			} else if start && i == len(linkLine)-1 {
				// End of input
				output = append(output, linkLine[startIndex:])
				break
			}
		}
	} else {
		// Regex match fails
		return output, errors.New("\nERROR: invalid data format, link input poorly formatted" + "\ninput: " + linkLine)
	}
	if len(output) != 2 || !foundHyphen {
		// If output has not been filled correctly
		return output, errors.New("\nERROR: invalid data format, incorrect link input format" +
			"\nexactly 2 valid input room names, seperated by a hyphen, required in input string" +
			"\ninput: " + linkLine)
	}
	return output, nil
}

/*
ReadLinks reads file contents in the form of an input slice of strings and checks the data for the
specified room linkages. ReadLinks writes valid links to the global Network and NetworkMap variables.
If an error in the input is found it is returned. Otherwise a nil value is returned.
*/
func readLinks(fileContents []string) error {
	linkCounter := 0
	for _, line := range fileContents {
		if RegexLink.MatchString(line) {
			linkCounter++
			linkSlice, errParseLinks := parseLinks(line)
			if errParseLinks != nil {
				return errParseLinks
			}
			errWriteLinks := writeLinks(linkSlice)
			if errWriteLinks != nil {
				return errWriteLinks
			}
		}
	}
	if linkCounter < 1 {
		return errors.New("\nERROR: invalid data format, no link input data found")
	}
	return nil
}

/*
GetRoomName takes an input []string of room details and returns the first valid, non-empty string
element as the room name, as well as an integer for the next index of the slice after the name.
In the event that a valid room name cannot be found / matched, a non-nil error is returned.
*/
func getRoomName(roomDetails []string) (string, int, error) {
	for i, detail := range roomDetails {
		if RegexString.MatchString(detail) {
			return detail, i + 1, nil
		}
	}
	// If no room name is found / matched
	return "", 0, errors.New("\nERROR: invalid data format, found entry with no room name")
}

/*
GetRoomCoords takes an input []string of room details, an already parsed room name string as well
as an index integer of the element after the room name in the []string. It returns its x- and y-
coordinates as a []int. In the event that valid coordinates cannot be found / matched, a non-nil
error is returned.
*/
func getRoomCoords(roomDetails []string, roomName string, index int) ([]int, error) {
	var coords []int
	nbr := 0
	var errCoord error

	for i := index; i < len(roomDetails); i++ {
		// Check for empty / whitespace element
		if RegexEmpty.MatchString(roomDetails[i]) {
			continue
		} else if RegexInt.MatchString(roomDetails[i]) {
			// Check for possible overflow with excessively long coordinate
			if !RegexIntLim.MatchString(roomDetails[i]) {
				return coords, errors.New("\nERROR: invalid data format, coordinates may not exceed 10 digits \ngot:  " + roomDetails[i])
			}
			// Record coordinates
			nbr, errCoord = strconv.Atoi(roomDetails[i])
			if errCoord != nil {
				return coords, errors.New("\nERROR: invalid data format, problem in parsing coordinates for room \" " + roomName + " \"")
			} else {
				coords = append(coords, nbr)
			}
		}
	}
	if len(coords) != 2 {
		return coords, errors.New("\nERROR: invalid data format, problem in parsing coordinates for room \" " + roomName + " \"")
	}
	return coords, nil
}

/*
WriteRoom takes file contents as an input slice of strings, along with a starting index integer where
room details are expected, as well as the expected class of room (start, end, intermediate). The
function also writes to the global NetworkMap map variable, initialising the room name key with
a string slice with capacity equal to the global TotalRoomNbr integer variable.
*/
func writeRoom(roomName, roomClass string, roomCoords []int) (Room, error) {
	output := Room{}

	if roomClass != "start" && roomClass != "intermediate" && roomClass != "end" {
		return output, errors.New("\nERROR: invalid data format, room must either have class " +
			"<start>, <intermediate>, or <end>")
	}
	output.Name = roomName
	output.Class = roomClass
	output.Coords = roomCoords
	output.AntID = 0
	output.Next = nil

	// Initialise Links element with maximum capacity so as to
	// avoid slice appending errors later
	output.Links = make([]*Room, 0, TotalRoomNbr-1)

	// Add to global NetworkMap variable with an empty slice value with capacity
	// equal to total number of rooms - 1 (theoretical max number of links per room)
	errInitialiseMapKey := initialiseMapKey(output.Name)
	if errInitialiseMapKey != nil {
		return output, errInitialiseMapKey
	}
	return output, nil
}

/*
ParseRoom scans a formatted string (which has already passed a regexp check in Read Rooms), as well as
its class string ("start", "intermediate" or "end"). It then calls functions GetRoomName and GetRoomCoords
to parse the relevant data, and feeds the outputs to a Room structure which is then returned. If any errors
are returned by the internal function calls, these are returned.
*/
func parseRoom(roomLine, roomClass string) (Room, error) {
	output := Room{}
	roomDetails := strings.Split(roomLine, " ")
	var roomCoords []int

	if !RegexRoom.MatchString(roomLine) {
		return output, errors.New("\nERROR: invalid data format, error with entry: " + roomLine)
	}

	// Extract room name
	roomName, index, errParse := getRoomName(roomDetails)
	if errParse != nil {
		return output, errParse
	}

	// Extract room coordinates
	roomCoords, errParse = getRoomCoords(roomDetails, roomName, index)
	if errParse != nil {
		return output, errParse
	}

	// Write data to Room struct
	output, errParse = writeRoom(roomName, roomClass, roomCoords)
	if errParse != nil {
		return output, errParse
	}
	return output, nil
}

/*
CheckRoomDuplicates takes a Room struct as input and compares its elements for duplicates in
the global Network variable ([]Room). In the event of a duplicate name or coordinates, a
non-nil error is returned.
*/
func checkRoomDuplicates(inputRoom Room) error {
	for _, existingRoom := range Network {
		if reflect.DeepEqual(existingRoom.Name, inputRoom.Name) {
			return errors.New("\nERROR: invalid data format, duplicate room names detected: " + inputRoom.Name)
		} else if reflect.DeepEqual(existingRoom.Coords, inputRoom.Coords) {
			return errors.New("\nERROR: invalid data format, duplicate room coordinates detected for rooms " +
				existingRoom.Name + " and " + inputRoom.Name)
		}
	}
	return nil
}

/*
ReadRooms reads file contents in the form of an input slice of strings and checks the data for the
ant colony rooms specified. Properties of the rooms are written to the global Network struct whilst
also checking for errors. If an error in the input is found it is returned. Otherwise a nil value
is returned. Errors
*/
func readRooms(fileContents []string) error {
	startLabel := false
	endLabel := false
	roomEntry := Room{}
	var errRead error
	var errDuplicates error

	for _, line := range fileContents {
		// Check if comment-line or label
		if RegexComment.MatchString(line) || RegexEmpty.MatchString(line) {
			continue
		} else if RegexStart.MatchString(line) && !startLabel {
			startLabel = true
			continue
		} else if RegexEnd.MatchString(line) && !endLabel {
			endLabel = true
			continue
		}

		// Check if end / start labels active at the same time (no room entry inbetween)
		if startLabel && endLabel {
			Network = []Room{} // empty / reset global Network variable
			return errors.New("\nERROR: invalid data format, no room entries discovered between start and end room labels")
		} else if RegexRoom.MatchString(line) {
			// Parse and write room data to global Network variable
			if startLabel {
				roomEntry, errRead = parseRoom(line, "start")
				startLabel = false
			} else if endLabel {
				roomEntry, errRead = parseRoom(line, "end")
				endLabel = false
			} else if RegexRoom.MatchString(line) {
				roomEntry, errRead = parseRoom(line, "intermediate")
			}
			// Return error if discovered
			if errRead != nil {
				Network = []Room{} // empty / reset global Network variable
				return errRead
			}
			// Check for duplicates
			errDuplicates = checkRoomDuplicates(roomEntry)
			if errDuplicates != nil {
				Network = []Room{} // empty / reset global Network variable
				return errDuplicates
			}
			Network = append(Network, roomEntry)

			// Write Start / End rooms
			if roomEntry.Class == "start" {
				Start = &Network[len(Network)-1]
			} else if roomEntry.Class == "end" {
				End = &Network[len(Network)-1]
			}
		}
	}

	if startLabel || endLabel {
		Network = []Room{} // empty / reset global Network variable
		return errors.New("\nERROR: invalid data format, missing start and / or end room data entry")
	}
	return nil
}

/*
CountRooms takes file contents as an input slice of strings, and counts the number of
lines which correspond to a room-and-coordinate entry (e.g. RoomA 3 2). It writes this
count to the global TotalRoomNbr integer variable, whilst also writing over the global
Network ([]Room) variable, giving it a max capacity equivalent to the preliminary room
total (which may be even less depending on duplicates etc.). The global NetworkMap
variable is also reinitialised to avoid later "assignment to entry in nil map".
Finally, the function returns an error value, which is not nil in the event of less than
2 room-coordinate entries being found, or too little / too many start & end room labels
(##start / ##end).
*/
func countRooms(fileContents []string) error {
	count := 0
	startLabel := false
	endLabel := false

	for _, line := range fileContents {
		if RegexStart.MatchString(line) && !startLabel {
			startLabel = true
		} else if RegexStart.MatchString(line) && startLabel {
			return errors.New("\nERROR: invalid data format, multiple start room labels ( ##start ) detected")
		} else if RegexEnd.MatchString(line) && !endLabel {
			endLabel = true
		} else if RegexEnd.MatchString(line) && endLabel {
			return errors.New("\nERROR: invalid data format, multiple end room labels ( ##end ) detected")
		} else if RegexRoom.MatchString(line) {
			count++
		}
	}

	// Validity checks
	if !startLabel {
		return errors.New("\nERROR: invalid data format, no start room label ( ##start ) found")
	} else if !endLabel {
		return errors.New("\nERROR: invalid data format, no end room label ( ##end ) found")
	} else if count < 2 {
		// Theoretically, at least 2 rooms required ("##start" and "##end")
		return errors.New("\nERROR: invalid data format, less than 2 valid room entries in input \ngot: " + strconv.Itoa(count))
	}
	Network = make([]Room, 0, count)
	TotalRoomNbr = count
	NetworkMap = make(map[string][]*Room, TotalRoomNbr)
	return nil
}

/*
ReadAnts takes file contents as an input slice of strings, and returns the number of
ants specified in the file, whilst also checking for errors in the data input. If an error
is found, it is returned. Otherwise the ant total as an integer is returned, along with a
nil error value.
*/
func readAnts(fileContents []string) error {
	var errReadAnts error
	foundAnts := false

	for _, line := range fileContents {
		// Check if comment-line
		if RegexComment.MatchString(line) {
			continue
		}
		// Perform validity check on ants
		if !foundAnts && RegexAnts.MatchString(line) {
			// Record ant total if found
			foundAnts = true
			TotalAntNbr, errReadAnts = strconv.Atoi(line)
			if errReadAnts != nil {
				return errors.New("\nERROR: invalid data format, unknown fault in ant number format")
			}
			continue
		} else if foundAnts && RegexAnts.MatchString(line) {
			// Keep scanning input in case multiple ant totals are specified
			TotalAntNbr = 0
			return errors.New("\nERROR: invalid data format, multiple ant inputs detected")
		}
	}
	if TotalAntNbr <= 0 {
		return errors.New("\nERROR: invalid data format, number of ants must be a positive integer")
	} else if TotalAntNbr > MaxAnts {
		return errors.New("\nERROR: invalid data format, maximum number of ants ( " + strconv.Itoa(MaxAnts) + " ) exceeded" +
			"\nfound: " + strconv.Itoa(TotalAntNbr))
	} else if !foundAnts {
		return errors.New("\nERROR: invalid data format, no ant input found")
	}
	return nil
}

/*
FindDirPath takes a full "folderPath" string and trims the end to the final instance of an input
root directory "rootDir" string. It thus does not account for sub-directories with the same name
as the root directory. i.e. if folder path directory is "folder1/folder2/folder2", and the
root directory is "folder2", the folder path directory will not be trimmed. An error value is
returned in the event that the root directory is not found.
*/
func findDirPath(folderPath string, rootDir string) (string, error) {
	for i := len(folderPath) - len(rootDir); i >= 0; i-- {
		if folderPath[i:i+len(rootDir)] == rootDir {
			folderPath = folderPath[0 : i+len(rootDir)]
			return folderPath, nil
		}
	}
	return folderPath, errors.New("\nERROR: invalid data format, the requested root directory " +
		"could not be found from the input file path")
}

/*
ReadFile takes a "fileName" string, ans well as root directory ("rootDir") string and writes the
target file (with "fileName") to a slice of strings. It checks if the input file is an "example" file,
which in turn is found in "./system/examples/", otherwise the file is assumed to lie within the root
*/
func readFile(fileName string, rootDir string) ([]string, error) {
	var fileContents []string

	// Establish filepath / root directory
	// (in case being accessed by program outside of root directory)
	currentDir, errWD := os.Getwd()
	if errWD != nil {
		return fileContents, errWD
	}
	rootDir, errRD := findDirPath(currentDir, rootDir)
	if errRD != nil {
		return fileContents, errRD
	}

	// Change filename / path if file is an example / test-file
	regexExample := regexp.MustCompile(`^((example)|(badexample)|(poorexample))\d*(\.txt)\z`)
	if regexExample.MatchString(fileName) {
		fileName = rootDir + "/sys/examples/" + fileName
	}

	// Write file contents to data structure
	file, errReadFile := os.ReadFile(fileName)
	if errReadFile != nil {
		return fileContents, errors.New("\nERROR: invalid data format, the specified file could not be read / found")
	} else if len(file) == 0 {
		return fileContents, errors.New("\nERROR: invalid data format, the input file is empty")
	}
	data := strings.ReplaceAll(string(file), "\r\n", "\n")
	fileContents = strings.Split(data, "\n")
	return fileContents, nil
}

/*
checkValidLines takes an input slice of strings and checks that each line conforms to at least one
formatting standard for a valid input, ie. is a valid ant number format, or a valid room format, or
a valid link format, or a valid comment / title line (start / end room). The function returns an
error value, which is non-nil if the file contains a line which does not confirm to the aforementioned
formatting guidlines.
*/
func checkValidLines(fileContents []string) error {
	for _, line := range fileContents {
		if !RegexAnts.MatchString(line) && !RegexComment.MatchString(line) &&
			!RegexEmpty.MatchString(line) && !RegexRoom.MatchString(line) &&
			!RegexEnd.MatchString(line) && !RegexStart.MatchString(line) &&
			!RegexLink.MatchString(line) {
			return errors.New("\nERROR: invalid data format, the specified file contains lines " +
				"with incorrect formatting, eg.: " + line)
		}
	}
	return nil
}

/*
Setup is a global function which takes a file name as an input and calls the local sys functions
to process / parse the file and populate the relevant global variables (Network, NetworkMap,
TotalRoomNbr & TotalAntNbr) to be used by functions in the lem-in/colony package. A non-nil
error is returned if any errors with the input are found.
*/
func Setup(fileName string) error {
	if !RegexFileName.MatchString(fileName) {
		return errors.New("\nERROR: invalid data format, input file must be a valid .txt file")
	}

	fileContents, readFileErr := readFile(fileName, "lem-in")
	if readFileErr != nil {
		return readFileErr
	}
	readAntsErr := readAnts(fileContents)
	if readAntsErr != nil {
		return readAntsErr
	}
	countRoomsErr := countRooms(fileContents)
	if countRoomsErr != nil {
		return countRoomsErr
	}
	readRoomsErr := readRooms(fileContents)
	if readRoomsErr != nil {
		return readRoomsErr
	}
	readLinksErr := readLinks(fileContents)
	if readLinksErr != nil {
		return readLinksErr
	}
	generalErr := checkValidLines(fileContents)
	if generalErr != nil {
		return generalErr
	}
	return nil
}
