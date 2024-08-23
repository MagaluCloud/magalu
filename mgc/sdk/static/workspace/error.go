package workspace

import "fmt"

type WorkspaceError struct {
	Name string
	Err  error
}

func (e WorkspaceError) Unwrap() error {
	return e.Err
}

func (e WorkspaceError) Error() string {
	return fmt.Sprintf("workspace %s: %s", e.Name, e.Err.Error())
}
