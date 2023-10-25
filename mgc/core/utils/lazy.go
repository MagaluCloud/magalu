package utils

import "errors"

type LoadWithError[T any] func() (T, error)
type OnceWithError func() error

var ErrorNotLoaded = errors.New("not loaded")

func NewLazyLoaderWithError[T any](loader LoadWithError[T]) LoadWithError[T] {
	var value T
	var err error = ErrorNotLoaded
	return func() (T, error) {
		if loader != nil {
			value, err = loader()
			loader = nil
		}
		return value, err
	}
}

func NewLazyOnceWithError(loader OnceWithError) OnceWithError {
	var err error = ErrorNotLoaded
	return func() error {
		if loader != nil {
			err = loader()
			loader = nil
		}
		return err
	}
}

func NewLazyLoader[T any](loader func() T) func() T {
	var value T
	return func() T {
		if loader != nil {
			value = loader()
			loader = nil
		}
		return value
	}
}

func NewLazyOnce(loader func()) func() {
	return func() {
		if loader != nil {
			loader()
			loader = nil
		}
	}
}
