package telemetry

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

const (
	StateFileName = "telemetry.yaml"

	keyDisabled           = "telemetry.disabled"
	keyNoticeShown        = "telemetry.notice_shown"
	keyCredentialsSetAt   = "telemetry.credentials_set_at"
	keyFirstValueRecorded = "telemetry.first_value_recorded"
	keyInstallationID     = "telemetry.installation_id"
)

var _ StateStore = (*ViperStateStore)(nil)

type ViperStateStore struct {
	path string
}

func NewViperStateStore(dir string) *ViperStateStore {
	return &ViperStateStore{path: filepath.Join(dir, StateFileName)}
}

func (s *ViperStateStore) Path() string {
	return s.path
}

func (s *ViperStateStore) newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(s.path)
	v.SetConfigType("yaml")
	return v
}

func (s *ViperStateStore) Load() (State, error) {
	v := s.newViper()
	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return State{}, nil
		}
		return State{}, err
	}

	state := State{
		Disabled:           v.GetBool(keyDisabled),
		NoticeShown:        v.GetBool(keyNoticeShown),
		FirstValueRecorded: v.GetBool(keyFirstValueRecorded),
		InstallationID:     v.GetString(keyInstallationID),
	}

	if raw := v.GetString(keyCredentialsSetAt); raw != "" {
		if t, err := time.Parse(time.RFC3339Nano, raw); err == nil {
			t = t.UTC()
			state.CredentialsSetAt = &t
		}
	}

	return state, nil
}

func (s *ViperStateStore) Save(state State) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	v := s.newViper()
	v.Set(keyDisabled, state.Disabled)
	v.Set(keyNoticeShown, state.NoticeShown)
	v.Set(keyFirstValueRecorded, state.FirstValueRecorded)
	if state.InstallationID != "" {
		v.Set(keyInstallationID, state.InstallationID)
	}
	if state.CredentialsSetAt != nil {
		v.Set(keyCredentialsSetAt, state.CredentialsSetAt.UTC().Format(time.RFC3339Nano))
	}

	return s.writeAtomically(v)
}

// writeAtomically grava num arquivo temporário do mesmo diretório e depois renomeia
// para o destino. Assim, execuções simultâneas ou um processo interrompido no meio
// nunca deixam o telemetry.yaml corrompido
func (s *ViperStateStore) writeAtomically(v *viper.Viper) error {
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".telemetry-*.yaml")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := tmp.Close(); err != nil {
		return err
	}
	if err := v.WriteConfigAs(tmpPath); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}
