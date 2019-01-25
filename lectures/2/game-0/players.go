// players
package main

import (
	// "fmt"
	"strings"
)

var (
	Player        Players
	PlayersInGame map[string]*Players
)

type thinkP struct {
	keyFunc func(*Players, interface{}) bool
	argFunc interface{}
	answers map[bool]string
}

type Players struct {
	Name      string
	inventory []*Item
	thoughts  map[*Room]thinkP
	position  *Room
}

func (p *Players) inInventory(it interface{}) bool {
	var out, ok bool = true, false
	val := it.([]*Item)
	for _, v := range val {
		_, ok = p.Checker(v.Name)
		out = out && ok
	}
	return out
}

func (t *thinkP) checkThink(p *Players) string {
	return t.answers[t.keyFunc(p, t.argFunc)]
}

func (p *Players) checkMyThink(r *Room) string {
	if v, ok := p.thoughts[r]; ok {
		return v.checkThink(p)
	}
	return ""
}

func (p *Players) createPlayer(Name string, inventory []*Item,
	thoughts map[*Room]thinkP, position *Room) {
	p.Name = Name
	p.inventory = inventory
	p.thoughts = thoughts
	p.position = position
	PlayersInGame[Name] = p
	return
}

func (p *Players) LookAround(...string) (out string) {
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

func (p *Players) Move(s ...string) string {
	n := p.position.Neighbors
	for _, v := range n {
		if s[0] == v.Name {
			if p.position.checkDoor(v) {
				return "дверь закрыта"
			}
			p.DelPlayerRoom()
			p.position = v
			p.AddPlayerRoom()
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
			p.position.WhoInMe = DelElemSliceString(p.position.WhoInMe, i)
		}
	}
	return
}

func (p *Players) addItemFromRoom(s string) {
	for i, v := range p.position.items {
		if v.Name == s {
			p.position.items = DelElemSliceItem(p.position.items, i)
			p.inventory = append(p.inventory, v)
			return
		}
	}
	return
}

func (p *Players) Robe(s ...string) string {
	if _, ok := p.Checker(s[0]); ok {
		return s[0] + " - уже одето"
	}
	if p.position.Checker(s[0]) {
		p.addItemFromRoom(s[0])
		p.position.refreshRoomDes()
		return "вы одели: " + s[0]
	}
	return "не могу одеть: " + s[0]
}

func (p *Players) Take(s ...string) string {
	if _, ok := p.Checker(MapOfItems["рюкзак"].Name); ok {
		if p.position.Checker(s[0]) {
			p.addItemFromRoom(s[0])
			p.position.refreshRoomDes()
			return "предмет добавлен в инвентарь: " + s[0]
		}
		return "нет такого"
	}
	return "некуда класть"
}

func (p *Players) Apply(s ...string) string {
	if it, ok := p.Checker(s[0]); ok {
		if it.iType == "key" {
			return p.applyKeyDoor(it, s[1])
		}
	}
	return "нет предмета в инвентаре - " + s[0]
}

func (p *Players) applyKeyDoor(it *Item, s string) string {
	for _, d := range p.position.doors {
		if s == d.Name && it == d.Key {
			d.ChengeStatus()
			return d.CheckStatus()
		}
	}
	return "не к чему применить"
}

func (p *Players) ExitGame(...string) string {
	return "Exit"
}

func (p *Players) Checker(s string) (*Item, bool) {
	for _, v := range p.inventory {
		if v.Name == s {
			return v, true
		}
	}
	return &Item{}, false
}

func InitPlayers() {
	PlayersInGame = map[string]*Players{}
	Player.createPlayer(
		"Player1",
		[]*Item{},
		map[*Room]thinkP{
			World["домой"]["домой"]["кухня"]: thinkP{
				keyFunc: (*Players).inInventory,
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
