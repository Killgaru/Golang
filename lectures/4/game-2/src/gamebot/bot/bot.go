package main

import (
	// "encoding/json"
	"gamebot/game"
	// "io/ioutil"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/telegram-bot-api.v4"
)

// для вендоринга используется GB
// сборка проекта осуществляется с помощью gb build
// установка зависимостей - gb vendor fetch gopkg.in/telegram-bot-api.v4
// установка зависимостей из манифеста - gb vendor restore

const (
	WebhookURL    = "https://tg-room-game-bot.herokuapp.com/"
	DefTimeOfLife = 15 * time.Minute
)

var (
	allActions = [][]tgbotapi.KeyboardButton{
		[]tgbotapi.KeyboardButton{
			tgbotapi.KeyboardButton{Text: "осмотреться"},
			tgbotapi.KeyboardButton{Text: "идти"},
			tgbotapi.KeyboardButton{Text: "одеть"},
		},
		[]tgbotapi.KeyboardButton{
			tgbotapi.KeyboardButton{Text: "взять"},
			tgbotapi.KeyboardButton{Text: "применить"},
			tgbotapi.KeyboardButton{Text: "сказать"},
		},
		[]tgbotapi.KeyboardButton{
			tgbotapi.KeyboardButton{Text: "сказать_игроку"},
			// tgbotapi.KeyboardButton{Text: "Exit"},
		},
	}
	Admin = struct {
		ID     int64
		inGame bool
	}{
		ID:     705987198,
		inGame: false,
	}
	NotPlayers = map[int64]string{}
)

func keyboardFromSS(ss []string) (out [][]tgbotapi.KeyboardButton) {
	for i, j := 0, 0; i+j < len(ss); i++ {
		out = append(out, []tgbotapi.KeyboardButton{})
		for j = 0; j < 3 && i+j < len(ss); j++ {
			out[i] = append(out[i], tgbotapi.KeyboardButton{Text: ss[i+j]})
		}
	}
	return
}

func keyboardFromRooms(r *game.Room) (out [][]tgbotapi.KeyboardButton) {
	sd := r.ShortDescrip
	substr := "можно пройти - "
	ind := strings.Index(sd, substr)
	ss := strings.Split(sd[ind+len(substr):], ", ")
	out = keyboardFromSS(ss)
	return
}

func keyboardFromRoomItem(r *game.Room) (out [][]tgbotapi.KeyboardButton) {
	var ss []string
	for _, v := range r.Items {
		ss = append(ss, v.Name)
	}
	out = keyboardFromSS(ss)
	return
}

func keyboardFromInventory(p *game.Player) (out [][]tgbotapi.KeyboardButton) {
	var ss []string
	for _, v := range p.Inventory {
		ss = append(ss, v.Name)
	}
	out = keyboardFromSS(ss)
	return
}

func keyboardFromRoomPlayers(r *game.Room) (out [][]tgbotapi.KeyboardButton) {
	var ss []string
	for _, v := range r.WhoInMe {
		ss = append(ss, v.Name)
	}
	out = keyboardFromSS(ss)
	return
}

func keyboardFromRoomDoorsItem(r *game.Room) (out [][]tgbotapi.KeyboardButton) {
	var ss []string
	for _, v := range r.Doors {
		ss = append(ss, v.Name)
	}
	for _, v := range r.Items {
		ss = append(ss, v.Name)
	}
	out = keyboardFromSS(ss)
	return
}

func caseForKeyboard(p *game.Player, parts []string) (keyboard interface{}, msg, msg2 string) {
	//keyboard изменить на интерфейс и тут сразу строить клавы, а там делать свич по типам
	switch parts[0] {
	case "идти":
		keyboard = tgbotapi.NewReplyKeyboard(keyboardFromRooms(p.Position)...)
		msg = "Куда бы пойтиии...?"
	case "одеть":
		if len(p.Position.Items) == 0 {
			msg2 = "В этой комнате нет ничего что можно одеть."
			break
		}
		keyboard = tgbotapi.NewReplyKeyboard(keyboardFromRoomItem(p.Position)...)
		msg = "Что бы одеть?"
	case "взять":
		if len(p.Position.Items) == 0 {
			msg2 = "В этой комнате нет ничего что можно взять."
			break
		}
		keyboard = tgbotapi.NewReplyKeyboard(keyboardFromRoomItem(p.Position)...)
		msg = "Что бы взять?"
	case "применить":
		if len(p.Inventory) == 0 {
			msg2 = "Инвентарь пуст. Нечего применять."
			break
		}
		if len(p.Position.Items) == 0 && len(p.Position.Doors) == 0 {
			msg2 = "Здесь не к чему применить что-либо."
			break
		}
		switch len(parts) {
		case 1:
			keyboard = tgbotapi.NewReplyKeyboard(keyboardFromInventory(p)...)
			msg = "Что бы выбрать?"
			break
		case 2:
			keyboard = tgbotapi.NewReplyKeyboard(keyboardFromRoomDoorsItem(p.Position)...)
			msg = "К чему применить?"
		}
	case "сказать":
		keyboard = tgbotapi.NewRemoveKeyboard(true)
		msg = "Что сказать?"
	case "сказать_игроку":
		if len(p.Position.WhoInMe) <= 1 {
			msg2 = "В этой комнате нет других игроков."
		}
		switch len(parts) {
		case 1:
			keyboard = tgbotapi.NewReplyKeyboard(keyboardFromRoomPlayers(p.Position)...)
			msg = "Кому сказать?"
			break
		case 2:
			keyboard = tgbotapi.NewRemoveKeyboard(true)
			msg = "Что сказать?"
		}
	}
	return
}

func main() {
	quit := make(chan bool)
	// Heroku прокидывает порт для приложения в переменную окружения PORT
	port := os.Getenv("PORT")
	bot, err := tgbotapi.NewBotAPI("761200557:AAHfMkIpOWaLZJBqauX-rw69oGm5QmuVlps")
	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	// Устанавливаем вебхук
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		log.Fatal(err)
	}

	game.InitGame()

	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":"+port, nil)

	// получаем все обновления из канала updates
	for update := range updates {
		var message tgbotapi.MessageConfig
		var wg sync.WaitGroup
		txt := update.Message.Text
		log.Println("received text: ", txt)
		log.Printf("from user:\n ID:%v\n Username:%v\n",
			update.Message.From.ID, update.Message.From.UserName)
		log.Println("Players in game: ", game.PlayersInGame)

		p, ok := game.FindThisUserInGame(update.Message.Chat.ID)
		if ok {
			p.TimeLife.Reset(DefTimeOfLife)
		}
		if _, ok2 := NotPlayers[update.Message.Chat.ID]; !ok && !ok2 {
			NotPlayers[update.Message.Chat.ID] = ""
		}
		switch txt {
		case "/start":
			if !ok {
				message = tgbotapi.NewMessage(update.Message.Chat.ID,
					"Чтобы начать введите имя персонажа.")
				NotPlayers[update.Message.Chat.ID] = "/start"
			} else {
				message = tgbotapi.NewMessage(update.Message.Chat.ID,
					"Вы уже в игре.")
			}
			bot.Send(message)
			continue
		case "/restart":
			if Admin.ID == update.Message.Chat.ID {
				var PlayersID = []int64{}
				if !Admin.inGame {
					message = tgbotapi.NewMessage(Admin.ID, "Игра перезагружаеться...")
					message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(message)
				}
				for _, p := range game.PlayersInGame {
					PlayersID = append(PlayersID, p.ID)
					quit <- true
					message = tgbotapi.NewMessage(p.ID, "Игра перезагружаеться...")
					message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(message)
				}

				game.InitGame()
				if !Admin.inGame {
					message = tgbotapi.NewMessage(Admin.ID,
						fmt.Sprintf("Игра перезагружена.\n"+
							"/start - создание персонажа"))
					message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					bot.Send(message)
				}
				for _, v := range PlayersID {
					message = tgbotapi.NewMessage(v,
						fmt.Sprintf("Игра перезагружена.\n"+
							"Ваш персонаж был удалён.\n"+
							"/start - создание персонажа"))
					message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
					if Admin.ID == v {
						Admin.inGame = false
					}
					bot.Send(message)
				}
				continue
			}
			break
		case "/help":
			if Admin.ID == update.Message.Chat.ID {
				message = tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf("/start - создание персонажа\n"+
						"/restart - перезапуск игры\n"+
						"/help - список доступных комманд"))
			} else {
				message = tgbotapi.NewMessage(update.Message.Chat.ID,
					fmt.Sprintf("/start - создание персонажа\n"+
						"/help - список доступных комманд"))
			}
			bot.Send(message)
			continue
		default:
			if !ok {
				if np := NotPlayers[update.Message.Chat.ID]; np == "/start" {
					game.AddPlayer(game.NewPlayer(txt), update.Message.Chat.ID)
					game.PlayersInGame[txt].TimeLife = time.NewTimer(DefTimeOfLife)
					go func(chan bool) {
						select {
						case <-game.PlayersInGame[txt].TimeLife.C:
							for _, it := range game.PlayersInGame[txt].Inventory {
								it.SpawnItem()
							}
							game.PlayersInGame[txt].DelPlayerRoom()
							delete(game.PlayersInGame, txt)
							if Admin.ID == update.Message.Chat.ID {
								Admin.inGame = false
							}
							message = tgbotapi.NewMessage(update.Message.Chat.ID,
								fmt.Sprintf("Ваш персонаж был удалён из-за длительного отсутствия в игре\n"+
									"/start - создание персонажа"))
							message.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
							bot.Send(message)
						case <-quit:
							game.PlayersInGame[txt].TimeLife.Stop()
						}
					}(quit)
					delete(NotPlayers, update.Message.Chat.ID)
					message = tgbotapi.NewMessage(update.Message.Chat.ID,
						txt+" для начала осмотрись")
					message.ReplyMarkup = tgbotapi.NewReplyKeyboard(allActions...)
					if Admin.ID == update.Message.Chat.ID {
						Admin.inGame = true
					}
				} else {
					message = tgbotapi.NewMessage(update.Message.Chat.ID,
						"Чтобы создать персонажа, нажми /start")
				}
				bot.Send(message)
				continue
			}
		}
		if _, ok := game.Commands[txt]; ok {
			p.ComFragment = txt
		}
		if p.ComFragment != "" {
			if p.ComFragment != txt {
				p.ComFragment += " " + txt
			}
			parts := strings.Split(p.ComFragment, " ")
			c := game.Commands[parts[0]]
			if len(parts)-1 < c.MinNumArg {
				// тут меняются кнопочки
				keyboard, msg, msg2 := caseForKeyboard(p, parts)
				if msg2 != "" {
					message = tgbotapi.NewMessage(update.Message.Chat.ID, msg2)
					message.ReplyMarkup = tgbotapi.NewReplyKeyboard(allActions...)
					p.ComFragment = ""
					bot.Send(message)
					continue
				}
				message = tgbotapi.NewMessage(update.Message.Chat.ID, msg)
				switch k := keyboard.(type) {
				case tgbotapi.ReplyKeyboardMarkup:
					message.ReplyMarkup = k
				case tgbotapi.ReplyKeyboardRemove:
					message.ReplyMarkup = k
				}
				bot.Send(message)
				continue
			}
			wg.Add(1)
			go func() {
				message = tgbotapi.NewMessage(update.Message.Chat.ID, p.HandleOutput())
				message.ReplyMarkup = tgbotapi.NewReplyKeyboard(allActions...)
				wg.Done()
			}()
			p.HandleInput(p.ComFragment)
			wg.Wait()
			p.ComFragment = ""
			bot.Send(message)
			continue
		}
		wg.Add(1)
		go func() {
			message = tgbotapi.NewMessage(update.Message.Chat.ID, p.HandleOutput())
			message.ReplyMarkup = tgbotapi.NewReplyKeyboard(allActions...)
			wg.Done()
		}()
		p.HandleInput(txt)
		wg.Wait()
		bot.Send(message)
	}
}
