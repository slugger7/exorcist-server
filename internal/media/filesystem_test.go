package media_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	. "github.com/slugger7/exorcist/internal/media"
)

func compareFileArrays(t *testing.T, got, want []File) {
	if len(got) != len(want) {
		t.Errorf("Array lengths differed: got %v but wanted %v", len(got), len(want))
	} else {
		for i, f := range got {
			if f.Path != want[i].Path {
				t.Errorf("Path in file array was different at index %v: got '%v' but wanted '%v'", i, f.Path, want[i].Path)
			}
			if f.FileName != want[i].FileName {
				t.Errorf("Filename in file array was different at index %v: got '%v' but wanted '%v'", i, f.FileName, want[i].FileName)
			}
		}
	}
}

func Test_GetFilesByExtensions(t *testing.T) {
	got, _ := GetFilesByExtensions("./test_data", []string{".toml"})

	want := []File{
		{
			Path:     "test_data/folder_1/folder_1_file.toml",
			FileName: "folder_1_file.toml",
			Name:     "",
		},
		{
			Path:     "test_data/folder_2/subFolder2/subfile.toml",
			FileName: "subfile.toml",
			Name:     "",
		},
		{
			Path:     "test_data/root_file.toml",
			FileName: "root_file.toml",
			Name:     "",
		},
	}

	compareFileArrays(t, got, want)
}

func Test_GetTitleOfFile_GivenAFileWithoutAnExtension_ShouldReturnOriginal(t *testing.T) {
	filename := "some_filename_without_extension"
	title := GetTitleOfFile(filename)

	if title != filename {
		t.Errorf("given a file without an extension (%v) it did not return the original filename, instead got (%v)", filename, title)
	}
}

func Test_GetTitleOfFile_GivenAFileWithAnExtension_ShouldReturnFilenameWithoutExtension(t *testing.T) {
	filename := "some_file_with.extension"
	expectedTitle := "some_file_with"

	actualTitle := GetTitleOfFile(filename)
	if actualTitle != expectedTitle {
		t.Errorf("Expected title (%v) did not match actual title (%v)", expectedTitle, actualTitle)
	}
}

func Test_GetRelativePath_WithFileBeingInTheRoot_ShouldReturnOnlyFile(t *testing.T) {
	path := "/file.something"
	root := "/"
	expected := "file.something"

	actual := GetRelativePath(root, path)

	if expected != actual {
		t.Errorf("Actual relative path '%v' did not match expected relative path '%v'", actual, expected)
	}
}

func Test_GetRelativePath_WithFileBeingInASubfolder_ShouldReturnPathToSubfolder(t *testing.T) {
	path := "/root/subfolder/file.something"
	root := "/root/"
	expected := "subfolder/file.something"

	actual := GetRelativePath(root, path)

	if expected != actual {
		t.Errorf("Actual relative path '%v' did not match expected relative path '%v'", actual, expected)
	}
}

func Test_FindNonExistentVideos_WithTwoEmptySlices_ShouldReturnEmptySlice(t *testing.T) {
	existingMedia := []model.Media{}
	files := []File{}

	nonExistentMedia := FindNonExistentMedia(existingMedia, files)

	if len(nonExistentMedia) != 0 {
		t.Error("A file did not exist even though there was nothing to match. Phone your God.")
	}
}

func Test_FindNonExistentVideos_WithTwoSlicesThatHaveNoDifference_ShouldReturnEmptySlice(t *testing.T) {
	existingMedia := []model.Media{
		{
			Path: "some path",
		},
	}
	files := []File{
		{
			Path: "some path",
		},
	}

	nonExistentMedia := FindNonExistentMedia(existingMedia, files)

	if len(nonExistentMedia) != 0 {
		t.Error("Lists had identical relative paths and should have returned an empty list but didn't")
	}
}

func Test_FindNonExistentVideos_WithAllTheFilesExistingAndExtraFiles_ShouldReturnEmptySlice(t *testing.T) {
	existingMedia := []model.Media{
		{
			Path: "some path",
		},
	}
	files := []File{
		{
			Path: "some path",
		},
		{
			Path: "another path that should not make a difference",
		},
	}

	nonExistentMedia := FindNonExistentMedia(existingMedia, files)

	if len(nonExistentMedia) != 0 {
		t.Error("Found a relative path that should exist")
	}
}

func Test_FindNonExistentVideos_WithAVideoThatDoesNotExist_ShouldReturnASliceContainingVideo(t *testing.T) {
	existingMedia := []model.Media{
		{
			Path: "some path",
		},
	}
	files := []File{
		{
			Path: "another path that should not make a difference",
		},
	}

	nonExistentMedia := FindNonExistentMedia(existingMedia, files)

	if len(nonExistentMedia) != 1 {
		t.Error("No elements were returned in the non existent video search, when one was expected")
	}
	if !(nonExistentMedia[len(nonExistentMedia)-1].Path == existingMedia[len(existingMedia)-1].Path) {
		t.Error("Returned path did not match expected relative path")
	}
}
