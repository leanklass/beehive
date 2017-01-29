/*
 *    Copyright (C) 2016 Gonzalo Izquierdo
 *                  2017 Christian Muehlhaeuser
 *
 *    This program is free software: you can redistribute it and/or modify
 *    it under the terms of the GNU Affero General Public License as published
 *    by the Free Software Foundation, either version 3 of the License, or
 *    (at your option) any later version.
 *
 *    This program is distributed in the hope that it will be useful,
 *    but WITHOUT ANY WARRANTY; without even the implied warranty of
 *    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *    GNU Affero General Public License for more details.
 *
 *    You should have received a copy of the GNU Affero General Public License
 *    along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *    Authors:
 *      Gonzalo Izquierdo <lalotone@gmail.com>
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// Package telegrambee is a Bee that can connect to Telegram.
package telegrambee

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/muesli/beehive/bees"
	telegram "gopkg.in/telegram-bot-api.v4"
)

type TelegramBee struct {
	bees.Bee

	// Telegram bot API Key
	apiKey string
	// Bot API client
	bot *telegram.BotAPI
}

func (mod *TelegramBee) Action(action bees.Action) []bees.Placeholder {
	outs := []bees.Placeholder{}

	switch action.Name {
	case "send":
		chatId := ""
		text := ""
		action.Options.Bind("chatId", &chatId)
		action.Options.Bind("text", &text)

		cid, err := strconv.Atoi(chatId)
		if err != nil {
			panic("Invalid telegram chat ID")
		}

		msg := telegram.NewMessage(int64(cid), text)
		_, err = mod.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message %v", err)
		}
	}

	return outs
}

func (mod *TelegramBee) Run(eventChan chan bees.Event) {
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := mod.bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		ev := bees.Event{
			Bee:  mod.Name(),
			Name: "message",
			Options: []bees.Placeholder{
				{
					Name:  "text",
					Type:  "string",
					Value: update.Message.Text,
				},
				{
					Name:  "chatID",
					Type:  "string",
					Value: strconv.FormatInt(update.Message.Chat.ID, 10),
				},
				{
					Name:  "userID",
					Type:  "string",
					Value: strconv.Itoa(update.Message.From.ID),
				},
			},
		}
		eventChan <- ev
	}
}

func (mod *TelegramBee) Stop() {
	log.Println("Stopping the Telegram bee")
}

func (mod *TelegramBee) ReloadOptions(options bees.BeeOptions) {
	mod.SetOptions(options)

	apiKey := getApiKey(&options)
	bot, err := telegram.NewBotAPI(apiKey)
	if err != nil {
		panic("Authorization failed, make sure the Telegram API key is correct")
	}
	log.Printf("TELEGRAM: Authorized on account %s", bot.Self.UserName)

	mod.apiKey = apiKey
	mod.bot = bot
}

// Gets the Bot's API key from a file, the recipe config or the
// TELEGRAM_API_KEY environment variable.
func getApiKey(options *bees.BeeOptions) string {
	var apiKey string
	options.Bind("apiKey", &apiKey)

	if strings.HasPrefix(apiKey, "file://") {
		buf, err := ioutil.ReadFile(strings.TrimPrefix(apiKey, "file://"))
		if err != nil {
			panic("Error reading API key file " + apiKey)
		}
		apiKey = string(buf)
	}

	if strings.HasPrefix(apiKey, "env://") {
		buf := strings.TrimPrefix(apiKey, "env://")
		apiKey = os.Getenv(string(buf))
	}

	return strings.TrimSpace(apiKey)
}
