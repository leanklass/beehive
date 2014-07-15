/*
 *    Copyright (C) 2014 Christian Muehlhaeuser
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
 *      Christian Muehlhaeuser <muesli@gmail.com>
 */

// beehive's web-module.
package webbee

import (
	"encoding/json"
	"github.com/hoisie/web"
	"github.com/muesli/beehive/modules"
	"io/ioutil"
	"log"
)

type WebBee struct {
	modules.Module

	addr        string

	eventChan chan modules.Event
}

func (mod *WebBee) Run(cin chan modules.Event) {
	mod.eventChan = cin
	web.Run(mod.addr)
}

func (mod *WebBee) Action(action modules.Action) []modules.Placeholder {
	outs := []modules.Placeholder{}
	return outs
}

func (mod *WebBee) GetRequest(ctx *web.Context) {
	ev := modules.Event{
		Bee:  mod.Name(),
		Name: "get",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}

	for k, v := range ctx.Params {
		log.Println("GET param:", k, "=", v)

		ph := modules.Placeholder{
			Name: k,
			Type: "string",
			Value: v,
		}
		ev.Options = append(ev.Options, ph)
	}

	mod.eventChan <- ev
}

func (mod *WebBee) PostRequest(ctx *web.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(b, &payload)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	ev := modules.Event{
		Bee:  mod.Name(),
		Name: "post",
		Options: []modules.Placeholder{
			modules.Placeholder{
				Name:  "json",
				Type:  "map",
				Value: payload,
			},
			modules.Placeholder{
				Name:  "ip",
				Type:  "string",
				Value: "tbd",
			},
		},
	}
	mod.eventChan <- ev
}
