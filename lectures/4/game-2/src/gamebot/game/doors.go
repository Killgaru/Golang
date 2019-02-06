// doors
package game

// import "fmt"

var (
	MapOfDoors map[string]*Door //MOD
)

type Door struct {
	Name    string
	Key     *Item
	Status  bool
	inRooms [2]*Room
}

func createDoor(Name string, key *Item, status bool, r [2]*Room) {
	d := Door{
		Name:    Name,
		Key:     key,
		Status:  status,
		inRooms: r,
	}
	d.addToMOD()
	d.addToRoom()
	return
}

func (d *Door) addToMOD() {
	if _, ok := MapOfDoors[d.Name]; !ok {
		MapOfDoors[d.Name] = d
	}
	return
}

func (d *Door) addToRoom() {
	for _, r := range d.inRooms {
		r.Doors = append(r.Doors, d)
	}
	return
}

func (d *Door) ChengeStatus() {
	if d.Status {
		d.Status = false
	} else {
		d.Status = true
	}
	return
}

func (d *Door) CheckStatus() string {
	if d.Status {
		return "дверь закрыта"
	}
	return "дверь открыта"
}

func InitDoors() {
	MapOfDoors = map[string]*Door{}
	createDoor(
		"дверь",
		MapOfItems["ключи"],
		true,
		[2]*Room{
			World["домой"]["домой"]["коридор"],
			World["улица"]["улица"]["улица"],
		},
	)
	return
}
