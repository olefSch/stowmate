package main

import "testing"

func TestFormatVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
		want    string
	}{
		{
			name:    "all set",
			version: "1.0.0",
			commit:  "abc1234",
			date:    "2026-01-01T00:00:00Z",
			want:    "1.0.0 (abc1234, 2026-01-01T00:00:00Z)",
		},
		{
			name:    "only version",
			version: "dev",
			want:    "dev",
		},
		{
			name:    "version and commit",
			version: "1.0.0",
			commit:  "abc1234",
			want:    "1.0.0 (abc1234)",
		},
		{
			name:    "version and date",
			version: "1.0.0",
			date:    "2026-01-01T00:00:00Z",
			want:    "1.0.0 (2026-01-01T00:00:00Z)",
		},
		{
			name: "all empty",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatVersion(tt.version, tt.commit, tt.date); got != tt.want {
				t.Errorf("formatVersion(%q, %q, %q) = %q, want %q", tt.version, tt.commit, tt.date, got, tt.want)
			}
		})
	}
}
