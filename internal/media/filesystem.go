package media

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type File struct {
	Name      string
	FileName  string
	Path      string
	Extension string
}

func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error calculating MD5 hash: %w", err)
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func GetRelativePath(root, path string) string {
	return strings.Replace(path, root, "", 1)
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
					Name:     GetTitleOfFile(d.Name()),
					FileName: filepath.Base(d.Name()),
					Path:     path,
				}

				ret = append(ret, file)
			}
		}

		return nil
	})

	return ret, reterr
}
