package room

type Room struct {
	Id   int
	Name string
}

func BuildRoom(id int, name string) Room {
	return Room{id, name}
}
