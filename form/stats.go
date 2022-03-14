package form

import (
	"github.com/df-mc/dragonfly/server/player/form"
	"strconv"
	"velvet/db"
)

func Stats(target string, d db.PlayerData) form.Menu {
	var s string

	keys := []string{"Player", "Rank", "Kills", "Deaths"}
	values := []string{target + "\n", d.Rank, strconv.Itoa(int(d.Kills)), strconv.Itoa(int(d.Deaths))}

	l := len(keys)
	for i := 0; i < l; i++ {
		k := keys[i]
		v := values[i]
		s += "§6" + k + ": §b" + v + "\n"
	}

	return form.NewMenu(NopSubmit, "§6"+target+"'s Stats").WithButtons(form.NewButton("Exit", "")).WithBody(s)
}
