package version

import "testing"

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		count    int
		arch     string
		isDirty  bool
		expected string
	}{
		{10, "amd64", false, "b10-amd64"},
		{10, "amd64", true, "b10-amd64-dev"},
		{0, "arm64", false, "b0-arm64"},
		{145, "arm64", true, "b145-arm64-dev"},
	}

	for _, tc := range tests {
		result := FormatVersion(tc.count, tc.arch, tc.isDirty)
		if result != tc.expected {
			t.Errorf("FormatVersion(%d, %s, %t) = %s; want %s", tc.count, tc.arch, tc.isDirty, result, tc.expected)
		}
	}
}
