// rooms
package main

import (
	"fmt"
	"strings"
	// "sync"
)

var (
	World map[string]Locations
)

type funcRoomDes func(*Room, ...string)

type Locations map[string]map[string]*Room

type Room struct {
	Name         string
	Location     string
	shortDescrip string
	description  string
	defaultDes   string
	items        []*Item
	doors        []*Door
	Neighbors    []*Room
	tempNei      []string
	WhoInMe      map[string]*Player /*[]string*/
	funcsRoomDes []funcRoomDes
}

func createRoom(Name, Loc, shortDes, des, defaultDes string,
	nei []string, funcsRoomDes []funcRoomDes) {
	cr := Room{
		Name:         Name,
		Location:     Loc,
		shortDescrip: shortDes,
		description:  des,
		defaultDes:   defaultDes,
		tempNei:      nei,
		funcsRoomDes: funcsRoomDes,
	}
	cr.addToWorld()
	cr.setRoomDes()
	cr.addNei()
	cr.WhoInMe = make(map[string]*Player)
	return
}

func (r *Room) tempNeiSplit(i int) (s1, s2 string, err error) {
	ss := strings.SplitN(r.tempNei[i], "; ", 2)
	if len(ss) != 2 {
		err = fmt.Errorf("Некорректно задан сосед: %v, для комнаты: %s",
			r.tempNei[i], r.Name)
		return
	}
	s1 = ss[0]
	s2 = ss[1]
	return
}

func (r *Room) addNei() {
	var (
		loc, room string
		err       error
	)
	for i := 0; i < len(r.tempNei); i++ {
		loc, room, err = r.tempNeiSplit(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		if l1, ok := World[loc]; ok {
			if r1, ok := l1[loc][room]; ok {
				r.Neighbors = append(r.Neighbors, r1)
				r.tempNei = DelElemSliceString(r.tempNei, i)
				r1.addNei()
			}
		}

	}
	return
}

func (r *Room) addToWorld() {
	if _, ok := World[r.Location]; !ok {
		World[r.Location] = Locations{}
	}
	if _, ok := World[r.Location][r.Location]; !ok {
		World[r.Location][r.Location] = map[string]*Room{}
	}
	if _, ok := World[r.Location][r.Location][r.Name]; !ok {
		World[r.Location][r.Location][r.Name] = r
	}
	return
}

func (r *Room) checkDoor(r2 *Room) bool {
	if len(r.doors) == 0 || len(r2.doors) == 0 {
		return false
	}
	for _, v1 := range r.doors {
		for _, v2 := range r2.doors {
			if v1 == v2 {
				return v1.Status
			}
		}

	}
	return true
}

func (r *Room) desClean() {
	r.description = ""
	return
}

func (r *Room) desWalk(des string) (out string) {
	var (
		s, room, loc, nei string
		err               error
	)
	s = "можно пройти - "
	if des != "" {
		s = " " + s
	}
	for i := range r.tempNei {
		loc, room, err = r.tempNeiSplit(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		if r.Location == loc {
			nei = room
		} else {
			nei = loc
		}
		if i != len(r.tempNei)-1 {
			s += fmt.Sprintf("%s, ", nei)
		} else {
			s += fmt.Sprintf("%s", nei)
		}

	}
	out = des + s
	return
}

func (r *Room) desFromItems(...string) {
	if len(r.items) != 0 {
		var desItem, out string
		m := make(map[string]string)
		d := []string{}
		for _, v := range r.items {
			desItem = v.description[r]
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

func (r *Room) setRoomDes() {
	for _, v := range r.funcsRoomDes {
		v(r)
	}
	r.description = r.desWalk(r.description)
	r.shortDescrip = r.desWalk(r.shortDescrip)
	return
}

func (r *Room) refreshRoomDes(ss ...string) {
	var out string
	for _, v := range ss {
		out += v
	}
	s := strings.Split(r.description, ". ")
	r.desClean()
	for _, v := range r.funcsRoomDes {
		v(r, out)
	}
	r.description += fmt.Sprint(" ", s[len(s)-1])
	return
}

func (r *Room) setDefDesAsDes(...string) {
	r.description = r.defaultDes
	return
}

func (r *Room) Checker(s string) bool {
	for _, v := range r.items {
		if v.Name == s {
			return true
		}
	}
	return false
}

func InitRooms() {
	World = map[string]Locations{}
	createRoom(
		"кухня", "домой",
		"кухня, ничего интересного.",
		"",
		"ты находишься на кухне, на столе чай.",
		[]string{"домой; коридор"},
		[]funcRoomDes{(*Room).setDefDesAsDes},
	)
	createRoom(
		"коридор", "домой",
		"ничего интересного.",
		"",
		"ничего интересного.",
		[]string{"домой; кухня", "домой; комната", "улица; улица"},
		[]funcRoomDes{(*Room).setDefDesAsDes},
	)
	createRoom(
		"комната", "домой",
		"ты в своей комнате.",
		"",
		"на столе: ключи, конспекты, на стуле - рюкзак.",
		[]string{"домой; коридор"},
		[]funcRoomDes{(*Room).desFromItems},
	)
	createRoom(
		"улица", "улица",
		"на улице весна.",
		"",
		"на улице весна.",
		[]string{"домой; коридор"},
		[]funcRoomDes{(*Room).setDefDesAsDes},
	)
	// fmt.Printf("кухня, соседи: %v\nкухня временные соседи: %v\n",
	// 	World["домой"]["домой"]["кухня"].Neighbors,
	// 	World["домой"]["домой"]["кухня"].tempNei)
	return
}
