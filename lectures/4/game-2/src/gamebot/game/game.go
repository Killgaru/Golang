// main
package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	// "sync"
	// "time"
)

type command struct {
	Name      string
	MinNumArg int
	function  func(*Player, ...string)
}

var (
	Commands = map[string]command{
		"осмотреться": command{
			Name:      "осмотреться",
			MinNumArg: 0,
			function:  (*Player).LookAround},
		"идти": command{
			Name:      "идти",
			MinNumArg: 1,
			function:  (*Player).Move},
		"одеть": command{
			Name:      "одеть",
			MinNumArg: 1,
			function:  (*Player).Robe},
		"взять": command{
			Name:      "взять",
			MinNumArg: 1,
			function:  (*Player).Take},
		"применить": command{
			Name:      "применить",
			MinNumArg: 2,
			function:  (*Player).Apply},
		"сказать": command{
			Name:      "сказать",
			MinNumArg: 1,
			function:  (*Player).SayRoom},
		"сказать_игроку": command{
			Name:      "сказать_игроку",
			MinNumArg: 2,
			function:  (*Player).SayPlayer},
		"Exit": command{
			Name:      "Exit",
			MinNumArg: 0,
			function:  (*Player).ExitGame},
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

func FindThisUserInGame(id int64) (p *Player, ok bool) {
	for _, v := range PlayersInGame {
		if v.ID == id {
			ok = true
			p = v
			break
		}
	}
	return
}

func handleCommand(p *Player, s1 string) {
	s := strings.Split(s1, " ")
	if v, ok := Commands[s[0]]; ok {
		v.function(p, s[1:]...)
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

func InitGame() {
	InitRooms()
	InitItems()
	InitDoors()
	InitPlayers()
}

// func Game0() {
// 	var (
// 		wg              sync.WaitGroup
// 		playername, out string
// 	)
// 	InitGame()
// 	fmt.Println("Enter your name")
// 	playername = Reader()
// 	AddPlayer(NewPlayer(playername))
// 	fmt.Println(playername, ", welcome in Game1")
// 	for {
// 		wg.Add(1)
// 		go func() {
// 			out = PlayersInGame[playername].HandleOutput()
// 			wg.Done()
// 		}()
// 		PlayersInGame[playername].HandleInput(Reader())
// 		wg.Wait()
// 		if out != "Exit" {
// 			fmt.Println(out)
// 		} else {
// 			break
// 		}
// 	}
// 	return
// }

// func main() {
// 	// initGame()
// 	Game0()
// }
