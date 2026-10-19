package alphasync

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/argon2"
)

// The environments, named as the 'env' config names them. Tenants are not the
// same in both, so each one has its own manifests.
const (
	EnvProd    = "prod"
	EnvPreProd = "pre-prod"
)

var Envs = []string{EnvProd, EnvPreProd}

const (
	keySalt = "mgc-alpha-v2:"

	// RFC 9106, second recommended option. Changing any of them moves every
	// manifest: publish them all again under a new keySalt.
	keyTime    = 3
	keyMemory  = 64 * 1024
	keyThreads = 4
	keyLength  = 32
)

func CheckEnv(env string) error {
	if !slices.Contains(Envs, env) {
		return fmt.Errorf("unknown environment %q: use %v", env, Envs)
	}
	return nil
}

// The environment the CLI points at: the 'env' config, or prod without a
// valid one, as the login does.
func CurrentEnv(configured string) string {
	if slices.Contains(Envs, configured) {
		return configured
	}
	return EnvProd
}

// tenantDigits gives the 32 hex digits of a tenant id. Only a UUID v4 is
// accepted: its 122 random bits are what keeps the tenant out of reach of
// whoever reads the bucket.
func tenantDigits(tenant string) (string, error) {
	id, err := uuid.Parse(strings.TrimSpace(tenant))
	if err != nil || id.Version() != 4 || id.Variant() != uuid.RFC4122 {
		return "", errors.New("the tenant id is not a UUID v4")
	}
	return hex.EncodeToString(id[:]), nil
}

// TenantKey names the manifest of a tenant without revealing the tenant. The
// version digit (always 4) and the variant digit (2 fixed bits) are left out,
// so 120 random bits remain, and argon2id makes trying them slow.
func TenantKey(env, tenant string) (string, error) {
	if err := CheckEnv(env); err != nil {
		return "", err
	}
	digits, err := tenantDigits(tenant)
	if err != nil {
		return "", err
	}

	input := digits[:12] + digits[13:16] + digits[17:]
	key := argon2.IDKey([]byte(input), []byte(keySalt+env), keyTime, keyMemory, keyThreads, keyLength)
	return hex.EncodeToString(key), nil
}

func ManifestPath(env, tenant string) (string, error) {
	key, err := TenantKey(env, tenant)
	if err != nil {
		return "", err
	}
	return "alpha/tenants/" + key + ".json", nil
}
