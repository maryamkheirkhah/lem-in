package main

import (
	lemin "lem-in/Lemin"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		lemin.ReadFile(args[1])
	}

}
