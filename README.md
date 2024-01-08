# webhookdb-cli

_Don't want to install this yourself?_
_Run [the WebhookDB CLI from the browser](https://webhookdb.com/terminal/)
or at the `/terminal` route of your self-hosted WebhookDB instance._

Command Line Interface for WebhookDB ([https://github.com/webhookdb/webhookdb](https://github.com/webhookdb/webhookdb)).
WebhookDB replicates any API into a database,
so you have immediate, reliable access to all your data.

To create an account and get started, run:

	webhookdb auth login

The CLI will guide you from there.

You also have quick access to the WebhookDB documentation:

	webhookdb docs html
	webhookdb docs tui

Use `webhookdb docs html` or
visit [https://webhookdb.com/docs/cli](https://webhookdb.com/docs/cli) for usage instructions.

Or check out [MANUAL.md](https://github.com/webhookdb/webhookdb-cli/blob/main/MANUAL.md)
to see every command and option.

## Privacy and Telemetry

The CLI collects information when unhandled exceptions are raised.
Set `WEBHOOKDB_PRIVACY` to any non-empty value to opt out of this
and any other telemetry we may add in the future.

## Releasing

Releases are automated. See `.github/workflows/release.yml`.
There is some additional work for releasing via Homebrew and the web terminal.

The process for releasing is:

- Go to [webhookdb/homebrew-webhookdb](https://github.com/webhookdb/homebrew-webhookdb)
  and make sure there is an empty `next` branch.
  You can use `make create-fresh-next-branch` from the `homebrew-webhookdb` repo for this.
- Tag a commit, ie `git tag 0.9.2`
- Push the tag, ie `git push origin 0.9.2`
- When it finishes, a Draft release will be built.
- A commit will also have been added to the homebrew repo's `next` branch.
- Edit the GitHub release, and publish it.
- Merge the changes from `homebrew-webhookdb` into `main`.
    - We cannot have goreleaser automatically push to `main`
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
