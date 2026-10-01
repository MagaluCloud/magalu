# Preview commands: `mgc beta` and `mgc alpha`

> **Status:** proposal, describing release 1. `mgc beta` (specs built into the
> CLI) and loading alpha specs from `~/.config/mgc/.preview/alpha` are built
> (branch `feat/preview-commands`); `mgc alpha download` is not built yet. See
> [design.md](design.md).

Preview commands let you try features before they ship in the official CLI. They
live under their own command, so the official commands never change because of
them.

| Level | Who can use it | How you get it |
| --- | --- | --- |
| `mgc beta` | Any customer | It ships with `mgc`: update `mgc` to get new beta features |
| `mgc alpha` | Customers invited by Magalu | An ID shared with you by Magalu |

A feature usually moves `mgc alpha <feature>` → `mgc beta <feature>` → `mgc <feature>`.
Once it reaches the official CLI, use the official command.

## Beta

`mgc beta --help` lists the products that have something in beta. A beta product
has every official command of that product, with the same names, plus what is in
beta:

```bash
mgc load-balancer network-loadbalancers list        # official
mgc beta load-balancer network-loadbalancers list   # same command, in beta
mgc beta load-balancer application list             # only in beta
```

So moving a script to beta, or back, only means adding or removing `beta`.

## Alpha

Magalu sends you an ID for the feature you were invited to test:

```bash
mgc alpha download 3f6c1d2e-8a4b-4c1f-9e2d-7b5a0c9d1e42
mgc alpha ai models list
mgc alpha remove ai
```

`mgc alpha` does not show up in `mgc --help` until you download something, but
`mgc alpha download` always works. Downloading the same ID again updates it.

Alpha features are kept in `.preview/alpha` in the MGC config dir
(`~/.config/mgc`, or `$MGC_CONFIG_DIR` when set). Updating `mgc` does not touch
that folder, so your alpha features survive CLI upgrades.

## Good to know

- **Preview commands may change or go away** without notice.
- **Access is still checked by the API.** Having a preview command does not grant
  access to a feature; your organization must have it enabled.
- **Alpha specs are checked before use.** An empty or invalid spec is refused, both
  on download and when you run its commands. A later release will also require a
  Magalu signature.
- **A broken alpha feature does not break the CLI.** `mgc alpha <feature>` shows
  the error and every other command keeps working.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| A beta feature you read about is not in `mgc beta --help` | Your `mgc` is older than the release that brought it. Update `mgc`. |
| `mgc alpha <feature>` says the spec is empty or invalid | The file was damaged or edited. Download it again. |
| `mgc alpha <feature>` says it needs a newer `mgc` | The feature uses something your CLI does not have yet. Update `mgc`. |
| A beta command disappeared after an update | It was promoted to the official CLI (drop `beta`) or retired. |
