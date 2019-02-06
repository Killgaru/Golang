// players
package game

import (
	"fmt"
	"strings"
	"time"
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
	Name        string
	ID          int64
	ComFragment string
	TimeLife    *time.Timer
	Inventory   []*Item
	thoughts    map[*Room]thinkP
	Position    *Room
	output      chan string
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

func (p *Player) createPlayer(id int64, inventory []*Item,
	thoughts map[*Room]thinkP, position *Room) {
	p.ID = id
	p.Inventory = inventory
	p.thoughts = thoughts
	p.Position = position
	p.AddPlayerRoom()
	return
}

func AddPlayer(p *Player, id int64) {
	p.createPlayer(
		id,
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
	if _, ok := p.Position.WhoInMe[p.Name]; !ok {
		p.Position.WhoInMe[p.Name] = p
	}
	return
}

func (p *Player) DelPlayerRoom() {
	if _, ok := p.Position.WhoInMe[p.Name]; ok {
		delete(p.Position.WhoInMe, p.Name)
	}
	return
}

func (p *Player) addItemFromRoom(s string) {
	for i, v := range p.Position.Items {
		if v.Name == s {
			p.Position.Items = DelElemSliceItem(p.Position.Items, i)
			p.Inventory = append(p.Inventory, v)
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
	s := p.checkMyThink(p.Position)
	if s != "" {
		ss := strings.Split(p.Position.description, ". ")
		if len(ss) < 2 {
			out = s + " " + p.Position.description
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
		out = p.Position.description
	}
	return
}

func (p *Player) Checker(s string) (*Item, bool) {
	for _, v := range p.Inventory {
		if v.Name == s {
			return v, true
		}
	}
	return &Item{}, false
}

func (p *Player) applyKeyDoor(it *Item, s string) string {
	for _, d := range p.Position.Doors {
		if s == d.Name && it == d.Key {
			d.ChengeStatus()
			return d.CheckStatus()
		}
	}
	return "не к чему применить"
}

func (p *Player) playersInRoom() (out string) {
	var toggle bool
	pir := p.Position.WhoInMe //Players In Room
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
	if p.Position.checkDoor(r) {
		out = "дверь закрыта"
		return
	}
	p.DelPlayerRoom()
	p.Position = r
	p.AddPlayerRoom()
	out = p.Position.ShortDescrip
	return
}

func (p *Player) LookAround(...string) {
	out := p.addMyThink() + p.playersInRoom()
	p.output <- out
	return
}

func (p *Player) Move(s ...string) {
	var out string
	n := p.Position.Neighbors
	for _, v := range n {
		if p.Position.Location != v.Location {
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
	v := MapOfItems[s[0]]
	switch {
	case v.iType != "clothes":
		out = "не могу одеть: " + s[0]
	case ok:
		out = s[0] + " - уже одето"
	case p.Position.Checker(s[0]):
		p.addItemFromRoom(s[0])
		p.Position.refreshRoomDes()
		out = "вы одели: " + s[0]
	}
	p.output <- out
	return
}

func (p *Player) Take(s ...string) {
	var out string
	if _, ok := p.Checker(MapOfItems["рюкзак"].Name); ok {
		if p.Position.Checker(s[0]) {
			p.addItemFromRoom(s[0])
			p.Position.refreshRoomDes()
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
		} else {
			out = fmt.Sprintf("Не могу применить %s к %s\n", s[0], s[1])
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
	for _, pir := range p.Position.WhoInMe {
		pir.output <- out
	}
	return
}

func (p *Player) SayPlayer(ss ...string) {
	var out string
	if p2, ok := p.Position.WhoInMe[ss[0]]; ok {
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
