package lemin

type room struct {
	name string
	x    int
	y    int
}

func SetRoom(room_name string, room_x int, room_y int) room {
	return room{
		name: room_name,
		x:    room_x,
		y:    room_y,
	}
}
