package main

import (
	"flag"
	"io/fs"
	"net/http"
	"fmt"
	"log"
	"strings"
	"strconv"
	"github.com/tilezen/go-tilepacks/tilepack"
	"regexp"
	"path/filepath"
)

type TileRequest struct {
	Tile *tilepack.Tile
	Layer string
	ContentType string
}

type TileParser interface {
	Parse(string) (*TileRequest, error)
}

type SimpleTileParser struct {
	TileParser
	re *regexp.Regexp
}

func NewSimpleTileParser() (TileParser, error) {
	
	re, err := regexp.Compile(`\/([^\/]+)\/(\d+)\/(\d+)\/(\d+)\.(\w+)$`)

	if err != nil {
		return nil, err
	}

	p := &SimpleTileParser{
		re: re,
	}

	return p, nil
}

func (p *SimpleTileParser) Parse(path string) (*TileRequest, error) {

	match := p.re.FindStringSubmatch(path)
	
	if match == nil {
		return nil, fmt.Errorf("invalid tile path")
	}

	layer := match[1]
	// ext := match[5]
	
	z, _ := strconv.ParseUint(match[2], 10, 32)
	x, _ := strconv.ParseUint(match[3], 10, 32)
	y, _ := strconv.ParseUint(match[4], 10, 32)

	tile := &tilepack.Tile{
		Z: uint(z),
		X: uint(x),
		Y: uint(y),
	}
	
	tile_req := &TileRequest{
		Tile: tile,
		Layer: layer,
		ContentType: "",	// Fix me
	}
		
	return tile_req, nil		
}

func MBTilesHandler(catalog map[string]tilepack.MbtilesReader) (http.HandlerFunc, error) {

	p, err := NewSimpleTileParser()

	if err != nil {
		return nil, err
	}

	return MBTilesHandlerWithParser(catalog, p)
}

func MBTilesHandlerWithParser(catalog map[string]tilepack.MbtilesReader, p TileParser) (http.HandlerFunc, error) {	
		
	fn := func(w http.ResponseWriter, r *http.Request) {
		
		tile_req, err := p.Parse(r.URL.Path)
		
		if err != nil {
			http.NotFound(w, r)
			return
		}

		reader, ok := catalog[tile_req.Layer]

		if !ok {
			http.NotFound(w, r)
			return
		}
		
		result, err := reader.GetTile(tile_req.Tile)
		
		if err != nil {
			log.Printf("Error getting tile: %+v", err)
			http.NotFound(w, r)
			return
		}

		if result.Data == nil {
			http.NotFound(w, r)
			return
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")
		
		if strings.Contains(acceptEncoding, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
		} else {
			log.Printf("Requester doesn't accept gzip but our mbtiles have gzip in them")
		}

		// FIX ME
		// w.Header().Set("Content-Type", tile_req.ContentType)
		w.Write(*result.Data)
	}

	return fn, nil
}

func ReaderCatalogFromFS(db_fs fs.FS) (map[string]tilepack.MbtilesReader, error){
	return ReaderCatalogFromFSWithPattern(db_fs, "*.db")
}

func ReaderCatalogFromFSWithPattern(db_fs fs.FS, pat string) (map[string]tilepack.MbtilesReader, error){	

	catalog := make(map[string]tilepack.MbtilesReader)
	
	files, err := fs.Glob(db_fs, pat)

	if err != nil {
		return nil, fmt.Errorf("Failed to glob FS, %v", err)
	}

	for _, path := range files {

		fname := filepath.Base(path)
		ext := filepath.Ext(fname)

		layer := strings.Replace(fname, ext, "", -1)

		r, err := tilepack.NewMbtilesReader(path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new MBTiles reader for '%s', %v", path, err)
		}

		catalog[layer] = r
	}

	return catalog, nil
}

func main() {

	flag.Parse()

	var db_fs fs.FS

	catalog, err := ReaderCatalogFromFS(db_fs)

	if err != nil {
		log.Fatalf("Failed to create reader catalog, %v", err)
	}

	handler, err := MBTilesHandler(catalog)

	if err != nil {
		log.Fatalf("Failed to create tile handler, %v", err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", handler)

	addr := "localhost:8080"
	err = http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatalf("Failed to serve requests, %v", err)
	}
}
