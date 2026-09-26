package yq

import "testing"

func TestYqAsset(t *testing.T) {
	tests := []struct{ goos, goarch, want string }{
		{`linux`, `amd64`, `yq_linux_amd64`},
		{`linux`, `arm64`, `yq_linux_arm64`},
		{`darwin`, `arm64`, `yq_darwin_arm64`},
		{`windows`, `amd64`, ``},
		{`linux`, `386`, ``},
	}
	for _, tt := range tests {
		if got := yqAsset(`v4.53.6`, tt.goos, tt.goarch); got != tt.want {
			t.Errorf("%s/%s = %q, want %q", tt.goos, tt.goarch, got, tt.want)
		}
	}
}
