package tile

import (
	"fmt"
	"github.com/tilezen/go-tilepacks/tilepack"
	"regexp"
	"strconv"
)

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
		Tile:        tile,
		Layer:       layer,
		ContentType: "", // Fix me
	}

	return tile_req, nil
}
