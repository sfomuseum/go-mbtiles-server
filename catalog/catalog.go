package catalog

import (
	"fmt"
	"github.com/tilezen/go-tilepacks/tilepack"
	"io/fs"
	"path/filepath"
	"strings"
)

func NewCatalogFromFS(db_fs fs.FS) (map[string]tilepack.MbtilesReader, error) {
	return NewCatalogFromFSWithPattern(db_fs, "*.db")
}

func NewCatalogFromFSWithPattern(db_fs fs.FS, pat string) (map[string]tilepack.MbtilesReader, error) {

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
