# Workspaces hold auth and runtime configuration for the MGC CLI, like tokens and log filter settings.
Users can create as many workspace as they choose to. Auth and config operations will affect only the
current workspace, so users can alter and switch between workspace without loosing the previous configuration

## Usage:
```bash
Usage:
  ./mgc profile workspaces [flags]
  ./mgc profile workspaces [command]
```

## Product catalog:
- Commands:
- create      Creates a new workspace
- delete      Deletes the workspace with the specified name
- get         Get current workspace.
- list        List all available workspaces
- set         Sets workspace to be used

## Other commands:
- Additional Commands:
- select      call "list", prompt selection and then "set"

## Flags:
```bash
Flags:
  -h, --help   help for workspaces
```

