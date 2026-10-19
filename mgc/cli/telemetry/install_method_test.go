package telemetry

import "testing"

func TestInstallMethod(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if got := DetectInstallMethod(tc.path); got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
