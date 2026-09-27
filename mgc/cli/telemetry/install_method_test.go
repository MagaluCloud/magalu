package telemetry

import "testing"

func TestInstallMethod(t *testing.T) {
	cases := []struct {
		path string
		want string
	}{
		{"/opt/homebrew/Cellar/mgc/1.4.2/bin/mgc", InstallMethodHomebrew},
		{"/usr/local/Cellar/mgc/1.4.2/bin/mgc", InstallMethodHomebrew},
		{"/home/linuxbrew/.linuxbrew/Cellar/mgc/1.4.2/bin/mgc", InstallMethodHomebrew},

		{"/home/user/.local/bin/mgc", InstallMethodManual},
		{"/usr/local/bin/mgc", InstallMethodManual},
		{"/home/user/Downloads/mgc", InstallMethodManual},
		{"", InstallMethodManual},

		{"/snap/mgc/42/bin/mgc", InstallMethodSnap},

		{"/usr/bin/mgc", InstallMethodSystem},
	}

	for _, c := range cases {
		if got := DetectInstallMethod(c.path); got != c.want {
			t.Errorf("%s: got %s, want %s", c.path, got, c.want)
		}
	}
}
