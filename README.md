# Optigit

A dashboard that allows you to track the outstanding pull requests
and issues on orgs/repositories of your choosing.

![Screenshot](screenshot.png)

## Testing

```
$ make
$ GITHUB_TOKEN=1234 DATABASE=sqlite3:test.db ORGS=test ./optigit
```

## Configuration

Optigit is configured entirely via environment variables.

- **ORGS** (required) A space-separated list of Github org names
  to watch.

- **DATABASE** (required) The full DSN of the database to use for
  storing state.  If this is not set, the `$VCAP_SERVICES`
  environment variable (for Cloud Foundry apps) will also be
  consulted.  One of these _MUST_ be provided.

- **GITHUB_TOKEN** (required) The Github Access Token to use when
  communicating with the Github API.

- **BIND** (optional) The IP and port (`ip:port`) to bind to and
  listen for incoming HTTP web requests.  Takes precedence over
  `$PORT`

- **PORT** (optional) The port to bind to and listen for incoming
  HTTP web requests.  All interfaces will be bound on the given
  port, if specified.  Defaults to `3000`

- **REFERSH_INTERVAL** (optional) How often to refresh data from
  Github.  Accepts units in minute (35m, 15m, etc.), hours (4h,
  1h, etc.), and days (1d, 7d, etc.).  Defaults to `1d`

- **IGNORE** (optional) A space-separated list of usernames to
  ignore issues and PRs from, by default.  This can be useful if
  you want to ignore your own bookkeeping issues, or if you want
  to focus on external OSS contributors at the expense of your
  in-house team.  Defaults to an empty list.
