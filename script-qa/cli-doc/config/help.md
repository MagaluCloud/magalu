# Configuration values are available to be set so that they persist between
different executions of the MgcSDK. They reside in a YAML file when set.
Config values may also be loaded via Environment Variables. Any Config available
(see 'list') may be exported as an env variable in uppercase with the 'MGC_' prefix

## Usage:
```bash
Usage:
  ./cli config [flags]
  ./cli config [command]
```

## Product catalog:
- Commands:
- delete      Delete/unset a Config value that had been previously set
- get         Get a specific Config value that has been previously set
- get-schema  Get the JSON Schema for the specified Config
- list        List all available Configs
- set         Set a specific Config value in the configuration file

## Other commands:
- Flags:
- -h, --help   help for config

## Flags:
```bash
Global Flags:
      --cli.show-cli-globals   Show all CLI global flags on usage text
```

