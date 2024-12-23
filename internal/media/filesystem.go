package media

import (
	"io/fs"
	"path/filepath"
	"slices"
)

func GetFilesByExtensions(root string, extensions []string) (ret []string, reterr error) {
	reterr = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if slices.Contains(extensions, filepath.Ext(d.Name())) {
				ret = append(ret, filepath.Join(path, d.Name()))
			}
		}

		return nil
	})

	return ret, nil
}
