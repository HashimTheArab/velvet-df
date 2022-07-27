package healing

type source struct{}

type (
	// SourceKit is a custom healing source applied when a player is healed by a kit.
	SourceKit struct{ source }
	// SourceKill is a custom healing source applied when a player is healed by killing another player.
	SourceKill struct{ source }
	// SourceDeath is a custom healing source applied when a player is killed, this is because we don't actually kill the player.
	SourceDeath struct{ source }
)

// HealingSource ...
func (source) HealingSource() {}
