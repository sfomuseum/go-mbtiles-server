package tile

import (
	"github.com/tilezen/go-tilepacks/tilepack"
)

type TileRequest struct {
	Tile        *tilepack.Tile
	Layer       string
	ContentType string
}
