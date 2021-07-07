# go-mbtiles-server

Go package for serving a collection of MBTiles (SQLite) databases via an HTTP endpoint.

## Important

This is work in progress. Documentation is incomplete.

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/sfomuseum/go-mbtiles-server.svg)](https://pkg.go.dev/github.com/sfomuseum/go-mbtiles-server)

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
> ./bin/server -h
Usage:
  ./bin/server [options]

Options:
  -root string
    	The path to a directory of MBTiles (SQLite) databases. Valid databases must have one of the following suffixes: .db, .sqlite, .sqlite3, .mbtiles
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
```

For example:

```
$> ./bin/server  -root /usr/local/tiles/sqlite -server-uri http://localhost:8181
2021/07/07 11:17:20 Listening on http://localhost:8181
```

And then:

```
$> curl -I http://localhost:8181/aerial/2019/16/10493/40167.png
HTTP/1.1 200 OK
Content-Length: 131970
Content-Type: image/png
Date: Wed, 07 Jul 2021 18:19:28 GMT
```

