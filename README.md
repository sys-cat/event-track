# How to Use this API

```
$ go run main.go
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] POST   /                         --> main.addRecord (3 handlers)
[GIN-debug] POST   /event/add                --> main.addEvent (3 handlers)
[GIN-debug] POST   /event/report             --> main.getReport (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

## root access

| url | method |
|---|---|
| / | POST |

### request example

```
curl -v -X POST -d '{"rid":"rid-test","event":100,"referer":"http://localhost:8080","environment":"staging"}' http://url:8080/
```

### request paramaters

| name        | type    | null | example                   |
|-------------|---------|------|---------------------------|
| rid         | string  | no   | "rid-test"                |
| event       | integer | no   | 100                       |
| referer     | string  | yes  | "http://test.com"         |
| environment | string  | no   | "staging" or "production" |

## Add event

| url | method |
|---|---|
| /event/add | POST |

### request example

```
curl -v -X POST -d '{"id":100, "name":"sample-event"}' http://url:8080/event/add
```

### request paramaters

| name        | type    | null | example                   |
|-------------|---------|------|---------------------------|
| id         | int  | no   | 100                |
| name       | string | no   | "sample-event"                       |

## Get Daily Event Report

| url | method |
|---|---|
| /event/report | POST |

### request example

```
curl -v -X POST -d '{"id":100, "environment":"staging"}' http://url:8080/event/report
```

### request paramaters

| name        | type    | null | example                   |
|-------------|---------|------|---------------------------|
| id         | int  | no   | 100                |
| environment | string  | no   | "staging" or "production" |
