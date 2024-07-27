
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

## Usage

To create a new satisfaction report, submit this:

```
POST /submit/satisfaction

headers:
X-Feedback-Submit-Token: <value of API_SUBMIT_TOKEN>

payload:
{
  "satisfied": <bool>,
  "uuid": "<new client-side generated UUID>",
	<...> other data
}
```

To update a satisfaction report, submit this:
```
PATCH /submit/satisfaction

headers:
X-Feedback-Submit-Token: <value of API_SUBMIT_TOKEN>

payload:
{
  "satisfied": <bool>,
  "uuid": "<existing UUID>",
	<...> other data
}
```

Rules:
- Trying to POST without X-Feedback-Submit-Token will return `HTTP 401 Unauthorized`
- Trying to POST without either `.satisfied` or `.uuid` will return `HTTP 400 Bad Request`
- Trying to POST with an *existing* `.uuid` will return `HTTP 409 Conflict`
- Trying to PATCH with a *new* `.uuid` will return `HTTP 404 Not Found`
