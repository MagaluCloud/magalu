package profile_manager

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"testing"

	"slices"

	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/spf13/afero"
)

type testCaseProfileManager struct {
	name          string
	expectedError error
	expectedFs    []utils.TestFsEntry
	providedFs    []utils.TestFsEntry
	run           func(m *ProfileManager) error
}

func checkFs(afs afero.Fs, expected []utils.TestFsEntry) (err error) {
	existingFiles := 0
	err = afero.Walk(afs, "/", func(path string, info fs.FileInfo, e error) (err error) {
		if e != nil {
			return e
		}
		if path == "/" {
			return nil
		}
		fsEntry, err := utils.FindFsEntry(path, expected)
		if err != nil {
			return
		}
		if fsEntry.Mode != info.Mode() {
			return fmt.Errorf("%s: expected mode %x, got %x", path, fsEntry.Mode, info.Mode())
		}
		if fsEntry.Mode&fs.ModeDir == 0 {
			var data []byte
			if data, err = afero.ReadFile(afs, path); err != nil {
				return
			}
			if !bytes.Equal(fsEntry.Data, data) {
				return fmt.Errorf("%s: expected data %q, got %q", path, fsEntry.Data, data)
			}
		}
		existingFiles++
		return nil
	})
	if err != nil {
		return
	}

	if len(expected) != existingFiles {
		return fmt.Errorf("expected %d FS entries, got %d", len(expected), existingFiles)
	}

	return
}

func createProfileManagerGetTest(name string, expectedError error) testCaseProfileManager {
	return testCaseProfileManager{
		name:          fmt.Sprintf("ProfileManager.Get(%q)", name),
		expectedError: expectedError,
		run: func(m *ProfileManager) error {
			p, err := m.Get(name)
			if err != nil {
				return err
			}
			if name != p.Name {
				return fmt.Errorf("expected name %q, got %q", name, p.Name)
			}
			return nil
		},
	}
}

func createProfileManagerCurrentTest(testName string, profileName string, provided []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	return testCaseProfileManager{
		name:       fmt.Sprintf("ProfileManager.Current()==%q[%s]", profileName, testName),
		providedFs: provided,
		expectedFs: provided,
		run: func(m *ProfileManager) error {
			p := m.Current()
			if profileName != p.Name {
				return fmt.Errorf("expected name %q, got %q", profileName, p.Name)
			}
			return nil
		},
	}
}

func createProfileManagerSetCurrentTest(testName string, profileName string, provided []utils.TestFsEntry, expected []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(utils.MergeFsEntries(expected, provided))
	return testCaseProfileManager{
		name:       fmt.Sprintf("ProfileManager.SetCurrent(%q)[%s]", profileName, testName),
		providedFs: provided,
		expectedFs: expected,
		run: func(m *ProfileManager) error {
			p, err := m.Get(profileName)
			if err != nil {
				return err
			}
			return m.SetCurrent(p)
		},
	}
}

func createProfileManagerCreateTest(testName string, profileName string, expectedError error, provided []utils.TestFsEntry, expected []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(utils.MergeFsEntries(expected, provided))
	return testCaseProfileManager{
		name:          fmt.Sprintf("ProfileManager.Create(%q)[%s]", profileName, testName),
		providedFs:    provided,
		expectedFs:    expected,
		expectedError: expectedError,
		run: func(m *ProfileManager) error {
			p, err := m.Create(profileName)
			if err != nil {
				return err
			}
			if profileName != p.Name {
				return fmt.Errorf("expected name %q, got %q", profileName, p.Name)
			}
			return nil
		},
	}
}

func createProfileManagerDeleteTest(testName string, profileName string, expectedError error, provided []utils.TestFsEntry, expected []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(expected)
	return testCaseProfileManager{
		name:          fmt.Sprintf("ProfileManager.Delete(%q)[%s]", profileName, testName),
		providedFs:    provided,
		expectedFs:    expected,
		expectedError: expectedError,
		run: func(m *ProfileManager) error {
			p, err := m.Get(profileName)
			if err != nil {
				return err
			}
			return m.Delete(p)
		},
	}
}

func createProfileManagerListTest(testName string, profileNames []string, provided []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	slices.Sort(profileNames)
	return testCaseProfileManager{
		name:       fmt.Sprintf("ProfileManager.List()=%v[%s]", profileNames, testName),
		providedFs: provided,
		expectedFs: provided,
		run: func(m *ProfileManager) error {
			profiles := m.List()
			if len(profileNames) != len(profiles) {
				return fmt.Errorf("expected %d profiles, got %d", len(profileNames), len(profiles))
			}

			for i, p := range profiles {
				if profileNames[i] != p.Name {
					return fmt.Errorf("expected profile %d to be named %q, got %q", i, profileNames[i], p.Name)
				}
			}

			return nil
		},
	}
}

func createProfileManagerCopyTest(testName string, srcName string, dstName string, expectedError error, provided []utils.TestFsEntry, expected []utils.TestFsEntry) testCaseProfileManager {
	provided = utils.AutoMkdirAll(provided)
	expected = utils.AutoMkdirAll(utils.MergeFsEntries(provided, expected))
	return testCaseProfileManager{
		name:          fmt.Sprintf("ProfileManager.Copy(%q,%q)[%s]", srcName, dstName, testName),
		providedFs:    provided,
		expectedFs:    expected,
		expectedError: expectedError,
		run: func(m *ProfileManager) error {
			src, err := m.Get(srcName)
			if err != nil {
				return err
			}
			dst, err := m.Get(dstName)
			if err != nil {
				return err
			}
			return m.Copy(src, dst)
		},
	}
}

func TestProfileManager(t *testing.T) {
	dir, err := buildMGCPath()
	if err != nil {
		t.Errorf("buildMGCPath: %s", err.Error())
	}

	tests := []testCaseProfileManager{
		// Get()
		createProfileManagerGetTest("default", nil),
		createProfileManagerGetTest("current", errorNameNotAllowed),
		createProfileManagerGetTest("*+&le", errorInvalidName),
		// Current()
		createProfileManagerCurrentTest("empty-fs", defaultProfileName, nil),
		createProfileManagerCurrentTest("empty-file", defaultProfileName, []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte(""),
			},
		}),
		createProfileManagerCurrentTest("provided", "a-profile-name", []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte("a-profile-name"),
			},
		}),
		// SetCurrent()
		createProfileManagerSetCurrentTest("empty-fs", defaultProfileName, nil, []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte(defaultProfileName),
			},
		}),
		createProfileManagerSetCurrentTest("provided", "other-name",
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, currentProfileNameFile),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("a-profile-name"),
				},
			}, []utils.TestFsEntry{
				{
					Path: path.Join(dir, currentProfileNameFile),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("other-name"),
				},
			},
		),
		// Create()
		createProfileManagerCreateTest("empty-fs", defaultProfileName, nil, nil, []utils.TestFsEntry{
			{
				Path: path.Join(dir, defaultProfileName),
				Mode: utils.DIR_PERMISSION | fs.ModeDir,
			},
		}),
		createProfileManagerCreateTest("other-profile", "a-profile-name", nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "other-name"),
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			}, []utils.TestFsEntry{
				{
					Path: path.Join(dir, "a-profile-name"),
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			},
		),
		createProfileManagerCreateTest("existing-profile", "existing-profile", errorProfileAlreadyExists,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "existing-profile"),
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			}, nil,
		),
		createProfileManagerCreateTest("invalid-name", "*+&le", errorInvalidName, nil, nil),
		// Delete()
		createProfileManagerDeleteTest("empty-fs", defaultProfileName, errorDeleteCurrentNotAllowed, nil, nil),
		createProfileManagerDeleteTest("missing-profile", "a-profile-name", nil, nil, nil),
		createProfileManagerDeleteTest("existing-profile", "existing-name", nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "existing-name"),
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			}, []utils.TestFsEntry{
				{
					Path: dir,
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			},
		),
		createProfileManagerDeleteTest("existing-profile-with-files", "existing-name", nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "existing-name/some-file"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("some contents"),
				},
			}, []utils.TestFsEntry{
				{
					Path: dir,
					Mode: utils.DIR_PERMISSION | fs.ModeDir,
				},
			},
		),
		// List()
		createProfileManagerListTest("empty-fs", []string{defaultProfileName}, nil),
		createProfileManagerListTest("missing-current", []string{"a-profile-name"}, []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte("a-profile-name"),
			},
		}),
		createProfileManagerListTest("existing-current-should-not-duplicate", []string{"a-profile-name"}, []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte("a-profile-name"),
			},
			{
				Path: path.Join(dir, "a-profile-name"),
				Mode: utils.DIR_PERMISSION | fs.ModeDir,
			},
		}),
		createProfileManagerListTest("multiple", []string{defaultProfileName, "a-profile-name", "other-name"}, []utils.TestFsEntry{
			{
				Path: path.Join(dir, currentProfileNameFile),
				Mode: utils.FILE_PERMISSION,
				Data: []byte("a-profile-name"),
			},
			{
				Path: path.Join(dir, defaultProfileName),
				Mode: utils.DIR_PERMISSION | fs.ModeDir,
			},
			{
				Path: path.Join(dir, "other-name"),
				Mode: utils.DIR_PERMISSION | fs.ModeDir,
			},
		}),
		// Copy()
		createProfileManagerCopyTest("empty-fs", defaultProfileName, "a-profile-name", nil, nil, nil),
		createProfileManagerCopyTest("copy", defaultProfileName, "a-profile-name", nil,
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, defaultProfileName, "some-file"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("contents-here"),
				},
				{
					Path: path.Join(dir, defaultProfileName, "other-file"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("other-contents-here"),
				},
			},
			[]utils.TestFsEntry{
				{
					Path: path.Join(dir, "a-profile-name", "some-file"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("contents-here"),
				},
				{
					Path: path.Join(dir, "a-profile-name", "other-file"),
					Mode: utils.FILE_PERMISSION,
					Data: []byte("other-contents-here"),
				},
			},
		),
		createProfileManagerCopyTest("self-copy", defaultProfileName, defaultProfileName, errorCopyToSelf, nil, nil),
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			m := &ProfileManager{dir, fs}
			err := utils.PrepareFs(fs, tc.providedFs)
			if err != nil {
				t.Errorf("could not prepare provided FS: %s", err.Error())
			}
			err = tc.run(m)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected err == %#v, found: %#v", tc.expectedError, err)
			}
			err = checkFs(fs, tc.expectedFs)
			if err != nil {
				t.Errorf("unexpected FS state: %s", err.Error())
			}
		})
	}
}
