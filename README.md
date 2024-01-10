[![Latest Release](https://img.shields.io/github/v/release/webhookdb/webhookdb-cli?color=blue&sort=semver)](https://github.com/webhookdb/webhookdb-cli/releases/latest)
[![Docker](https://github.com/webhookdb/webhookdb-cli/actions/workflows/deploy-dockerhub.yml/badge.svg)](https://hub.docker.com/r/webhookdb/webhookdb-cli/tags)
[![Tests](https://github.com/webhookdb/webhookdb-cli/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/webhookdb/webhookdb-cli/actions/workflows/pr-checks.yml)
[![Release](https://github.com/webhookdb/webhookdb-cli/actions/workflows/release.yml/badge.svg)](https://github.com/webhookdb/webhookdb-cli/actions/workflows/release.yml)
[![Docs](https://img.shields.io/badge/docs-purple)](https://docs.webhookdb.com/docs/cli-reference.html)

# webhookdb-cli

_Don't want to install this yourself?_
_Run [the WebhookDB CLI from the browser](https://webhookdb.com/terminal/)
or at the `/terminal` route of your self-hosted WebhookDB instance._

_Self-hosting? Visit `https://<webhookdb host>/terminal` to run against your own WebhookDB server._

Command Line Interface for WebhookDB ([https://github.com/webhookdb/webhookdb](https://github.com/webhookdb/webhookdb)).
WebhookDB replicates any API into a database,
so you have immediate, reliable access to all your data.

## Installing

Use the docker container (note you need to mount `/root` to persist auth info between invocations):

```
$ docker run -it webhookdb/webhookdb-cli:latest version
0.13.0 (acb64d0f)
$ docker run -v ~/.webhookdb/dockercli:/root -it webhookdb/webhookdb-cli:latest auth login
```

On MacOS, install from Homebrew:

```
$ brew install webhookdb/webhookdb-cli/webhookdb
$ webhookdb version
0.13.0 (acb64d0f)
```

On Linux, grab the binary from the latest release (package managers coming soon):

- Download the latest Linux `tar.gz` file from <https://github.com/webhookdb/webhookdb-cli/releases/latest>
- Unzip the file: `tar -xvf webhookdb_X.X.X_linux_x86_64.tar.gz`
- Move `./webhookdb` to your execution path, like `/usr/local/bin`.

On Windows, grab the executable from the zip file:

- Download the latest Windows `zip` file from <https://github.com/webhookdb/webhookdb-cli/releases/latest>
- Unzip the `webhookdb_X.X.X_windows_x86_64.zip` file.
- Run the unzipped `webhookdb.exe` file.

## Usage

To create an account and get started, run:

	webhookdb auth login

The CLI will guide you from there.

You also have quick access to the WebhookDB documentation:

- `webhookdb docs guide` to go to <https://docs.webhookdb.com/docs/getting-started/>.
- `webhookdb docs manual` or `webhookdb docs html` to go
  to [the CLI Reference](https://docs.webhookdb.com/docs/cli-reference.html).
- `webhookdb docs tui` to render [MANUAL.md](https://github.com/webhookdb/webhookdb-cli/blob/main/MANUAL.md)
  in your terminal.

## Privacy and Telemetry

The CLI collects information when unhandled exceptions are raised.
Set `WEBHOOKDB_PRIVACY` to any non-empty value to opt out of this
and any other telemetry we may add in the future.

## Releasing

Releases are automated. See `.github/workflows/release.yml`.
A new release is automatically drafted when a tag is (manually) pushed;
when the release is committed, a Dockerhub build is triggered.

There is some additional work for releasing via Homebrew and the web terminal.

The process for releasing is:

- Go to [webhookdb/homebrew-webhookdb-cli](https://github.com/webhookdb/homebrew-webhookdb-cli)
  and make sure there is an empty `next` branch.
  You can use `make create-fresh-next-branch` from the `homebrew-webhookdb-cli` repo for this.
- Tag a commit, ie `git tag 0.9.2`
- Push the tag, ie `git push origin 0.9.2`
- When it finishes, a Draft release will be built.
- A commit will also have been added to the homebrew repo's `next` branch.
- Edit the GitHub release, and publish it.
    - This deploys to Dockerhub (
      see [deploy-dockerhub.yml](https://github.com/webhookdb/webhookdb-cli/blob/main/.github/workflows/deploy-dockerhub.yml)).
- Merge the changes from `homebrew-webhookdb-cli` into `main`.
    - We cannot have goreleaser automatically push to the tap's `main`
      because it would refer to the draft release in the active formula.
      So we have to make the formula change active once the release is published.
- Update the code
  [used to serve the terminal](https://github.com/webhookdb/webhookdb/blob/main/lib/webterm/static/index.html#L54)
  to refer to the new version. Then deploy the change.

## Feedback

Please send us an email to [hello@webhookdb.com](mailto:hello@webhookdb.com)
or open an issue in one of the [webhookdb repositories](https://github.com/webhookdb).

## License

Licensed under the [Apache License 2.0 license](https://github.com/webhookdb/webhookdb-cli/blob/main/LICENSE).

Copyright (c) Lithic Technology LLC. All rights reserved.
