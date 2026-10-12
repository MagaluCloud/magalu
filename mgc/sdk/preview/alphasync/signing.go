package alphasync

import (
	"crypto/ed25519"
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
)

const signedPrefix = "mgc-alpha-spec-v1\n"

//go:embed trusted_keys.json
var trustedKeysJSON []byte

type Keys map[string]ed25519.PublicKey

func TrustedKeys() Keys {
	keys, err := ParseKeys(trustedKeysJSON)
	if err != nil {
		panic(err)
	}
	return keys
}

func ParseKeys(data []byte) (Keys, error) {
	var encoded map[string]string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return nil, fmt.Errorf("trusted keys do not parse: %w", err)
	}

	keys := Keys{}
	for id, value := range encoded {
		key, err := hex.DecodeString(value)
		if err != nil || len(key) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("trusted key %q is not an Ed25519 public key in hex", id)
		}
		keys[id] = ed25519.PublicKey(key)
	}
	return keys, nil
}

type signedFields struct {
	KeyID         string `json:"key_id"`
	Module        string `json:"module"`
	SHA256        string `json:"sha256"`
	MinCLIVersion string `json:"min_cli_version"`
}

func SignedPayload(spec ManifestSpec) []byte {
	fields, err := json.Marshal(signedFields{
		KeyID:         spec.KeyID,
		Module:        spec.Module,
		SHA256:        strings.ToLower(spec.SHA256),
		MinCLIVersion: spec.MinCLIVersion,
	})
	if err != nil {
		panic(err)
	}
	return append([]byte(signedPrefix), fields...)
}

func Sign(key ed25519.PrivateKey, keyID string, spec ManifestSpec) ManifestSpec {
	spec.KeyID = keyID
	spec.SHA256 = strings.ToLower(spec.SHA256)
	spec.Signature = base64.StdEncoding.EncodeToString(ed25519.Sign(key, SignedPayload(spec)))
	return spec
}

func (k Keys) Verify(spec ManifestSpec) error {
	if spec.KeyID == "" || spec.Signature == "" {
		return fmt.Errorf("alpha module %q is not signed", spec.Module)
	}
	key, ok := k[spec.KeyID]
	if !ok {
		return fmt.Errorf("alpha module %q is signed with key %q, which this mgc does not know: update mgc", spec.Module, spec.KeyID)
	}
	signature, err := base64.StdEncoding.DecodeString(spec.Signature)
	if err != nil || !ed25519.Verify(key, SignedPayload(spec), signature) {
		return fmt.Errorf("alpha module %q has an invalid signature", spec.Module)
	}
	return nil
}

func (k Keys) KeyIDOf(private ed25519.PrivateKey) (string, error) {
	public, ok := private.Public().(ed25519.PublicKey)
	if !ok {
		return "", errors.New("not an Ed25519 private key")
	}
	for id, key := range k {
		if key.Equal(public) {
			return id, nil
		}
	}
	return "", errors.New("this private key does not pair with any trusted key")
}

func VerifiedStore(store preview.Store, dir string, keys Keys) preview.Store {
	return &verifiedStore{Store: store, dir: dir, keys: keys}
}

type verifiedStore struct {
	preview.Store
	dir  string
	keys Keys
}

func (s *verifiedStore) Read(name string) (preview.Package, error) {
	pkg, err := s.Store.Read(name)
	if err != nil {
		return pkg, err
	}

	index, err := readIndex(s.dir)
	if err != nil {
		return preview.Package{}, err
	}
	spec, ok := index[name]
	if !ok || spec.Module != name {
		return preview.Package{}, errors.New("it was not installed by 'mgc alpha sync': run it again")
	}
	if err := s.keys.Verify(spec); err != nil {
		return preview.Package{}, err
	}

	sum := sha256.Sum256(pkg.Spec)
	if hex.EncodeToString(sum[:]) != strings.ToLower(spec.SHA256) {
		return preview.Package{}, errors.New("it changed since it was synced: run 'mgc alpha sync'")
	}
	return pkg, nil
}

func readIndex(dir string) (map[string]ManifestSpec, error) {
	file := filepath.Join(dir, indexFile)
	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var index map[string]ManifestSpec
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("%s is damaged: run 'mgc alpha sync': %w", file, err)
	}
	return index, nil
}
