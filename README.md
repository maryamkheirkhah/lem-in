# LEM-IN, grit:lab golang project

## 1. PROJECT OVERVIEW

This project consists of creating a a digital version of an ant farm, where the quickest way to get 'n' number of ants across a colony (composed of rooms and tunnels) needs to be found and printed to the terminal. The objectives can be outlined as follows:   

> **1.1.** At the beginning of the game, all the ants are in the room ##start. The goal is to bring them to the room ##end with as few moves as possible.
> 
> **1.2.** The shortest path is not necessarily the simplest.
> 
> **1.3.** Some colonies will have many rooms and many links, but no path between ##start and ##end.
> 
> **1.4.** Some will have rooms that link to themselves, sending your path-search spinning in circles. Some will have too many/too few ants, no ##start or ##end, duplicated rooms, links to unknown rooms, rooms with invalid coordinates and a variety of other invalid or poorly-formatted input. In those cases the program will return an error message ERROR: invalid data format. If you wish, you can elaborate a more specific error message (example: ERROR: invalid data format, invalid number of Ants or ERROR: invalid data format, no start room found).

## 2. SUMMARY OF ALGORITHM  
  
The algorithm is written entirely in GO, and adopts its own module "lem-in", with two packages: "sys" and "routing".  
  
"**sys**" contains all those system functions involved in reading, interpreting and error-checking the input. It is called via the main file using its one global function "*Setup*", which in turn calls a host of local functions that reads the input, returns errors where applicable, and otherwise compiles the ant room network and writes to Global variables which are used as input for functions within the "*routing*" package. The read order is distinct and listed below. It allows for input files to be a little 'muddled' but still imported as long as all data is included and follows the guidleines listed in *3. FORMATTING RULES / INTERPRETATION*.  
  
> **2.1.** Number of ants.  
>  
> **2.2.** Number of rooms (+ check for start & end room).  
>  
> **2.3.** Reading of rooms (names, coordinates).  
>  
> **2.4.** Reading of links (+ check for valid rooms, valid links).
  
"**routing**" contains all those functions involved in analysis of the input room network (global variable). Included in this is a ***depth-first-search*** of the network, whereby a list of all possible routes, from the start room to end room, are compiled. Each route is then "mapped" according to conflicts with all other routes. A ***recursive*** "tournament" strategy is then adopted, whereby all possible non-conflicting route combinations are explored, rated and compared to the current top-rated combination. Once the optimal route combination has been returned, the individual rooms featured in the combinations are given a value pointing to the next room in the route. This allows for each route to be used as a linked-list, enabling efficient routing of all ants through the network. This package also includes those functions involved in printing the solution to the terminal.  

## 3. FORMATTING RULES / INTERPRETATION
  
The stipulated requirements (including those for formatting) have been interpreted as follows:  

> **3.1.** Exactly one "ant number formatted line" must be included in the input file, which consists of a *positive integer*  
>  
> **3.2.** A room name must consist of alphanumeric characters only, and may not include characters such as " *-* ", " *,* " or *blank space*.
>  
> **3.3.** Exactly one *start* and one *end* room must be included in the input file, and will be interpreted as the first "room formatted line" *after* the respective " *##start* " and " *##end* " lines.
>  
> **3.4.** A "room formatted line" is interpreted as a room name, separated by *one-or-more* blank spaces, followed by two integer values which are separated by *one-or-more* blank spaces (e.g. " *room1 10 2* ").
> 
> **3.5.** A "link formatted line" consists of exactly two room names separated by a hyphen (" *-* ") with *zero-or-more* blank spaces on either side (e.g. " *room1-room2* ").  
>  
> **3.6** A link may NOT specify a room linking to itself (e.g. " *room1-room1* ").
>  
> **3.7.** Blank spaces at the beginning and end of lines are permitted but ignored, and thus are not included in room names etc.  
>  
> **3.8.** Lines beginning with a ***single*** " *#* " are considered comment lines, and are ignored.  
>  
> **3.9.** A valid input file is a *text file* (**.txt**), consisting of **only alphanumeric characters** in its file name, and with contents adhering to the guidelines listed above.

## 4. USAGE  
  
Simply clone this repository and open a terminal with the repo's main folder as its working directory. Place any desired input file in this main directory, and run the program from the terminal as follows:  

> ***go run . <name_of_input_file>***  
>  
> e.g. " *go run . audit.txt* "  

Alternatively. A directory of example files already exist within the repo (*./sys/examples*) and can be called directly by their file name without specifying their file path (e.g. " *go run . example00.txt* ").

 