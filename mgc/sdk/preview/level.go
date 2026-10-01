package preview

// Level is a preview namespace, such as "mgc beta", and the text it shows.
type Level struct {
	Name    string
	Summary string
}

var (
	Beta = Level{
		Name:    "beta",
		Summary: "Preview commands released",
	}
	Alpha = Level{
		Name:    "alpha",
		Summary: "Early preview commands",
	}
)

// Levels lists every preview level, in the order they are merged.
var Levels = []Level{Beta, Alpha}
