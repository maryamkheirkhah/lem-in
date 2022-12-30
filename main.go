package main

import (
	"lem-in/routing"
	"lem-in/sys"
	"log"
	"os"
)

/*
See README.md in main repository.
*/
func main() {
	args := os.Args
	if len(args) == 2 {
		var errLemIn error
		errLemIn = sys.Setup(args[1])
		if errLemIn != nil {
			log.Fatal(errLemIn)
		}
		errLemIn = routing.Run()
		if errLemIn != nil {
			log.Fatal(errLemIn)
		}
	} else {
		log.Fatal("\nERROR: invalid data format \n" + "please enter only one argument, " +
			"corresponding to the name of the input file")
	}
}

/*
This project is meant to make you code a digital version of an ant farm.

Create a program lem-in that will read from a file (describing the ants and the colony) given in the arguments.

Upon successfully finding the quickest path, lem-in will display the content of the file passed as argument and each move the ants make from room to room.

How does it work?

You make an ant farm with tunnels and rooms.
You place the ants on one side and look at how they find the exit.
You need to find the quickest way to get n ants across a colony (composed of rooms and tunnels).

At the beginning of the game, all the ants are in the room ##start. The goal is to bring them to the room ##end with as few moves as possible.
The shortest path is not necessarily the simplest.
Some colonies will have many rooms and many links, but no path between ##start and ##end.
Some will have rooms that link to themselves, sending your path-search spinning in circles. Some will have too many/too few ants, no ##start or ##end, duplicated rooms, links to unknown rooms, rooms with invalid coordinates and a variety of other invalid or poorly-formatted input. In those cases the program will return an error message ERROR: invalid data format. If you wish, you can elaborate a more specific error message (example: ERROR: invalid data format, invalid number of Ants or ERROR: invalid data format, no start room found).

You must display your results on the standard output in the following format :

number_of_ants
the_rooms
the_links

Lx-y Lz-w Lr-o ...

x, z, r represents the ants numbers (going from 1 to number_of_ants) and y, w, o represents the rooms names.

A room is defined by "name coord_x coord_y", and will usually look like "Room 1 2", "nameoftheroom 1 6", "4 6 7".

The links are defined by "name1-name2" and will usually look like "1-2", "2-5".

Instructions
You need to create tunnels and rooms.
A room will never start with the letter L or with # and must have no spaces.
You join the rooms together with as many tunnels as you need.
A tunnel joins only two rooms together never more than that.
A room can be linked to an infinite number of rooms and by as many tunnels as deemed necessary.
Each room can only contain one ant at a time (except at ##start and ##end which can contain as many ants as necessary).
To be the first to arrive, ants will need to take the shortest path or paths. They will also need to avoid traffic jams as well as walking all over their fellow ants.
You will only display the ants that moved at each turn, and you can move each ant only once and through a tunnel (the room at the receiving end must be empty).
The rooms names will not necessarily be numbers, and in order.
Any unknown command will be ignored.
The program must handle errors carefully. In no way can it quit in an unexpected manner.
The coordinates of the rooms will always be int.
Your project must be written in Go.
The code must respect the good practices.
It is recommended to have test files for unit testing.
*/
