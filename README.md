# bun starter kit

[![build workflow](https://github.com/go-bun/bun-starter-kit/actions/workflows/build.yml/badge.svg)](https://github.com/go-bun/bun-starter-kit/actions)

Bun starter kit consists of:

- [treemux](https://github.com/vmihailenco/treemux)
- [bun](https://github.com/uptrace/bun)
- Hooks to initialize the app.
- CLI to run HTTP server and migrations, for example, `go run cmd/bun/*.go db help`.
- [example](example) package that shows how to load fixtures and test handlers.

## Quickstart

To start using this kit, clone the repo:

```shell
git clone https://github.com/go-bun/bun-starter-kit.git
```

Make sure you have correct information in `app/config/test.yaml` and then run migrations (database
must exist before running):

```shell
go run cmd/bun/main.go -env=test db init
go run cmd/bun/main.go -env=test db migrate
```

Then run the tests in [example](example) package:

```shell
cd example
go test
```

See [documentation](https://bun.uptrace.dev/guide/starter-kit.html) for more info.
