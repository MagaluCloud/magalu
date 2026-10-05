package telemetry

import "testing"

func TestCommandInfoAction(t *testing.T) {
	testCases := []struct {
		name string
		info CommandInfo
		want string
	}{
		{"nested command", CommandInfo{Path: []string{"virtual-machine", "instances", "list"}}, "virtualmachine.instances.list"},
		{"single segment", CommandInfo{Path: []string{"version"}}, "version"},
		{"upper case and hyphens", CommandInfo{Path: []string{"Object-Storage", "Buckets", "Create"}}, "objectstorage.buckets.create"},
		{"multiple hyphens", CommandInfo{Path: []string{"block-storage", "volume-attachments", "attach"}}, "blockstorage.volumeattachments.attach"},
		{"empty path", CommandInfo{}, UnknownAction},
		{"unknown command", CommandInfo{Path: []string{"virtual-machine"}, UnknownCommand: true}, UnknownAction},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.info.Action(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCommandInfoProduct(t *testing.T) {
	testCases := []struct {
		name string
		info CommandInfo
		want string
	}{
		{"nested command", CommandInfo{Path: []string{"virtual-machine", "images", "list"}}, "virtual machine"},
		{"upper case and hyphens", CommandInfo{Path: []string{"Object-Storage", "Buckets", "Create"}}, "object storage"},
		{"settings command", CommandInfo{Path: []string{"config", "list"}}, "config"},
		{"empty path", CommandInfo{}, UnknownAction},
		{"unknown command", CommandInfo{Path: []string{"virtual-machine"}, UnknownCommand: true}, UnknownAction},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.info.Product(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCommandInfoResourceType(t *testing.T) {
	testCases := []struct {
		name string
		info CommandInfo
		want string
	}{
		{"nested command", CommandInfo{Path: []string{"virtual-machine", "instances", "list"}}, "virtual_machine_instances"},
		{"two segments", CommandInfo{Path: []string{"auth", "login"}}, "auth"},
		{"upper case", CommandInfo{Path: []string{"Object-Storage", "Buckets", "Create"}}, "object_storage_buckets"},
		{"single segment has no resource", CommandInfo{Path: []string{"version"}}, ""},
		{"empty path", CommandInfo{}, ""},
		{"unknown command", CommandInfo{Path: []string{"virtual-machine", "instances", "list"}, UnknownCommand: true}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.info.ResourceType(); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCommandInfoIsLogin(t *testing.T) {
	testCases := []struct {
		name string
		info CommandInfo
		want bool
	}{
		{"auth login", CommandInfo{Path: []string{"auth", "login"}}, true},
		{"auth logout", CommandInfo{Path: []string{"auth", "logout"}}, false},
		{"auth group only", CommandInfo{Path: []string{"auth"}}, false},
		{"deeper login path", CommandInfo{Path: []string{"auth", "login", "extra"}}, false},
		{"login outside auth", CommandInfo{Path: []string{"other", "login"}}, false},
		{"unknown command", CommandInfo{Path: []string{"auth", "login"}, UnknownCommand: true}, false},
		{"empty path", CommandInfo{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.info.IsLogin(); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCommandInfoIsInfrastructure(t *testing.T) {
	testCases := []struct {
		name string
		info CommandInfo
		want bool
	}{
		{"product command", CommandInfo{Path: []string{"virtual-machine", "instances", "list"}}, true},
		{"another product", CommandInfo{Path: []string{"object-storage", "buckets", "list"}}, true},
		{"auth", CommandInfo{Path: []string{"auth", "login"}}, false},
		{"config", CommandInfo{Path: []string{"config", "set"}}, false},
		{"profile", CommandInfo{Path: []string{"profile", "ssh-keys", "list"}}, false},
		{"workspace", CommandInfo{Path: []string{"workspace", "list"}}, false},
		{"telemetry", CommandInfo{Path: []string{"telemetry", "status"}}, false},
		{"unknown command", CommandInfo{Path: []string{"virtual-machine"}, UnknownCommand: true}, false},
		{"empty path", CommandInfo{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.info.IsInfrastructure(); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
