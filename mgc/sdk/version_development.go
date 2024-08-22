//go:build !release

package sdk

// in order to use this, build with -buildvcs=true

var Version = func() string {
	return version + "-dev"
}()
