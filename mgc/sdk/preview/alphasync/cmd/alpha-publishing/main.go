package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/MagaluCloud/magalu/mgc/sdk/preview"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/alphasync"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/serverpin"
	"github.com/MagaluCloud/magalu/mgc/sdk/preview/specvalidator"
	"github.com/Masterminds/semver/v3"
	"github.com/google/uuid"
)

var validKeyID = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return errors.New("usage: alpha-publishing keygen|sign|invite|grant|revoke [flags]")
	}

	switch args[0] {
	case "keygen":
		return keygen(args[1:], stdout)
	case "sign":
		return sign(args[1:], stdout)
	case "invite":
		return invite(args[1:], stdout)
	case "grant":
		return grant(args[1:], stdout)
	case "revoke":
		return revoke(args[1:], stdout)
	default:
		return fmt.Errorf("unknown command %q: use keygen, sign, invite, grant or revoke", args[0])
	}
}

func keygen(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("keygen", flag.ContinueOnError)
	keysFile := flags.String("keys", "", "trusted keys file the CLI embeds")
	keyID := flags.String("key-id", "", "id of the new key, such as 2026-10")
	out := flags.String("out", "", "new file for the private key")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *keysFile == "" || *out == "" || !validKeyID.MatchString(*keyID) {
		return errors.New("keygen needs --keys, --out and a --key-id made of letters, digits, '.', '-' and '_'")
	}

	trusted, err := readKeysFile(*keysFile)
	if err != nil {
		return err
	}
	if _, ok := trusted[*keyID]; ok {
		return fmt.Errorf("key id %q is already trusted: pick a new one", *keyID)
	}

	public, private, err := ed25519.GenerateKey(nil)
	if err != nil {
		return err
	}

	// O_EXCL: never overwrite a private key that may be the only copy.
	file, err := os.OpenFile(*out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintln(file, hex.EncodeToString(private.Seed())); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	trusted[*keyID] = hex.EncodeToString(public)
	if err := writeKeysFile(*keysFile, trusted); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Public key %q added to %s: commit it, and release mgc with it before publishing.\n", *keyID, *keysFile)
	fmt.Fprintf(stdout, "Private key written to %s: store it in the vault, then delete the file.\n", *out)
	return nil
}

func sign(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("sign", flag.ContinueOnError)
	keysFile := flags.String("keys", "", "trusted keys file the CLI embeds")
	keyFile := flags.String("key", "", "private key file")
	module := flags.String("module", "", "command under 'mgc alpha', such as iam")
	specFile := flags.String("spec", "", "compiled spec")
	minCLIVersion := flags.String("min-cli-version", "", "oldest mgc that runs the spec, when it needs one")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if *keysFile == "" || *keyFile == "" || *module == "" || *specFile == "" {
		return errors.New("sign needs --keys, --key, --module and --spec")
	}
	if err := alphasync.CheckModule(*module); err != nil {
		return err
	}
	if *minCLIVersion != "" {
		if _, err := semver.NewVersion(*minCLIVersion); err != nil {
			return fmt.Errorf("--min-cli-version %q is not a version", *minCLIVersion)
		}
	}

	private, err := readPrivateKey(*keyFile)
	if err != nil {
		return err
	}
	encoded, err := os.ReadFile(*keysFile)
	if err != nil {
		return err
	}
	trusted, err := alphasync.ParseKeys(encoded)
	if err != nil {
		return err
	}
	keyID, err := trusted.KeyIDOf(private)
	if err != nil {
		return err
	}

	spec, err := os.ReadFile(*specFile)
	if err != nil {
		return err
	}
	if err := checkSpec(*module, spec); err != nil {
		return fmt.Errorf("the CLI would refuse this spec: %w", err)
	}

	sum := sha256.Sum256(spec)
	item := alphasync.Sign(private, keyID, alphasync.ManifestSpec{
		Module:        *module,
		Path:          "alpha/specs/" + *module + "/" + uuid.NewString() + ".openapi.yaml",
		SHA256:        hex.EncodeToString(sum[:]),
		MinCLIVersion: *minCLIVersion,
	})
	if err := trusted.Verify(item); err != nil {
		return err
	}

	encoder := json.NewEncoder(stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(item)
}

func invite(args []string, stdout io.Writer) error {
	flags := flag.NewFlagSet("invite", flag.ContinueOnError)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 0 {
		return errors.New("invite takes no arguments")
	}
	_, err := fmt.Fprintln(stdout, alphasync.NewInvite())
	return err
}

func checkSpec(module string, spec []byte) error {
	pinned, err := serverpin.Pin(spec)
	if err != nil {
		return err
	}
	extensionPrefix := "x-mgc"
	return preview.Check(preview.Alpha, preview.Package{Name: module, Spec: pinned}, specvalidator.New(), &extensionPrefix)
}

func readPrivateKey(file string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	seed, err := hex.DecodeString(strings.Join(strings.Fields(string(data)), ""))
	if err != nil || len(seed) != ed25519.SeedSize {
		return nil, fmt.Errorf("%s does not hold a private key: expected %d bytes in hex", file, ed25519.SeedSize)
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

func readKeysFile(file string) (map[string]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if _, err := alphasync.ParseKeys(data); err != nil {
		return nil, err
	}
	var keys map[string]string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func writeKeysFile(file string, keys map[string]string) error {
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, append(data, '\n'), 0o644)
}
