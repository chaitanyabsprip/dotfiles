package eza

import "testing"

func TestEzaAsset(t *testing.T) {
	tests := []struct{ goos, goarch, want string }{
		{`linux`, `amd64`, `eza_x86_64-unknown-linux-gnu.tar.gz`},
		{`linux`, `arm64`, `eza_aarch64-unknown-linux-gnu.tar.gz`},
		{`darwin`, `arm64`, ``}, // no macOS build; brew covers it
		{`linux`, `386`, ``},
	}
	for _, tt := range tests {
		if got := ezaAsset(`v0.23.5`, tt.goos, tt.goarch); got != tt.want {
			t.Errorf("%s/%s = %q, want %q", tt.goos, tt.goarch, got, tt.want)
		}
	}
}
