package db

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/upper/db/v4"
	"strings"
	"velvet/perm"
	"velvet/session"
)

// Entry is a player database entry.
type Entry struct {
	// XUID is the xuid of the player.
	XUID string `bson:"xuid"`
	// DisplayName is the display name of the player.
	DisplayName string `bson:"display_name"`
	// Name is the lowercase name of the player.
	Name string `bson:"name"`
	// DeviceID is the device id of the player.
	DeviceID string `bson:"device_id"`
	// Rank is the rank of the player.
	Rank string `bson:"rank"`
	// Kills are the amount of players the player has killed.
	Kills uint32 `bson:"kills"`
	// Deaths are the amount of times the player has died.
	Deaths uint32 `bson:"deaths"`
	// Punishments contains all the punishments of the player.
	Punishments Punishments `bson:"punishments"`
}

// Punishments contains the punishments of a player.
type Punishments struct {
	Ban  Punishment `bson:"ban"`
	Mute Punishment `bson:"mute"`
}

// Register registers a player into the database.
func Register(xuid, displayName, deviceID string) {
	entry := sess.Collection("players").Find(db.Cond{"xuid": xuid})
	if ok, _ := entry.Exists(); ok {
		_ = entry.Update(map[string]string{"display_name": displayName, "name": strings.ToLower(displayName), "device_id": deviceID})
	} else {
		_, _ = sess.Collection("players").Insert(Entry{
			XUID:        xuid,
			DisplayName: displayName,
			Name:        strings.ToLower(displayName),
			DeviceID:    deviceID,
		})
	}
}

// LoadSession loads a user session from the database.
func LoadSession(p *player.Player) (*session.Session, error) {
	var entry *Entry
	if err := findPlayer(p.Name()).One(entry); err != nil {
		return nil, err
	}
	return session.New(p,
		perm.GetRank(entry.Rank),
		entry.Kills,
		entry.Deaths,
	), nil
}

// SaveSession saves a user session to the database.
func SaveSession(session *session.Session) error {
	return SaveOfflinePlayer(&Entry{
		XUID:        session.Player.XUID(),
		DisplayName: session.Player.Name(),
		Name:        strings.ToLower(session.Player.Name()),
		DeviceID:    session.DeviceID(),
		Rank:        session.RankName(),
		Kills:       session.Kills(),
		Deaths:      session.Deaths(),
		Punishments: Punishments{}, // todo
	})
}

// Registered returns whether a player is registered
func Registered(id string) bool {
	ok, _ := findPlayer(id).Exists()
	return ok
}

// LoadOfflinePlayer returns an offline player entry for the given ign, if the player does not exist, an error will be returned.
func LoadOfflinePlayer(ign string) (*Entry, error) {
	var entry *Entry
	err := findPlayer(ign).One(entry)
	return entry, err
}

// LoadOfflinePlayers returns all player entries that match the given conditions.
func LoadOfflinePlayers(cond ...any) ([]*Entry, error) {
	var data []*Entry
	err := sess.Collection("players").Find(cond...).All(&data)
	return data, err
}

// SaveOfflinePlayer saves the entry of an offline player.
func SaveOfflinePlayer(entry *Entry) error {
	players := sess.Collection("players")
	found := players.Find(db.Cond{"xuid": entry.XUID})
	if ok, _ := found.Exists(); ok {
		return found.Update(entry)
	}
	_, err := players.Insert(entry)
	return err
}

// GetAlias will return all the names that have the same deviceID as the given ign.
// Zero values will be returned if the player has never joined before.
func GetAlias(ign string) (deviceID string, names []string) {
	p, err := LoadOfflinePlayer(ign)
	if err != nil {
		return
	}
	deviceID = p.DeviceID
	entries, _ := LoadOfflinePlayers(db.Cond{"device_id": deviceID})
	for _, e := range entries {
		name := "§e" + e.DisplayName
		if !e.Punishments.Ban.Expired() {
			name += " §l§cBANNED§r"
		}
		names = append(names, name)
	}
	return
}

// IsStaff will return whether a player has a staff rank.
func IsStaff(id string) bool {
	p, err := LoadOfflinePlayer(id)
	return err == nil && perm.StaffRanks.Contains(p.Rank)
}

// findPlayer is used internally to fetch a player entry from the database.
func findPlayer(id string) db.Result {
	return sess.Collection("players").Find(
		db.Or(db.Cond{"xuid": id}, db.Cond{"name": strings.ToLower(id)}),
	)
}
