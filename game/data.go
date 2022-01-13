package game

type kb struct {
	Horizontal float64
	Vertical   float64
}

type form struct {
	ResourcePath string
}

const (
	NoDebuff = "nodebuff"
	Diamond  = "diamond"
	Build    = "build"
)

var defaultKB = kb{0.398, 0.405}
