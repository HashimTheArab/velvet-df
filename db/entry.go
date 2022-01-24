package db

import (
	"time"
)

type BanEntry struct {
	XUID   string `db:"XUID"`
	IGN    string `db:"IGN"`
	Mod    string `db:"Mod"`
	Reason string `db:"Reason"`
	// Expires is a unix timestamp of when the ban will expire. If it is -1, the ban is a blacklist.
	Expires int64 `db:"Expires"`
}

// Blacklist returns whether the ban is a blacklist (permanent)
func (ban BanEntry) Blacklist() bool {
	return ban.Expires == -1
}

// Update saves the passed xuid for the ban.
func (ban BanEntry) Update(xuid string) {
	_, _ = db.Exec("UPDATE Bans set XUID=? WHERE IGN=?", xuid, ban.IGN)
}

// FormattedExpiration returns how long is left on the ban, formatted nicely.
func (ban BanEntry) FormattedExpiration() string {
	if ban.Expires == -1 {
		return "Never"
	}
	return (time.Second * time.Duration(ban.Expires-time.Now().Unix())).Round(time.Second).String()
}
