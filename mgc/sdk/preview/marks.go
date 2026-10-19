package preview

import (
	"slices"

	"github.com/MagaluCloud/magalu/mgc/core"
)

const (
	NewMark          = "[PREVIEW]"
	NewCommandsTrail = "(preview commands)"
	NewFlagsTrail    = "(preview flags)"

	// The group of the official products, as opposed to settings and levels.
	productsGroupID = "catalog"
)

// FlagDiff lists the flags a preview command has that its official one has
// not, named the way the help shows them.
type FlagDiff func(preview, official core.Executor) []string

// Marks tells the help what a preview module brings that the official product
// of the same name has not. A command is given as the chain of descriptors
// from its level down, such as alpha, iam, roles, list.
type Marks struct {
	root    func() core.Grouper
	diff    FlagDiff
	modules map[string]*moduleMarks
}

type moduleMarks struct {
	// Without an official product, the module is new as a whole and nothing
	// inside it is marked: marking everything tells nothing.
	compared   bool
	changes    []change
	everywhere []string
}

type change struct {
	// Descriptor names from the module down.
	path []string
	// The topmost new command or group; nothing under it is listed.
	isNew bool
	flags []string
}

func NewMarks(root func() core.Grouper, diff FlagDiff) *Marks {
	return &Marks{root: root, diff: diff, modules: map[string]*moduleMarks{}}
}

// Mark tells how the help shows a command or group: new when the official
// product lacks it, or with a trail when something under it is new.
func (m *Marks) Mark(chain []core.Descriptor) (isNew bool, trail string) {
	module, path := m.module(chain)
	if module == nil {
		return false, ""
	}
	if !module.compared {
		return len(path) == 1, ""
	}

	for _, c := range module.changes {
		if !isUnder(c.path, path) {
			continue
		}
		switch {
		case c.isNew && len(c.path) == len(path):
			return true, ""
		case c.isNew:
			trail = NewCommandsTrail
		case trail == "" && len(without(c.flags, module.everywhere)) > 0:
			trail = NewFlagsTrail
		}
	}
	return false, trail
}

// NewFlags lists the flags of a command that the official one has not. A new
// command has none: all of it is new.
func (m *Marks) NewFlags(chain []core.Descriptor) []string {
	module, path := m.module(chain)
	if module == nil {
		return nil
	}
	for _, c := range module.changes {
		if !c.isNew && slices.Equal(c.path, path) {
			return c.flags
		}
	}
	return nil
}

// Everywhere lists the new flags every command of the module has, which the
// trail leaves out and the help of the module tells once.
func (m *Marks) Everywhere(chain []core.Descriptor) []string {
	if len(chain) != 2 {
		return nil
	}
	if module, _ := m.module(chain); module != nil {
		return module.everywhere
	}
	return nil
}

func (m *Marks) module(chain []core.Descriptor) (*moduleMarks, []string) {
	if len(chain) < 2 || !IsLevel(chain[0].Name()) {
		return nil, nil
	}
	group, ok := chain[1].(core.Grouper)
	if !ok {
		return nil, nil
	}

	path := make([]string, len(chain)-1)
	for i, desc := range chain[1:] {
		path[i] = desc.Name()
	}

	key := chain[0].Name() + "/" + group.Name()
	module, ok := m.modules[key]
	if !ok {
		module = m.compareModule(group)
		m.modules[key] = module
	}
	return module, path
}

// A mark is not worth breaking the help for: a module that fails to compare
// has none, not even the one of a module the official tree lacks.
func (m *Marks) compareModule(group core.Grouper) *moduleMarks {
	official := m.official(group.Name())
	if official == nil {
		return &moduleMarks{}
	}

	module := &moduleMarks{compared: true}
	var commands [][]string
	err := m.compare(group, official, []string{group.Name()}, module, &commands)
	if err != nil {
		logger().Debugw("not marking what is new", "module", group.Name(), "error", err)
		return &moduleMarks{compared: true}
	}
	module.everywhere = inEvery(commands)
	return module
}

func (m *Marks) official(name string) core.Descriptor {
	root := m.root()
	if root == nil {
		return nil
	}
	product, err := root.GetChildByName(name)
	if err != nil || product.GroupID() != productsGroupID {
		return nil
	}
	return product
}

func (m *Marks) compare(desc, official core.Descriptor, path []string, module *moduleMarks, commands *[][]string) error {
	if desc.IsInternal() {
		return nil
	}

	switch desc := desc.(type) {
	case core.Grouper:
		official, ok := official.(core.Grouper)
		if !ok {
			module.changes = append(module.changes, change{path: path, isNew: true})
			return nil
		}
		_, err := desc.VisitChildren(func(child core.Descriptor) (bool, error) {
			counterpart, _ := official.GetChildByName(child.Name())
			return true, m.compare(child, counterpart, append(slices.Clip(path), child.Name()), module, commands)
		})
		return err

	case core.Executor:
		official, ok := official.(core.Executor)
		if !ok {
			module.changes = append(module.changes, change{path: path, isNew: true})
			return nil
		}
		flags := m.diff(desc, official)
		*commands = append(*commands, flags)
		if len(flags) > 0 {
			module.changes = append(module.changes, change{path: path, flags: flags})
		}
	}
	return nil
}

// isUnder tells whether path is a group above, or the command at, the given
// place.
func isUnder(path, place []string) bool {
	return len(path) >= len(place) && slices.Equal(path[:len(place)], place)
}

// The flags new in every command, when there is more than one to tell by.
func inEvery(commands [][]string) []string {
	if len(commands) < 2 {
		return nil
	}
	every := slices.Clone(commands[0])
	for _, flags := range commands[1:] {
		every = slices.DeleteFunc(every, func(flag string) bool { return !slices.Contains(flags, flag) })
	}
	return every
}

func without(flags, left []string) []string {
	return slices.DeleteFunc(slices.Clone(flags), func(flag string) bool { return slices.Contains(left, flag) })
}

func IsLevel(name string) bool {
	return slices.ContainsFunc(Levels, func(level Level) bool { return level.Name == name })
}
