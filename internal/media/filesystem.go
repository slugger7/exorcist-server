package media

import (
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
)

type File struct {
	Name      string
	Path      string
	Extension string
}

func GetTitleOfFile(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) == 1 {
		return filename
	}

	parts = parts[:len(parts)-1]

	return strings.Join(parts, ".")
}

func GetFilesByExtensions(root string, extensions []string) (ret []File, reterr error) {
	reterr = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if slices.Contains(extensions, filepath.Ext(d.Name())) {
				file := File{
					Name: GetTitleOfFile(d.Name()),
					Path: path,
				}

				ret = append(ret, file)
			}
		}

		return nil
	})

	return ret, reterr
}
