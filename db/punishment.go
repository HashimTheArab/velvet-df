package db

import (
	"time"
)

// Punishment is a struct used for storing punishments such as mutes and bans.
type Punishment struct {
	Mod        string    `bson:"mod"`
	Reason     string    `bson:"reason"`
	Permanent  bool      `bson:"permanent"`
	Expiration time.Time `bson:"expiration"`
}

// Remaining returns the remaining duration of the punishment.
func (p Punishment) Remaining() time.Duration {
	return time.Until(p.Expiration).Round(time.Second)
}

// Expired checks if the punishment has expired.
func (p Punishment) Expired() bool {
	return !p.Permanent && time.Now().After(p.Expiration)
}

// Update saves the passed xuid for the ban.
func (p Punishment) Update(ign, xuid string) {
	_ = findBan(ign).Update(map[string]string{"xuid": xuid})
}

// FormattedExpiration returns how long is left on the ban, formatted nicely.
func (p Punishment) FormattedExpiration() string {
	if p.Permanent {
		return "Never"
	}
	return (time.Second * p.Remaining()).String()
}
