//go:build !release

package cmd

func logValidationErr(e error) {
	logger().Warnw("result validation failed", "error", e.Error())
}
