# Feedback-api

Feedback API is a RESTful Go API allowing submission and updates to a user's satisfaction report for other apps.

The main goal is to be simple and allow partial submittions (i.e. without further comments or issue types), if a user decides to abandon the process in the middle of the report.

Thus, the first request by the front-end must be a POST with a newly generated UUID, and subsequent requests must be of type PATCH with the same UUID - more details below.

HTTP header `X-Feedback-Submit-Token` is a very rudimentary approach to prevent random submissions - this token should be known to your frontend, and as such, should not be considered a "secret". If necessary, one can rotate this token with every frontend update/deployment.

A secondary goal is to be a generic Feedback API, i.e. allow this to be used in a variety of projects. For some needs, it might be required to allow submissions from authenticated users only. Thus, this project *might* implement optional JWT token validation instead of the well-known Submit Token later on.

For running in production, it's recommended to have a reverse proxy in front with an IP-based ratelimiter in order to partially prevent spam attacks.

## Features

- JSON logs enabled by default (set `LOGS_JSON=false` to disable)
- Rudimentary OpenTelemetry tracing and exporting via OTLP gRPC
- Prometheus metrics (exported by default on `0.0.0.0:2222/metrics`) for HTTP and satisfaction values
- PostgreSQL as database (can be modified to support [other GORM DBs](https://gorm.io/docs/connecting_to_the_database.html))
- Automatic unexpected panic recovery (via the `gin.Recovery()` middleware)
- Automatic recovery after DB downtime

## Usage

To get issue types, query this:

```
GET /issues
```

To create a new satisfaction report, submit this:

```
POST /submit/satisfaction

HTTP headers:
X-Feedback-Submit-Token: <value of API_SUBMIT_TOKEN>

payload:
{
  "satisfied": <bool>,
  "uuid": "<new client-side generated UUID>",
  "issue_id": <int>, # optional
	<...> other data
}
```

To update a satisfaction report, submit this:
```
PATCH /submit/satisfaction

HTTP headers:
X-Feedback-Submit-Token: <value of API_SUBMIT_TOKEN>

payload:
{
  "satisfied": <bool>,
  "uuid": "<existing UUID>",
  "issue_id": <int>, # optional
	<...> other data
}
```

Rules:
- Trying to POST without X-Feedback-Submit-Token will return `HTTP 401 Unauthorized`
- Trying to POST without either `.satisfied` or `.uuid` will return `HTTP 400 Bad Request`
- Trying to POST with an *existing* `.uuid` will return `HTTP 409 Conflict`
- Trying to PATCH with a *new* `.uuid` will return `HTTP 404 Not Found`

## Local development

Run PostgreSQL, Grafana Tempo & Grafana with:
```shell
cd local-dev
docker compose up -d
```

Run app with the default Environment Variables and credentials (found in `.env`):
```shell
go run main.go
```

### Linting & formatting style

Code must comply with the [configured linters](.golangci.yaml) in `golangci-lint`.

#### Linting

To integrate linting into your IDE, follow [instructions here](https://golangci-lint.run/welcome/integrations/#editor-integration).

For VSCode users, install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go), and configure:

  - `go.lintTool` with value `golangci-lint`
  - `go.lintFlags` with value `--fast`

VSCode will promt you to install `golangci-lint` automatically.

For non-IDE users, `golangci-lint` can be installed with [these instructions](https://golangci-lint.run/welcome/install/#local-installation), and then run with `golangci-lint run`.

#### Formatting

Code must follow `gofmt` and `goimports` formatting conventions. For ease of use in IDEs, the `gopls` language server covers these rules.

For VSCode users, install the [Go extension](https://marketplace.visualstudio.com/items?itemName=golang.Go), select `default` in `go.formatTool`, which will use `gopls`. Make sure `go.useLanguageServer` and `editor.formatOnSave` are enabled (they are by default).

For non-IDE users, install `goimports` and run manually:
```shell
go install golang.org/x/tools/cmd/goimports@latest # install
goimports -d . # check what's wrong
goimports -w . # fix it
```
