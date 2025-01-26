package videoRepository_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

func Test_GetVideoWithoutChecksumStatement(t *testing.T) {
	actual := videoRepository.GetVideoWithoutChecksumStatement().DebugSql()

	expected := "\nSELECT video.id AS \"video.id\",\n     video.checksum AS \"video.checksum\",\n     video.relative_path AS \"video.relative_path\",\n     library_path.path AS \"library_path.path\"\nFROM public.video\n     INNER JOIN public.library_path ON (library_path.id = video.library_path_id)\nWHERE video.checksum IS NULL;\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_UpdateVideoChecksum(t *testing.T) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}

	checksum := "someChecksum"

	video := model.Video{
		ID:       newUuid,
		Checksum: &checksum,
	}

	actual := videoRepository.UpdateVideoChecksum(video).DebugSql()

	expected := fmt.Sprintf("\nUPDATE public.video\nSET checksum = '%v'::text\nWHERE video.id = '%v';\n", checksum, newUuid)
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}
