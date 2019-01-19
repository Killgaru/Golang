// main
package main

import (
	"fmt"
	"strings"
)

type funcRoomDes func(*Room)

type Door struct {
	Name   string
	fromTo []*Room
	status bool
}

type Room struct {
	Name         string
	shortDescrip string
	description  string
	defaultDes   string
	items        []Item
	doors        []*Door
	neighbors    []string
	WhoInMe      []string
	funcsRoomDes []funcRoomDes
}

type Players struct {
	Name      string
	inventory []Item
	position  *Room
}

type Item struct {
	Name        string
	description map[*string]string
}

func (d *Door) chengeStatus() {
	if d.status {
		d.status = false
	} else {
		d.status = true
	}
	return
}

func (d *Door) checkStatus() string {
	if d.status {
		return "дверь закрыта"
	}
	return "дверь открыта"
}

func (r *Room) createRoom(Name, shortDes, des, defaultDes string, item []Item, d []*Door,
	nei []string, WhoInMe []string, funcsRoomDes []funcRoomDes) {
	r.Name = Name
	r.shortDescrip = shortDes
	r.description = des
	r.defaultDes = defaultDes
	r.items = item
	r.neighbors = nei
	r.doors = d
	r.WhoInMe = WhoInMe
	r.funcsRoomDes = funcsRoomDes
	r.desShortWalk()
	r.refreshRoomDes()
	RoomsInGame[Name] = r
	return
}

func (p *Players) createPlayer(Name string, inventory []Item, position *Room) {
	p.Name = Name
	p.inventory = inventory
	p.position = position
	PlayersInGame[Name] = p
	return
}

func (d *Door) createDoor(Name string, fromTo []*Room, status bool) {
	d.Name = Name
	d.fromTo = fromTo
	d.status = status
	return
}

func (r *Room) checkDoor(r2 *Room) bool {
	if len(r.doors) == 0 || len(r2.doors) == 0 {
		return false
	}
	for _, v1 := range r.doors {
		for _, v2 := range v1.fromTo {
			if v2.Name == r2.Name {
				return v1.status
			}
		}

	}
	return true
}

func (r *Room) desClean() {
	r.description = ""
	return
}

func (r *Room) desWalk() {
	var s string = " можно пройти - "
	nei := r.neighbors
	for i, v := range nei {
		if i != len(nei)-1 {
			s += fmt.Sprintf("%s, ", v)
		} else {
			s += fmt.Sprintf("%s", v)
		}

	}
	r.description += s
	return
}

func (r *Room) desShortWalk() {
	var s string = " можно пройти - "
	nei := r.neighbors
	for i, v := range nei {
		if i != len(nei)-1 {
			s += fmt.Sprintf("%s, ", v)
		} else {
			s += fmt.Sprintf("%s", v)
		}

	}
	r.shortDescrip += s
	return
}

func (r *Room) desFromItems() {
	if len(r.items) != 0 {
		var desItem, out string
		m := make(map[string]string)
		d := []string{}
		for _, v := range r.items {
			desItem = v.description[&r.Name]
			if vm, ok := m[desItem]; ok {
				vm += ", " + v.Name
				m[desItem] = vm
			} else {
				m[desItem] = v.Name
				d = append(d, desItem)
			}
		}
		for i, v := range d {
			out += v + m[v]
			if i+1 == len(d) {
				out += "."
			} else {
				out += ", "
			}
		}
		r.description += out
	} else {
		r.description += "пустая комната."
	}
	return
}

func (r *Room) desFromInvPlayer() {
	if r.Name == kitchen.Name {
		for _, v := range r.WhoInMe {
			if vm, ok := PlayersInGame[v]; ok {
				if vm.Checker(myBackpack.Name) && vm.Checker(myKeys.Name) &&
					vm.Checker(myNotes.Name) {
					r.description += "надо идти в универ."
				} else {
					r.description += "надо собрать рюкзак и идти в универ."
				}
			}
		}
	}
	return
}

func (r *Room) refreshRoomDes() {
	for _, v := range r.funcsRoomDes {
		v(r)
	}
	return
}

func (r *Room) setDefDesAsDes() {
	r.description = r.defaultDes
	return
}

func (p *Players) lookAround(...string) string {

	return p.position.description
}

func (p *Players) move(s ...string) string {
	n := p.position.neighbors
	for _, v := range n {
		if s[0] == v {
			k := RoomsInGame[v]
			if p.position.checkDoor(k) {
				return "дверь закрыта"
			}
			p.DelPlayerRoom()
			p.position = k
			p.AddPlayerRoom()
			p.position.refreshRoomDes()
			return p.position.shortDescrip
		}
	}
	return "нет пути в " + s[0]
}

func (p *Players) AddPlayerRoom() {
	p.position.WhoInMe = append(p.position.WhoInMe, p.Name)
	return
}

func (p *Players) DelPlayerRoom() {
	for i, v := range p.position.WhoInMe {
		if v == p.Name {
			p.position.WhoInMe = append(append(p.position.WhoInMe[:0],
				p.position.WhoInMe[:i]...), p.position.WhoInMe[i+1:]...)
		}
	}
	return
}

type Checker interface {
	Checker(string) bool
}

func (p *Players) Checker(s string) bool {
	for _, v := range p.inventory {
		if v.Name == s {
			return true
		}
	}
	return false
}

func (r *Room) Checker(s string) bool {
	for _, v := range r.items {
		if v.Name == s {
			return true
		}
	}
	return false
}

func (p *Players) addItemFromRoom(s string) {
	for i, v := range p.position.items {
		if v.Name == s {
			p.position.items = append(append(p.position.items[:0],
				p.position.items[:i]...), p.position.items[i+1:]...)
			//fmt.Printf("***\n%v\n***\n", p.position.items)
			p.inventory = append(p.inventory, v)
			return
		}
	}
	return
}

func (p *Players) robe(s ...string) string {
	if p.Checker(s[0]) {
		return s[0] + " - уже одето"
	}
	if p.position.Checker(s[0]) {
		p.addItemFromRoom(s[0])
		p.position.refreshRoomDes()
		return "вы одели: " + s[0]
	}
	return "не могу одеть: " + s[0]
}

func (p *Players) take(s ...string) string {
	if p.Checker(myBackpack.Name) {
		if p.position.Checker(s[0]) {
			p.addItemFromRoom(s[0])
			p.position.refreshRoomDes()
			return "предмет добавлен в инвентарь: " + s[0]
		}
		return "нет такого"
	}
	return "некуда класть"
}

func (p *Players) apply(s ...string) string {
	if p.Checker(s[0]) {
		if v, ok := keysAndDoors[s[0]]; ok && v.Name == s[1] {
			v.chengeStatus()
			return v.checkStatus()
		}
		return "не к чему применить"
	}
	return "нет предмета в инвентаре - " + s[0]
}

func (p *Players) ExitGame(...string) string {
	return "Exit"
}

var (
	notInitiated                                 bool = true
	keysAndDoors                                      = map[string]*Door{}
	RoomsInGame                                       = map[string]*Room{}
	PlayersInGame                                     = map[string]*Players{}
	kitchen, corridor, myRoom, myStreet, myHouse Room
	Player                                       Players
	myBackpack                                   = Item{
		Name: "рюкзак",
		description: map[*string]string{
			&myRoom.Name: "на стуле - ",
		},
	}
	myKeys = Item{
		Name: "ключи",
		description: map[*string]string{
			&myRoom.Name: "на столе: ",
		},
	}
	myNotes = Item{
		Name: "конспекты",
		description: map[*string]string{
			&myRoom.Name: "на столе: ",
		},
	}
	myDoor   Door
	commands = map[string]func(*Players, ...string) string{
		"осмотреться": (*Players).lookAround,
		"идти":        (*Players).move,
		"одеть":       (*Players).robe,
		"взять":       (*Players).take,
		"применить":   (*Players).apply,
		"Exit":        (*Players).ExitGame,
	}
)

func handleCommand(s1 string) string {
	//fmt.Println(len(s))
	s := strings.Split(s1, " ")
	//fmt.Printf("  HandleCommand:\n from Reader: %s\n converted: %v", s1, s)
	switch len(s) {
	case 1:
		if v, ok := commands[s[0]]; ok {
			return v(&Player)
		}
	case 2:
		if v, ok := commands[s[0]]; ok {
			return v(&Player, s[1])
		}
	case 3:
		if v, ok := commands[s[0]]; ok {
			return v(&Player, s[1], s[2])
		}
	}
	return "неизвестная команда"
}

func Reader() (s1 string) {
	// var i int
	s := make([]string, 3)
	fmt.Scanln(&s[0], &s[1], &s[2])
	for _, v := range s {
		if v != "" {
			// 		i++
			if s1 == "" {
				s1 += fmt.Sprint(v)
			} else {
				s1 += fmt.Sprint(" ", v)
			}

		}
	}
	//fmt.Println("  Reader: ", s)
	// s = s[:i]
	return
}

func initGame() {
	Player.createPlayer("Player1", []Item{}, &kitchen)
	kitchen.createRoom(
		"кухня",
		"кухня, ничего интересного.",
		"ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ.",
		"ты находишься на кухне, на столе чай, ",
		[]Item{}, []*Door{},
		[]string{"коридор"}, []string{Player.Name},
		[]funcRoomDes{(*Room).desClean, (*Room).setDefDesAsDes,
			(*Room).desFromInvPlayer, (*Room).desWalk},
	)
	corridor.createRoom(
		"коридор",
		"ничего интересного.",
		"ничего интересного.",
		"ничего интересного.",
		[]Item{}, []*Door{&myDoor},
		[]string{"кухня", "комната", "улица"}, []string{},
		[]funcRoomDes{(*Room).desClean, (*Room).setDefDesAsDes, (*Room).desWalk},
	)
	myRoom.createRoom(
		"комната",
		"ты в своей комнате.",
		"",
		"на столе: ключи, конспекты, на стуле - рюкзак.",
		[]Item{myKeys, myNotes, myBackpack}, []*Door{},
		[]string{"коридор"}, []string{},
		[]funcRoomDes{(*Room).desClean, (*Room).desFromItems, (*Room).desWalk},
	)
	myStreet.createRoom(
		"улица",
		"на улице весна.",
		"на улице весна.",
		"на улице весна.",
		[]Item{}, []*Door{&myDoor},
		[]string{"домой"}, []string{},
		[]funcRoomDes{(*Room).desClean, (*Room).setDefDesAsDes, (*Room).desWalk},
	)
	myHouse.createRoom(
		"домой",
		corridor.shortDescrip,
		corridor.description,
		corridor.defaultDes,
		corridor.items,
		corridor.doors,
		corridor.neighbors,
		corridor.WhoInMe,
		corridor.funcsRoomDes,
	)
	myDoor.createDoor("дверь", []*Room{&corridor, &myStreet, &myHouse}, true)
	keysAndDoors[myKeys.Name] = &myDoor

}

func Game0() {
	if notInitiated {
		initGame()
		notInitiated = false
		fmt.Println("Welcome in Game0")
	}
	out := handleCommand(Reader())
	if out != "Exit" {
		fmt.Println(out)
		Game0()
	}
	return
}

func main() {
	// initGame()
	// fmt.Println(PlayersInGame)
	// fmt.Println(kitchen.description)
	// fmt.Println(corridor.description)
	Game0()
}
