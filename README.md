
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