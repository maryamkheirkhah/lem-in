package lemin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var FileSlice = []string{}

//var numberOfAnts int

func ReadFile(filename string) ([]string, error) {
	fileAddress := "examples/" + filename
	file, err := os.Open(fileAddress)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	// iterate over lines of file

	for scanner.Scan() {
		FileSlice = append(FileSlice, scanner.Text())
	}

	fmt.Println("file is", FileSlice)
	if err := scanner.Err(); err != nil {
		file.Close()

	}
	Makefarm()
	err = file.Close() // proper close instead of defer which apparently causes bugs

	return nil, err

}

func Makefarm() {
	//var err error
	numberOfAnts, err := strconv.Atoi(FileSlice[0])
	if err != nil {
		fmt.Println("fuck")
	}
	if FileSlice[1] == "##start" {
		start := CreateRoom(FileSlice[2])

		fmt.Println("start room is :", start)

	}
}
func CreateRoom(roomStr string) room {
	startStr := strings.Split(roomStr, " ")
	x, err1 := strconv.Atoi(startStr[1])
	if err1 != nil {
		fmt.Println(err1)
	}
	y, err2 := strconv.Atoi(startStr[2])
	if err2 != nil {
		fmt.Println(err1)
	}
	room := SetRoom(startStr[0], x, y)
	return room
}
