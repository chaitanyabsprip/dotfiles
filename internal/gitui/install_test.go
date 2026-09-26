package gitui

import "testing"

func TestGituiAsset(t *testing.T) {
	tests := []struct{ goos, goarch, want string }{
		{`linux`, `amd64`, `gitui-linux-x86_64.tar.gz`},
		{`linux`, `arm64`, `gitui-linux-aarch64.tar.gz`},
		{`darwin`, `arm64`, `gitui-mac.tar.gz`},
		{`darwin`, `amd64`, `gitui-mac-x86.tar.gz`},
		{`linux`, `386`, ``},
	}
	for _, tt := range tests {
		if got := gituiAsset(`v0.28.1`, tt.goos, tt.goarch); got != tt.want {
			t.Errorf("%s/%s = %q, want %q", tt.goos, tt.goarch, got, tt.want)
		}
	}
}
