// main
package main

import (
	"fmt"
	"strings"
	// "sync"
)

var (
	notInitiated bool = true
	commands          = map[string]func(*Players, ...string) string{
		"осмотреться": (*Players).LookAround,
		"идти":        (*Players).Move,
		"одеть":       (*Players).Robe,
		"взять":       (*Players).Take,
		"применить":   (*Players).Apply,
		"Exit":        (*Players).ExitGame,
	}
)

type Checker interface {
	Checker(string) bool
}

func DelElemSliceString(s []string, i int) (out []string) {
	out = append(append(s[:0], s[:i]...), s[i+1:]...)
	return
}

func DelElemSliceItem(it []*Item, i int) (out []*Item) {
	out = append(append(it[:0], it[:i]...), it[i+1:]...)
	return
}

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
	InitRooms()
	InitItems()
	InitDoors()
	InitPlayers()
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
	initGame()
	// fmt.Println(handleCommand("осмотреться"))
	// fmt.Println(Player.position.description)
	// fmt.Println(handleCommand("идти коридор"))
	// Game0()
}
