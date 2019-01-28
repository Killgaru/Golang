// players
package main

import (
	"fmt"
	"strings"
)

var (
	// Player        Players
	PlayersInGame = map[string]*Player{}
)

type thinkP struct {
	keyFunc func(*Player, interface{}) bool
	argFunc interface{}
	answers map[bool]string
}

type Player struct {
	Name      string
	inventory []*Item
	thoughts  map[*Room]thinkP
	position  *Room
	output    chan string
}

func NewPlayer(name string) *Player {
	var p Player
	p.Name = name
	if p.output != nil {
		close(p.output)
	}
	p.output = make(chan string)
	PlayersInGame[name] = &p
	return &p
}

func (p *Player) createPlayer(inventory []*Item,
	thoughts map[*Room]thinkP, position *Room) {
	p.inventory = inventory
	p.thoughts = thoughts
	p.position = position
	p.AddPlayerRoom()
	return
}

func addPlayer(p *Player) {
	p.createPlayer(
		[]*Item{},
		map[*Room]thinkP{
			World["домой"]["домой"]["кухня"]: thinkP{
				keyFunc: (*Player).inInventory,
				argFunc: []*Item{
					MapOfItems["рюкзак"],
					MapOfItems["ключи"],
					MapOfItems["конспекты"],
				},
				answers: map[bool]string{
					false: "надо собрать рюкзак и идти в универ.",
					true:  "надо идти в универ.",
				},
			},
		},
		World["домой"]["домой"]["кухня"])
	return
}

func (p *Player) AddPlayerRoom() {
	if _, ok := p.position.WhoInMe[p.Name]; !ok {
		p.position.WhoInMe[p.Name] = p
	}
	return
}

func (p *Player) DelPlayerRoom() {
	if _, ok := p.position.WhoInMe[p.Name]; ok {
		delete(p.position.WhoInMe, p.Name)
	}
	return
}

func (p *Player) addItemFromRoom(s string) {
	for i, v := range p.position.items {
		if v.Name == s {
			p.position.items = DelElemSliceItem(p.position.items, i)
			p.inventory = append(p.inventory, v)
			return
		}
	}
	return
}

func (p *Player) inInventory(it interface{}) bool {
	var out, ok bool = true, false
	val := it.([]*Item)
	for _, v := range val {
		_, ok = p.Checker(v.Name)
		out = out && ok
	}
	return out
}

func (t *thinkP) checkThink(p *Player) string {
	return t.answers[t.keyFunc(p, t.argFunc)]
}

func (p *Player) checkMyThink(r *Room) string {
	if v, ok := p.thoughts[r]; ok {
		return v.checkThink(p)
	}
	return ""
}

func (p *Player) addMyThink() (out string) {
	s := p.checkMyThink(p.position)
	if s != "" {
		ss := strings.Split(p.position.description, ". ")
		if len(ss) < 2 {
			out = s + " " + p.position.description
		} else {
			for i := 0; i < len(ss)-1; i++ {
				if i == len(ss)-2 {
					out += ss[i] + ", "
				} else {
					out += ss[i] + ". "
				}

			}
			out += s + " " + ss[len(ss)-1]
		}
	} else {
		out = p.position.description
	}
	return
}

func (p *Player) Checker(s string) (*Item, bool) {
	for _, v := range p.inventory {
		if v.Name == s {
			return v, true
		}
	}
	return &Item{}, false
}

func (p *Player) applyKeyDoor(it *Item, s string) string {
	for _, d := range p.position.doors {
		if s == d.Name && it == d.Key {
			d.ChengeStatus()
			return d.CheckStatus()
		}
	}
	return "не к чему применить"
}

func (p *Player) playersInRoom() (out string) {
	var toggle bool
	pir := p.position.WhoInMe //Players In Room
	if len(pir) > 1 {
		out = ". Кроме вас тут ещё "
		for s := range pir {
			if s != p.Name {
				if toggle {
					out += ", "
				}
				out += s
				toggle = true
			}
		}
	}
	return
}

func (p *Player) moveTo(r *Room) (out string) {
	if p.position.checkDoor(r) {
		out = "дверь закрыта"
		return
	}
	p.DelPlayerRoom()
	p.position = r
	p.AddPlayerRoom()
	out = p.position.shortDescrip
	return
}

func (p *Player) LookAround(...string) {
	out := p.addMyThink() + p.playersInRoom()
	p.output <- out
	return
}

func (p *Player) Move(s ...string) {
	var out string
	n := p.position.Neighbors
	for _, v := range n {
		if p.position.Location != v.Location {
			if s[0] == v.Location {
				out = p.moveTo(v)
				break
			} else {
				out = "нет пути в " + s[0]
			}
		} else {
			if s[0] == v.Name {
				out = p.moveTo(v)
				break
			} else {
				out = "нет пути в " + s[0]
			}
		}
	}
	p.output <- out
	return
}

func (p *Player) Robe(s ...string) {
	var out string
	_, ok := p.Checker(s[0])
	switch {
	case ok:
		out = s[0] + " - уже одето"
	case p.position.Checker(s[0]):
		p.addItemFromRoom(s[0])
		p.position.refreshRoomDes()
		out = "вы одели: " + s[0]
	default:
		out = "не могу одеть: " + s[0]
	}
	p.output <- out
	return
}

func (p *Player) Take(s ...string) {
	var out string
	if _, ok := p.Checker(MapOfItems["рюкзак"].Name); ok {
		if p.position.Checker(s[0]) {
			p.addItemFromRoom(s[0])
			p.position.refreshRoomDes()
			out = "предмет добавлен в инвентарь: " + s[0]
		} else {
			out = "нет такого"
		}
	} else {
		out = "некуда класть"
	}
	p.output <- out
	return
}

func (p *Player) Apply(s ...string) {
	var out string
	if it, ok := p.Checker(s[0]); ok {
		if it.iType == "key" {
			out = p.applyKeyDoor(it, s[1])
		}
	} else {
		out = "нет предмета в инвентаре - " + s[0]
	}
	p.output <- out
	return
}

func (p *Player) ExitGame(...string) {
	p.output <- "Exit"
	return
}

func (p *Player) SayRoom(ss ...string) {
	out := fmt.Sprintf("%s говорит:", p.Name)
	for _, v := range ss {
		out += " " + v
	}
	for _, pir := range p.position.WhoInMe {
		pir.output <- out
	}
	return
}

func (p *Player) SayPlayer(ss ...string) {
	var out string
	if p2, ok := p.position.WhoInMe[ss[0]]; ok {
		if len(ss) > 1 {
			out = fmt.Sprintf("%s говорит вам:", p.Name)
			for _, v := range ss[1:] {
				out += " " + v
			}
		} else {
			out = fmt.Sprintf("%s выразительно молчит, смотря на вас", p.Name)
		}
		p2.output <- out
	} else {
		out = "тут нет такого игрока"
		p.output <- out
	}
	return
}

func (p *Player) GetOutput() chan string {
	return p.output
}

func (p *Player) HandleInput(s string) {
	handleCommand(p, s)
	return
}

func (p *Player) HandleOutput() string {
	return <-p.GetOutput()
}

func InitPlayers() {
	PlayersInGame = map[string]*Player{}
	return
}
