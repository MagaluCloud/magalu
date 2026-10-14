package preview

type Level struct {
	Name    string
	Summary string
}

var (
	Beta = Level{
		Name:    "beta",
		Summary: "Commands in public preview",
	}
	Alpha = Level{
		Name:    "alpha",
		Summary: "Early preview commands",
	}
)

var Levels = []Level{Beta, Alpha}
