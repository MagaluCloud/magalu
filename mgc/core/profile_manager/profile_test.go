package profile_manager

import (
	"errors"
	"fmt"

	"io/fs"
	"path"
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"magalu.cloud/core/utils"
)

type testCaseProfile struct {
	name          string
	expectedData  []byte
	expectedError error
	expectedFs    []utils.TestFsEntry
	providedFs    []utils.TestFsEntry
	run           func(p *Profile) ([]byte, error)
}

func createProfileWriteTest(testName, name string, data []byte, expectedError error, provided, expected []utils.TestFsEntry) testCaseProfile {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(expected)
	return testCaseProfile{
		name:          fmt.Sprintf("profile.Write(%q)[%s]", name, testName),
		expectedError: expectedError,
		expectedFs:    expected,
		providedFs:    provided,
		run: func(p *Profile) ([]byte, error) {
			return nil, p.Write(name, data)
		},
	}
}

func createProfileReadTest(testName, name string, expectedData []byte, expectedError error, provided []utils.TestFsEntry) testCaseProfile {
	provided = utils.AutoMkdirAll(provided)
	return testCaseProfile{
		name:          fmt.Sprintf("profile.Read(%q)[%s]", name, testName),
		expectedData:  expectedData,
		expectedError: expectedError,
		expectedFs:    provided,
		providedFs:    provided,
		run: func(p *Profile) ([]byte, error) {
			return p.Read(name)
		},
	}
}

func createDeleteProfileTest(testName, name string, expectedError error, provided, expected []utils.TestFsEntry) testCaseProfile {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(expected)
	return testCaseProfile{
		name:          fmt.Sprintf("profile.Delete(%q)[%s]", name, testName),
		expectedError: expectedError,
		expectedFs:    expected,
		providedFs:    provided,
		run: func(p *Profile) ([]byte, error) {
			return nil, p.Delete(name)
		},
	}
}

func TestProfile(t *testing.T) {
	dir, err := buildMGCPath()
	if err != nil {
		t.Errorf("buildMGCPath: %s", err.Error())
	}

	tests := []testCaseProfile{
		// Write()
		createProfileWriteTest("empty-fs", "test", []byte("test"), nil, nil, []utils.TestFsEntry{
			{
				Path: path.Join(dir, "profile/test"),
				Mode: utils.FILE_PERMISSION,
				Data: []byte("test"),
			},
		}),
		createProfileWriteTest("existing file", "test", []byte("updated test"), nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("test"),
				},
			},
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("updated test"),
				},
			},
		),
		createProfileWriteTest("new file", "other-test", []byte("other-test"), nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("test"),
				},
			},
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("test"),
				},
				{
					Path: path.Join(dir, "profile/other-test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("other-test"),
				},
			},
		),
		// Read()
		createProfileReadTest("empty-fs", "test", nil, afero.ErrFileNotFound, nil),
		createProfileReadTest("existing file", "test", []byte("test"), nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("test"),
				},
			},
		),
		// Delete()
		createDeleteProfileTest("empty-fs", "test", nil, nil, nil),
		createDeleteProfileTest("existing file", "test", nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile/test"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("test"),
				},
			},
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "profile"),
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			},
		),
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			m := &ProfileManager{dir, fs}

			p, err := m.Get("profile")
			if err != nil {
				t.Errorf("ProfileManager.Get() failed: %s", err.Error())
			}

			err = utils.PrepareFs(fs, tc.providedFs)
			if err != nil {
				t.Errorf("could not prepare provided FS: %s", err.Error())
			}
			data, err := tc.run(p)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected err == %#v, found: %#v", tc.expectedError, err)
			}
			if !reflect.DeepEqual(data, tc.expectedData) {
				t.Errorf("expected data == %s, found: %s", string(tc.expectedData), string(data))
			}
			err = checkFs(fs, tc.expectedFs)
			if err != nil {
				t.Errorf("unexpected FS state: %s", err.Error())
			}
		})
	}
}
