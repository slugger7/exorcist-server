package media

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type File struct {
	Name      string
	FileName  string
	Path      string
	Extension string
	Size      int64
}

func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", errs.BuildError(err, "error opening file")
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", errs.BuildError(err, "error calculating MD5 hash")
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

func GetFileSize(path string) (int64, error) {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return fileinfo.Size(), nil
}

func GetFilesByExtensions(root string, extensions []string) (ret []File, reterr error) {
	reterr = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			if slices.Contains(extensions, filepath.Ext(d.Name())) {
				fileSize, err := GetFileSize(d.Name())
				if err != nil {
					return err
				}
				file := File{
					Name:     GetTitleOfFile(d.Name()),
					FileName: filepath.Base(d.Name()),
					Path:     path,
					Size:     fileSize,
				}

				ret = append(ret, file)
			}
		}

		return nil
	})

	return ret, reterr
}

func FindNonExistentMedia(existingVideos []model.Media, files []File) []model.Media {
	nonExsistentVideos := []model.Media{}
	for _, v := range existingVideos {
		if !slices.ContainsFunc(files, func(mediaFile File) bool {
			return mediaFile.Path == v.Path
		}) {
			nonExsistentVideos = append(nonExsistentVideos, v)
		}
	}
	return nonExsistentVideos
}
