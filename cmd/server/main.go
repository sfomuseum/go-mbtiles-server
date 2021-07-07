package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-mbtiles-server/catalog"
	"github.com/sfomuseum/go-mbtiles-server/http"
	"log"
	gohttp "net/http"
	"os"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")
	root := flag.String("root", "", "The path to a directory of MBTiles (SQLite) databases. Valid databases must have one of the following suffixes: .db, .sqlite, .sqlite3, .mbtiles")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	
	flag.Parse()

	ctx := context.Background()

	db_fs := os.DirFS(*root)

	c, err := catalog.NewCatalogFromFS(db_fs)

	if err != nil {
		log.Fatalf("Failed to create reader catalog, %v", err)
	}

	handler, err := http.MBTilesHandler(c)

	if err != nil {
		log.Fatalf("Failed to create tile handler, %v", err)
	}

	mux := gohttp.NewServeMux()

	mux.Handle("/", handler)

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Failed to create new server, %v", err)
	}

	log.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to serve requests, %v", err)
	}
}
