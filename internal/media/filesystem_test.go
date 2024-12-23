package media_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/slugger7/exorcist/internal/media"
)

func convertDirEntriesToString(dirs []os.DirEntry) (ret []string) {
	for _, d := range dirs {
		ret = append(ret, d.Name())
	}

	return ret
}

func compareStringArrays(t *testing.T, got, want []string) {
	if len(got) != len(want) {
		t.Errorf("Array lengths differed: got %v but wanted %v", len(got), len(want))
	} else {
		for i, v := range got {
			if v != want[i] {
				t.Errorf("Value in array was different at index %v: got '%v' but wanted '%v'", i, v, want[i])
			}
		}
	}
}

func printList(list []string) {
	for _, l := range list {
		fmt.Println(l)
	}
}

func TestGetFilesByExtensions(t *testing.T) {
	got, _ := GetFilesByExtensions("./test_data", []string{".toml"})

	want := []string{
		"test_data/folder_1/folder_1_file.toml/folder_1_file.toml",
		"test_data/folder_2/subFolder2/subfile.toml/subfile.toml",
	}

	compareStringArrays(t, got, want)
}
