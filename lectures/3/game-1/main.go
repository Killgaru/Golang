// main
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	// "time"
)

var (
	commands = map[string]func(*Player, ...string){
		"осмотреться":    (*Player).LookAround,
		"идти":           (*Player).Move,
		"одеть":          (*Player).Robe,
		"взять":          (*Player).Take,
		"применить":      (*Player).Apply,
		"сказать":        (*Player).SayRoom,
		"сказать_игроку": (*Player).SayPlayer,
		"Exit":           (*Player).ExitGame,
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

func handleCommand(p *Player, s1 string) {
	s := strings.Split(s1, " ")
	if v, ok := commands[s[0]]; ok {
		v(p, s[1:]...)
	} else {
		p.output <- fmt.Sprint("неизвестная команда")
	}
	return
}

func Reader() (s string) {
	scaner := bufio.NewScanner(os.Stdin)
	scaner.Scan()
	s = scaner.Text()
	return
}

func initGame() {
	InitRooms()
	InitItems()
	InitDoors()
	InitPlayers()
}

func Game0() {
	var (
		wg              sync.WaitGroup
		playername, out string
	)
	initGame()
	fmt.Println("Enter your name")
	playername = Reader()
	addPlayer(NewPlayer(playername))
	fmt.Println(playername, ", welcome in Game1")
	for {
		wg.Add(1)
		go func() {
			out = PlayersInGame[playername].HandleOutput()
			wg.Done()
		}()
		PlayersInGame[playername].HandleInput(Reader())
		wg.Wait()
		if out != "Exit" {
			fmt.Println(out)
		} else {
			break
		}
	}
	return
}

func main() {
	// initGame()
	Game0()
}
