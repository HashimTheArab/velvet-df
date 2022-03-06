package session

const (
	FlagAdmin uint32 = 1 << iota
	FlagBuilder
	FlagStaff
	FlagVanished
	FlagBuilding
)
