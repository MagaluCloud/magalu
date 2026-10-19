package viperstate

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/MagaluCloud/magalu/mgc/cli/telemetry"
	"github.com/spf13/viper"
)

const (
	keyDisabled           = "telemetry.disabled"
	keyNoticeShown        = "telemetry.notice_shown"
	keyCredentialsSetAt   = "telemetry.credentials_set_at"
	keyFirstValueRecorded = "telemetry.first_value_recorded"
	keyInstallationID     = "telemetry.installation_id"
)

var _ telemetry.StateStore = (*Store)(nil)

type Store struct {
	path string
}

func NewStore(dir string) *Store {
	return &Store{path: filepath.Join(dir, telemetry.StateFileName)}
}

func (s *Store) Path() string {
	return s.path
}

func (s *Store) newViper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(s.path)
	v.SetConfigType("yaml")
	return v
}

func (s *Store) Load() (telemetry.State, error) {
	v := s.newViper()
	if err := v.ReadInConfig(); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return telemetry.State{}, nil
		}
		return telemetry.State{}, err
	}

	state := telemetry.State{
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

func (s *Store) Save(state telemetry.State) error {
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
func (s *Store) writeAtomically(v *viper.Viper) error {
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
