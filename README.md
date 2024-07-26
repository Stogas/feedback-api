
### Local development

Run PostgreSQL with:
```shell
cd local-dev
docker compose up -d
```

PostgreSQL will be reachable at:
```
localhost:5432
user: test
password: test
database: test
```

Run app with:
```shell
go run main.go
```

Before commiting, make sure your code complies with `gofmt`:
```shell
gofmt -d . # check what's wrong
gofmt -w . # fix it
```
TIP: automatic formatting is available in the vscode Go extension. Make sure to enable vscode's `editor.formatOnSave` feature to utilize this. Also, make sure to use the golang.go formatter (provided by the Go extension) instead of Prettier.
