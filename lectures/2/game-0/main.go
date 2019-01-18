// main
package main

import (
	"fmt"
)

type Door struct {
	fromTo []string
	status bool
}

type Room struct {
	Name         string
	shortDescrip string
	description  string
	items        []Item
	doors        []Door
	neighbors    []string
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

func (r *Room) createRoom(Name, shortDes, des string, item []Item, d []Door, nei []string) {
	r.Name = Name
	r.shortDescrip = shortDes
	r.description = des
	r.items = item
	r.neighbors = nei
	r.doors = d
	r.desWalk(nei)
	RoomsInGame[Name] = r
	return
}

func (r *Room) desWalk(nei []string) {
	var s string = " можно пройти - "
	for i, v := range nei {
		if i != len(nei)-1 {
			s += fmt.Sprintf("%s, ", v)
		} else {
			s += fmt.Sprintf("%s", v)
		}

	}
	r.description += s
	r.shortDescrip += s
	return
}

func (p *Players) lookAround(...string) string {
	return p.position.description
}

func (p *Players) move(s ...string) string {
	n := p.position.neighbors
	//fmt.Println(n, s[0])
	for _, v := range n {
		//fmt.Println(n, v, s[0], s[0] == v)
		if s[0] == v {
			k := RoomsInGame[v]
			p.position = k
			return p.position.shortDescrip
		}
	}
	return "нет пути в " + s[0]
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
		return "вы одели: " + s[0]
	}
	return "не могу одеть: " + s[0]
}

func (p *Players) take(s ...string) string {
	if p.Checker(s[0]) {
		return s[0] + " - уже взято"
	}
	if p.Checker(myBackpack.Name) {
		if p.position.Checker(s[0]) {
			p.addItemFromRoom(s[0])
			return "предмет добавлен в инвентарь: " + s[0]
		}
		return "нет такого"
	}
	return "некуда класть"
}

var (
	RoomsInGame                                  = map[string]*Room{}
	kitchen, corridor, myRoom, myStreet, myHouse Room
	Player                                       = Players{
		Name:      "Player1",
		inventory: []Item{},
		position:  &kitchen,
	}
	myBackpack = Item{
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
	myDoor = Door{
		fromTo: []string{
			corridor.Name + "-" + myStreet.Name,
			myStreet.Name + "-" + myHouse.Name,
		},
		status: true,
	}
	commands = map[string]func(*Players, ...string) string{
		"осмотреться": (*Players).lookAround,
		"идти":        (*Players).move,
		"одеть":       (*Players).robe,
		"взять":       (*Players).take,
	}
)

func handleCommand(s []string) string {
	//fmt.Println(len(s))
	switch len(s) {
	case 1:
		if v, ok := commands[s[0]]; ok {
			return v(&Player)
		}
	case 2:
		if v, ok := commands[s[0]]; ok {
			return v(&Player, s[1])
		}
	}

	return "неизвестная команда"
}

func Reader() (s []string) {
	var i int
	s = make([]string, 3)
	fmt.Scanln(&s[0], &s[1], &s[2])
	for _, v := range s {
		if v != "" {
			i++
		}
	}
	s = s[:i]
	return
}

func init() {
	kitchen.createRoom(
		"кухня",
		"кухня, ничего интересного.",
		"ты находишься на кухне, на столе чай, надо собрать рюкзак и идти в универ.",
		[]Item{}, []Door{}, []string{"коридор"},
	)
	corridor.createRoom(
		"коридор",
		"ничего интересного.",
		"ничего интересного.",
		[]Item{}, []Door{}, []string{"кухня", "комната", "улица"},
	)
	myRoom.createRoom(
		"комната",
		"ты в своей комнате.",
		"на столе: ключи, конспекты, на стуле - рюкзак.",
		[]Item{myKeys, myNotes, myBackpack},
		[]Door{}, []string{"коридор"},
	)
	myStreet.createRoom(
		"улица",
		"на улице весна.",
		"на улице весна.",
		[]Item{}, []Door{}, []string{"домой"},
	)
	myHouse.createRoom(
		"домой",
		corridor.shortDescrip,
		corridor.description,
		corridor.items,
		corridor.doors,
		corridor.neighbors,
	)

}

func main() {
	fmt.Println("Hello World!")
	// fmt.Println(kitchen.neighbors)
	// fmt.Println(corridor.neighbors)
	// fmt.Println(myRoom.neighbors)
	// fmt.Println(RoomsInGame)
	fmt.Println(handleCommand([]string{"осмотреться"}))
	fmt.Println(Player.position.Name)
	fmt.Println()
	fmt.Println(handleCommand([]string{"идти", "коридор"}))
	fmt.Println(Player.position.Name)
	fmt.Println()
	fmt.Println(handleCommand([]string{"идти", "комната"}))
	fmt.Println(Player.position.Name)
	fmt.Println()
	fmt.Println(handleCommand([]string{"осмотреться"}))
	fmt.Println(Player.position.Name)
	fmt.Println(Player.position.items)
	fmt.Println(Player.inventory)
	fmt.Println()
	fmt.Println(handleCommand([]string{"взять", "телефон"}))
	fmt.Println(Player.position.Name)
	fmt.Println()
	fmt.Println(handleCommand([]string{"одеть", "рюкзак"}))
	fmt.Println(Player.position.Name)
	fmt.Println(Player.position.items)
	fmt.Println(Player.inventory)
	fmt.Println()
	fmt.Println(handleCommand([]string{"взять", "телефон"}))
	fmt.Println(Player.position.Name)
	fmt.Println(Player.position.items)
	fmt.Println(Player.inventory)
	//fmt.Println(myBackpack.description[&myRoom.Name])

}
