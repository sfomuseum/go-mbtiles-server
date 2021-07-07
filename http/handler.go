package http

import (
	"github.com/sfomuseum/go-mbtiles-server/tile"
	"github.com/tilezen/go-tilepacks/tilepack"
	"github.com/aaronland/go-mimetypes"
	"log"
	gohttp "net/http"
	"strings"
)

func MBTilesHandler(catalog map[string]tilepack.MbtilesReader) (gohttp.HandlerFunc, error) {

	p, err := tile.NewSimpleTileParser()

	if err != nil {
		return nil, err
	}

	return MBTilesHandlerWithParser(catalog, p)
}

func MBTilesHandlerWithParser(catalog map[string]tilepack.MbtilesReader, p tile.TileParser) (gohttp.HandlerFunc, error) {

	fn := func(w gohttp.ResponseWriter, r *gohttp.Request) {

		tile_req, err := p.Parse(r.URL.Path)

		if err != nil {
			gohttp.NotFound(w, r)
			return
		}

		reader, ok := catalog[tile_req.Layer]

		if !ok {
			gohttp.NotFound(w, r)
			return
		}

		result, err := reader.GetTile(tile_req.Tile)

		if err != nil {
			log.Printf("Error getting tile: %+v", err)
			gohttp.NotFound(w, r)
			return
		}

		if result.Data == nil {
			gohttp.NotFound(w, r)
			return
		}

		acceptEncoding := r.Header.Get("Accept-Encoding")

		if strings.Contains(acceptEncoding, "gzip") {
			w.Header().Set("Content-Encoding", "gzip")
		} else {
			log.Printf("Requester doesn't accept gzip but our mbtiles have gzip in them")
		}

		t := mimetypes.TypesByExtension(tile_req.ContentType)

		if len(t) >= 1 {
			w.Header().Set("Content-Type", t[0])
		}
		
		w.Write(*result.Data)
	}

	return fn, nil
}
