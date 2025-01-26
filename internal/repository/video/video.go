package videoRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func GetVideoWithoutChecksumStatement() postgres.SelectStatement {
	selectStatement := table.Video.SELECT(table.Video.ID, table.Video.Checksum, table.Video.RelativePath, table.LibraryPath.Path).
		FROM(table.Video.INNER_JOIN(table.LibraryPath, table.LibraryPath.ID.EQ(table.Video.LibraryPathID))).
		WHERE(table.Video.Checksum.IS_NULL())

	repository.DebugCheckSelect(selectStatement)

	return selectStatement
}

func ExecuteChecksumStatement(db *sql.DB, statement postgres.SelectStatement) (data []struct {
	model.LibraryPath
	model.Video
}, err error) {
	err = statement.Query(db, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func UpdateVideoChecksum(video model.Video) postgres.UpdateStatement {
	statement := table.Video.UPDATE().
		SET(table.Video.Checksum.SET(postgres.String(*video.Checksum))).
		MODEL(video).
		WHERE(table.Video.ID.EQ(postgres.UUID(video.ID)))

	repository.DebugCheckUpdate(statement)

	return statement
}

func ExecuteUpdate(db *sql.DB, statement postgres.UpdateStatement) (err error) {
	_, err = statement.Exec(db)
	return err
}
