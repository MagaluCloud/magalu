package core

type AuthenticationError struct {
	Err error
}

func (e AuthenticationError) Unwrap() error {
	return e.Err
}

func (e AuthenticationError) Error() string {
	return e.Err.Error()
}
