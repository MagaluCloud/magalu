package telemetry

// NoopDispatcher descarta o evento sem erro
type NoopDispatcher struct{}

func (NoopDispatcher) Dispatch(Event) error {
	return nil
}
