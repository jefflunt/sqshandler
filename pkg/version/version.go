package version

import "fmt"

// Version holds the current compiled version of the sqshandler CLI.
// It is injected at build time using -ldflags.
var Version = "dev"

// FormatVersion returns the standardized version string format.
func FormatVersion(commitCount int, arch string, isDirty bool) string {
	v := fmt.Sprintf("b%d-%s", commitCount, arch)
	if isDirty {
		v += "-dev"
	}
	return v
}
