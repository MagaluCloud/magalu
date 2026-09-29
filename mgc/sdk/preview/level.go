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

// IsLevelName reports whether name is the command of a preview level.
func IsLevelName(name string) bool {
	for _, level := range Levels {
		if level.Name == name {
			return true
		}
	}
	return false
}
