package session

const vanishPrefix = "§c[V]§r "

// vanish is the custom gamemode for vanished players
type vanish struct{}

func (vanish) AllowsEditing() bool      { return true }
func (vanish) AllowsTakingDamage() bool { return false }
func (vanish) CreativeInventory() bool  { return false }
func (vanish) HasCollision() bool       { return false }
func (vanish) AllowsFlying() bool       { return true }
func (vanish) AllowsInteraction() bool  { return true }
func (vanish) Visible() bool            { return true }

// todo: fix vanish!
