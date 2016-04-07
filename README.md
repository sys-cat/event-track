# How to Use this API

routing is root only.

## request example

```
curl -v -X POST -d '{"rid":"rid-test","event":100,"referer":"http://localhost:8080","environment":"staging"}' http://url:8080/
```

## request paramaters

| name        | type    | null | example                   |
|-------------|---------|------|---------------------------|
| rid         | string  | no   | "rid-test"                |
| event       | integer | no   | 100                       |
| referer     | string  | yes  | "http://test.com"         |
| environment | string  | no   | "staging" or "production" |
