MagaLu CLI
==========

The Magalu Cloud command line interface (CLI), enables users to
authenticate, configure and control the cloud resource

## Goals

- Allow the user to authenticate itself
- Save the configuration
- Generate commands at runtime (No need for new binaries) by means of OpenAPI schema files

## Telemetry

The CLI collects pseudonymous usage data to help prioritize improvements. After
each command it sends one event about how the command ran, never flag values,
credentials, error messages or command output. Help, `--version`, `completion`
and `mgc telemetry` itself are not collected. The full list of collected fields
is in [docs/telemetry](./docs/telemetry/help.md).

A notice is shown once, in an interactive terminal, before collection starts.
The event is sent in the background by a separate process, so the command never
waits for the network. Network errors and server errors are retried a few times,
and the event is silently dropped if it still cannot be sent. With `--debug`,
the sending process writes its diagnostics to `telemetry-debug.txt` in the CLI
config directory.

To opt out, use any of:

- `mgc telemetry disable` (saved in `telemetry.yaml` in the CLI config directory,
  valid for every profile; `mgc telemetry enable` turns it back on)
- `MGC_CLI_TELEMETRY_OPTOUT=1`
- `DO_NOT_TRACK=1`

`mgc telemetry status` shows whether collection is enabled and what disabled it.
See the [privacy policy](https://magalu.cloud/termos-legais/politica-de-privacidade/).

## Development

See [DEVELOPMENT.md](./DEVELOPMENT.md)

## OpenAPI

See [sdk/openapi/README.md](../sdk/openapi/README.md)
