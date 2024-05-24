# Profiles hold auth and runtime configuration for the MgcSDK, like tokens and log filter settings.
Users can create as many profiles as they choose to. Auth and config operations will affect only the
current profile, so users can alter and switch between profiles without loosing the previous configuration

## Usage:
```bash
Usage:
  ./mgc profile [flags]
  ./mgc profile [command]
```

## Product catalog:
- Commands:
- create         Creates a new profile
- current        Shows current selected profile. Any changes to auth or config values will only affect this profile
- delete         Deletes the profile with the specified name
- list           List all available profiles
- set-current    Sets profile to be used

## Other commands:
- Additional Commands:
- select-current call "list", prompt selection and then "set-current"

## Flags:
```bash
Flags:
  -h, --help   help for profile
```

