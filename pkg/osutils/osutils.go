package osutils

import "runtime"

func GetOsArch() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}
