# Register new SSH key by providing a name and the public SSH key

## Usage:
```bash
The supported key types are: ssh-rsa, ssh-dss, ecdsa-sha, ssh-ed25519, sk-ecdsa-sha, sk-ssh-ed25519
```

## Product catalog:
- Usage:
- mgc profile ssh-keys create [flags]

## Other commands:
- Flags:
- -h, --help          help for create
- --key string    The SSH public key. The supported key types are: ssh-rsa, ssh-dss, ecdsa-sha, ssh-ed25519, sk-ecdsa-sha, sk-ssh-ed25519 (max character count: 16384) (required)
- --name string   The SSH Key name (max character count: 45) (required)
- -v, --version       version for create

## Flags:
```bash
Global Flags:
      --api-key string           Use your API key to authenticate with the API
  -U, --cli.retry-until string   Retry the action with the same parameters until the given condition is met. The flag parameters
                                 use the format: 'retries,interval,condition', where 'retries' is a positive integer, 'interval' is
                                 a duration (ex: 2s) and 'condition' is a 'engine=value' pair such as "jsonpath=expression"
  -t, --cli.timeout duration     If > 0, it's the timeout for the action execution. It's specified as numbers and unit suffix.
                                 Valid unit suffixes: ns, us, ms, s, m and h. Examples: 300ms, 1m30s
      --debug                    Display detailed log information at the debug level
      --env enum                 Environment to use (one of "pre-prod" or "prod") (default "prod")
      --no-confirm               Bypasses confirmation step for commands that ask a confirmation from the user
  -o, --output string            Change the output format. Use '--output=help' to know more details.
  -r, --raw                      Output raw data, without any formatting or coloring
      --server-url uri           Manually specify the server to use
```

