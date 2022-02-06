package perm

const (
	Owner     = "Owner"
	Admin     = "Admin"
	Mod       = "Mod"
	Famous    = "Famous"
	YouTube   = "YouTube"
	Hyperedge = "Hyperedge"
	Ravager   = "Ravager"
	VIP       = "VIP"
)

var ranks = map[string]*Rank{
	Owner: {
		Color:      "§4",
		ChatFormat: "§o§4Owner §r§a%v: §f%v",
	},
	Admin: {
		Color:      "§d",
		ChatFormat: "Admin §r§a%v: §f%v",
	},
	Mod: {
		Color:      "§2",
		ChatFormat: "Mod §r§a%v: §f%v",
	},
	Famous: {
		Color:      "§c",
		ChatFormat: "Famous §r§a%v: §f%v",
	},
	YouTube: {
		Color:      "§b",
		ChatFormat: "Media §r§a%v: §f%v",
	},
	Hyperedge: {
		Color:      "§6",
		ChatFormat: "Hyperedge §r§a%v: §f%v",
	},
	Ravager: {
		Color:      "§4",
		ChatFormat: "Ravager §r§a%v: §f%v",
	},
	VIP: {
		Color:      "§9",
		ChatFormat: "VIP §r§a%v: §f%v",
	},
}

func init() {
	for name, rank := range ranks {
		rank.Name = name
		rank.ChatFormat = "§8[" + rank.Color + name + "§8] " + "§r§a%v: §f%v"
	}
}

// GetRank will return the rank based on the name or nil if it doesn't exist.
func GetRank(name string) *Rank {
	if r, ok := ranks[name]; ok {
		return r
	}
	return nil
}

// Ranks will return a list of all ranks.
func Ranks() map[string]*Rank {
	return ranks
}
