package catalog

import (
	"fmt"
	"github.com/tilezen/go-tilepacks/tilepack"
	"io/fs"
	"path/filepath"
	"strings"
	"os"
	"regexp"
)

// MatchFunc is a custom function to signal whether a given path is a valid MBTiles database path or filename.
type MatchFunc func(string) (bool, error)

// NewCatalogFromFS return a lookup table of tilepack.MbtilesReader instances from a fs.FS instance. This method
// will match filenames (in 'db_fs') ending in `.*\.(sqlite|sqlite3|mbtiles|db)$`
func NewCatalogFromFS(db_fs fs.FS) (map[string]tilepack.MbtilesReader, error) {

	re, err := regexp.Compile(`.*\.(sqlite|sqlite3|mbtiles|db)$`)

	if err != nil {
		return nil, err
	}
	
	return NewCatalogFromFSWithRegexp(db_fs, re)
}

// NewCatalogFromFS return a lookup table of tilepack.MbtilesReader instances from a fs.FS instance whose
// children match the regular expression defined in 're'.
func NewCatalogFromFSWithRegexp(db_fs fs.FS, re *regexp.Regexp) (map[string]tilepack.MbtilesReader, error) {

	match_func := func(path string) (bool, error) {

		if re.MatchString(path){
			return true, nil
		}

		return false, nil
	}

	return NewCatalogFromFSWithMatchFunc(db_fs, match_func)
}

// NewCatalogFromFS return a lookup table of tilepack.MbtilesReader instances from a fs.FS instance whose
// children return true when compared against the MatchFunction defined in 'match_func'
func NewCatalogFromFSWithMatchFunc(db_fs fs.FS, match_func MatchFunc) (map[string]tilepack.MbtilesReader, error) {

	files := make([]string, 0)
	
	walk_fn := func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			return err
		}

		ok, err := match_func(path)

		if err != nil {
			return err
		}

		if ok {
			files = append(files, path)
		}

		return nil
	}

	err := fs.WalkDir(db_fs, ".", walk_fn)

	if err != nil {
		return nil, err
	}

	catalog := make(map[string]tilepack.MbtilesReader)
	
	for _, path := range files {

		fname := filepath.Base(path)
		ext := filepath.Ext(fname)

		layer := strings.Replace(fname, ext, "", -1)

		// START OF things we probably shouldn't be doing
		
		var abs_path string

		str_type := fmt.Sprintf("%T", db_fs)
		
		switch str_type {
		case "os.dirFS":
			root_fs := fmt.Sprintf("%s", db_fs)	// SEE THIS? It totally breaks the fs.FS contract
			abs_path = filepath.Join(root_fs, path)
		default:

			// Create tmp files and io.Copy the data in to place
			// Note the absence of any good way to clean up that tmp data
			// yet...
			
			return nil, fmt.Errorf("Unsupport fs.FS type '%T'", db_fs)
		}

		if abs_path == "" {
			return nil, fmt.Errorf("Invalid absolute path for '%s'", path)
		}

		_, err = os.Stat(abs_path)

		if err != nil {
			return nil, fmt.Errorf("Unable to locate '%s', %v", abs_path, err)
		}

		// END OF things we probably shouldn't be doing
		
		r, err := tilepack.NewMbtilesReader(abs_path)

		if err != nil {
			return nil, fmt.Errorf("Failed to create new MBTiles reader for '%s', %v", path, err)
		}

		catalog[layer] = r
	}

	return catalog, nil
}
