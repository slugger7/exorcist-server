package repository_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository"
)

// TODO: implement [snapshot tests](https://github.com/gkampitakis/go-snaps)

func Test_GetVideoWithoutChecksumStatement(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	actual, _ := ds.GetVideoWithoutChecksumStatement().Sql()

	expected := "\nSELECT video.id AS \"video.id\",\n     video.checksum AS \"video.checksum\",\n     video.relative_path AS \"video.relative_path\",\n     library_path.path AS \"library_path.path\"\nFROM public.video\n     INNER JOIN public.library_path ON (library_path.id = video.library_path_id)\nWHERE video.checksum IS NULL;\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_UpdateVideoChecksum(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}

	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}

	checksum := "someChecksum"

	video := model.Video{
		ID:       newUuid,
		Checksum: &checksum,
	}

	actual, _ := ds.UpdateVideoChecksum(video).Sql()

	expected := `
UPDATE public.video
SET checksum = $1::text
WHERE video.id = $2;
`
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func Test_MarkVideoAsNotExistingStatement(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}

	video := model.Video{
		ID:     newUuid,
		Exists: false,
	}

	actual, _ := ds.UpdateVideoExistsStatement(video).Sql()

	expected := `
UPDATE public.video
SET exists = $1::boolean
WHERE video.id = $2;
`
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func Test_GetVideosInLibraryPath(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}
	actual, _ := ds.GetVideosInLibraryPath(newUuid).Sql()

	expected := `
SELECT video.relative_path AS "video.relative_path",
     video.id AS "video.id"
FROM public.video
WHERE (video.library_path_id = $1) AND video.exists IS TRUE;
`
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_InsertVideosStatement_WithNoVideos_ShouldReturnNil(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	videos := []model.Video{}
	actual := ds.InsertVideosStatement(videos)

	if actual != nil {
		t.Errorf("Expected actual to be nil. Acutal: %v", actual)
	}
}

func Test_InsertVideosStatement_WithVideos_ShouldReturnStatement(t *testing.T) {
	ds := &repository.DatabaseService{
		Env: &environment.EnvironmentVariables{DebugSql: false},
	}
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}
	video := model.Video{
		LibraryPathID: newUuid,
		RelativePath:  "relativePath",
		Title:         "title",
		FileName:      "filename",
		Height:        69,
		Width:         420,
		Runtime:       1337,
		Size:          80085,
	}
	videos := []model.Video{video}
	actual, _ := ds.InsertVideosStatement(videos).Sql()

	expected :=
		`
INSERT INTO public.video (library_path_id, relative_path, title, file_name, height, width, runtime, size)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`
	if actual != expected {
		t.Errorf("Expected \n%v got \n%v", expected, actual)
	}
}
