package videoRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func GetVideoWithoutChecksumStatement() postgres.SelectStatement {
	selectStatement := table.Video.SELECT(table.Video.ID, table.Video.Checksum, table.Video.RelativePath, table.LibraryPath.Path).
		FROM(table.Video.INNER_JOIN(table.LibraryPath, table.LibraryPath.ID.EQ(table.Video.LibraryPathID))).
		WHERE(table.Video.Checksum.IS_NULL())

	repository.DebugCheck(selectStatement)

	return selectStatement
}

func UpdateVideoChecksum(video model.Video) postgres.UpdateStatement {
	statement := table.Video.UPDATE().
		SET(table.Video.Checksum.SET(postgres.String(*video.Checksum))).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	repository.DebugCheck(statement)

	return statement
}

func UpdateVideoExistsStatement(video model.Video) postgres.UpdateStatement {
	statement := table.Video.UPDATE().
		SET(table.Video.Exists.SET(postgres.Bool(video.Exists))).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	repository.DebugCheck(statement)

	return statement
}

func GetVideosInLibraryPath(libraryPathId uuid.UUID) postgres.SelectStatement {
	statement := table.Video.SELECT(table.Video.RelativePath, table.Video.ID).
		FROM(table.Video.Table).
		WHERE(table.Video.LibraryPathID.EQ(postgres.UUID(libraryPathId)).
			AND(table.Video.Exists.IS_TRUE()),
		)

	repository.DebugCheck(statement)

	return statement
}

func InsertVideosStatement(videos []model.Video) postgres.InsertStatement {
	if len(videos) == 0 {
		return nil
	}
	statement := table.Video.INSERT(
		table.Video.LibraryPathID,
		table.Video.RelativePath,
		table.Video.Title,
		table.Video.FileName,
		table.Video.Height,
		table.Video.Width,
		table.Video.Runtime,
		table.Video.Size,
	).
		MODELS(videos)

	repository.DebugCheck(statement)

	return statement
}

func QuerySelect(db *sql.DB, statement postgres.SelectStatement) (data []struct{ model.Video }, err error) {
	err = statement.Query(db, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
