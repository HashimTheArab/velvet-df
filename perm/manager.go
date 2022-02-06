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
	Owner:     {Color: "§4"},
	Admin:     {Color: "§d"},
	Mod:       {Color: "§2"},
	Famous:    {Color: "§c"},
	YouTube:   {Color: "§b"},
	Hyperedge: {Color: "§6"},
	Ravager:   {Color: "§4"},
	VIP:       {Color: "§9"},
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