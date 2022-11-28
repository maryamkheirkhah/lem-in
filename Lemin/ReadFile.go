package lemin

import (
	"bufio"
	"fmt"
	"os"
)

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	// iterate over lines of file
	i := 0
	numberOfAnts := scanner.Text()
	fmt.Println("numberOfAnts", numberOfAnts)
	for scanner.Scan() {
		fmt.Println(i, scanner.Text())
		if scanner.Text() == "##start" {
			fmt.Println("start")
		} else if scanner.Text() == "##end" {
			fmt.Println("end")
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		file.Close()

	}

	err = file.Close() // proper close instead of defer which apparently causes bugs

	return nil, err

}
