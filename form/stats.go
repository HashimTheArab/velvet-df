package form

import (
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strconv"
	"velvet/db"
	"velvet/session"
)

// StatsOffline is the stats form that shows an offline players stats.
func StatsOffline(d *db.Entry) form.Menu {
	return statsForm(d.DisplayName, d.Rank, d.Kills, d.Deaths)
}

// StatsOnline is the stats form that shows an online players stats.
func StatsOnline(s *session.Session) form.Menu {
	return statsForm(s.Player.Name(), s.RankName(), s.Kills(), s.Deaths())
}

// statsForm creates a stats form that can be used for both online and offline players.
func statsForm(target string, rank string, kills, deaths uint32) form.Menu {
	var s string

	keys := []string{"Player", "Rank", "Kills", "Deaths"}
	values := []string{target + "\n", rank, strconv.Itoa(int(kills)), strconv.Itoa(int(deaths))}

	l := len(keys)
	for i := 0; i < l; i++ {
		k := keys[i]
		v := values[i]
		s += "§6" + k + ": §b" + v + "\n"
	}

	return form.NewMenu(NopSubmit, text.Colourf("<gold>%s's Stats</gold>", target)).WithButtons(form.NewButton("Exit", "")).WithBody(s)
}
