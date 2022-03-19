# webhookdb-cli

_Don't want to install this yourself?_
_Run [the WebhookDB CLI from the browser](https://webhookdb.com/terminal/)._

Command Line Interface for WebhookDB ([https://webhookdb.com](https://webhookdb.com)).
WebhookDB allows you to query any API in real-time with SQL.

To create an account and get started, run:

	webhookdb auth login

The CLI will guide you from there.

You also have quick access to the WebhookDB documentation:

	webhookdb docs html
	webhookdb docs tui

Use `webhookdb docs html` or
visit [https://webhookdb.com/docs/cli](https://webhookdb.com/docs/cli) for usage instructions.

Or check out [MANUAL.md](https://github.com/lithictech/webhookdb-cli/blob/main/MANUAL.md)
to see every command and option.

## Privacy and Telemetry

The CLI collects information when unhandled exceptions are raised.
Set `WEBHOOKDB_PRIVACY` to any non-empty value to opt out of this
and any other telemetry we may add in the future.

## Releasing

Releases are automated. See `.github/workflows/release.yml`.
There is some additional work for releasing via Homebrew and the web terminal.

The process for releasing is:

- Go to [lithictech/homebrew-webhookdb](https://github.com/lithictech/homebrew-webhookdb)
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
- Update the code [in the WebhookDB website](https://github.com/lithictech/webhookdb-api/blob/main/webhookdb-website/static/webterm/index.html#L33)
  to refer to the new version. Then merge the change, the website deploys from `main`.

## Feedback

Please send us an email, [webhookdb@lithic.tech](mailto:webhookdb@lithic.tech)
Got feedback for us? Please don't hesitate to tell us on feedback.

## License

Copyright (c) Lithic Technology LLC. All rights reserved.

Licensed under the [Apache License 2.0 license](https://github.com/lithictech/webhookdb-cli/blob/main/LICENSE).
