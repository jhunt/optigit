# Optigit

A dashboard that allows you to track the outstanding pull requests
and issues on orgs/repositories of your choosing.

## Testing

Load test data into a SQLite database:

```
$ sqlite3 db/test.db < db/schema.sql && sqlite3 db/test.db < db/init.sql
```

Create and run the optigit binary:

```
$ make
$ GITHUB_TOKEN=1234 DATABASE=sqlite3:db/test.db ORGS=test ./optigit
```
