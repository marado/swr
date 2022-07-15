/*  Space Wars Rebellion Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v3"
)

var CommandFuncs = map[string]func(Entity, ...interface{}){
	"do_say":  do_say,
	"do_look": do_look,
	"do_who":  do_who,
}
var GMCommandFuncs = map[string]func(Entity, ...interface{}){
	"do_area_create": do_area_create,
	"do_area_set":    do_area_set,
	"do_area_remove": do_area_remove,
	"do_area_reset":  do_area_reset,
	"do_room_create": do_room_create,
	"do_room_edit":   do_room_edit,
	"do_room_find":   do_room_find,
	"do_room_remove": do_room_remove,
	"do_room_reset":  do_room_reset,
	"do_room_set":    do_room_set,
	"do_mob_create":  do_mob_create,
	"do_mob_stat":    do_mob_stat,
	"do_mob_find":    do_mob_find,
	"do_mob_remove":  do_mob_remove,
	"do_mob_reset":   do_mob_reset,
	"do_mob_set":     do_mob_set,
	"do_item_create": do_item_create,
	"do_item_stat":   do_item_stat,
	"do_item_find":   do_item_find,
	"do_item_remove": do_item_remove,
	"do_item_set":    do_item_set,
}

var Commands []*Command = make([]*Command, 0)

type Command struct {
	Name     string   `yaml:"name"`
	Keywords []string `yaml:"keywords,flow"`
	Level    uint     `yaml:"level"`
	Func     string   `yaml:"func"`
}

func CommandsLoad() {
	log.Printf("Loading commands list.")
	fp, err := ioutil.ReadFile("data/sys/commands.yml")
	ErrorCheck(err)
	err = yaml.Unmarshal(fp, &Commands)
	ErrorCheck(err)
	log.Printf("Commands successfully loaded.")
}
func command_map_to_func(name string) func(Entity, ...interface{}) {
	if k, ok := CommandFuncs[name]; ok {
		return k
	}
	if k, ok := GMCommandFuncs[name]; ok {
		return k
	}
	return do_nothing
}
func command_fuzzy_match(command string) []Command {
	ret := []Command{}
	for _, com := range Commands {
		for _, keyword := range com.Keywords {
			if strings.HasPrefix(strings.ToLower(command), strings.ToLower(keyword)) {
				ret = append(ret, *com)
			}
		}
	}
	return ret
}
func do_command(entity Entity, input string) {
	commands := command_fuzzy_match(input)
	args := strings.Split(input, " ")
	if len(commands) > 0 {
		args[0] = strings.TrimPrefix(args[0], "'")
		args[0] = strings.TrimPrefix(args[0], "\"")
		args[0] = strings.TrimPrefix(args[0], ".")
		args[0] = strings.TrimPrefix(args[0], ":")
		args[0] = strings.TrimPrefix(args[0], "!")
		command_map_to_func(commands[0].Func)(entity, args)
	} else {
		if entity.IsPlayer() {
			entity.Send("\r\nHuh?\r\n")
		}
	}
}