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
	"sort"
	"strings"
)

func do_help(entity Entity, args ...string) {
	player := entity.(*PlayerProfile)
	db := DB()
	if len(args) > 0 {
		help := db.GetHelp(strings.Join(args, " "))
		if len(help) > 0 {
			player.Send("\r\n&W%s&d\r\n", MakeTitle(help[0].Name, ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_CENTER))
			player.Send("&YKeywords: &g%v&d\r\n", help[0].Keywords)
			player.Send("&w%s&d\r\n", help[0].Desc)
		}
	} else {
		player.Send("\r\n&W%s&d\r\n", MakeTitle("Help", ANSI_TITLE_STYLE_NORMAL, ANSI_TITLE_ALIGNMENT_CENTER))
		keys := []string{}
		for i := 0; i < len(db.helps); i++ {
			keys = append(keys, db.helps[i].Keywords...)
		}
		sort.Strings(keys)
		buf := ""
		for i := 0; i < len(keys); i++ {
			buf += "&W" + keys[i] + "&d "
			if i%4 == 0 && i > 1 {
				buf += "\r\n"
			}
		}
		player.Send(buf)
	}

}