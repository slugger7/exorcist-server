//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Video = newVideoTable("public", "video", "")

type videoTable struct {
	postgres.Table

	// Columns
	ID            postgres.ColumnString
	LibraryPathID postgres.ColumnString
	RelativePath  postgres.ColumnString
	Title         postgres.ColumnString
	FileName      postgres.ColumnString
	Height        postgres.ColumnInteger
	Width         postgres.ColumnInteger
	Runtime       postgres.ColumnFloat
	Size          postgres.ColumnInteger
	Checksum      postgres.ColumnString
	Added         postgres.ColumnTimestamp
	Deleted       postgres.ColumnBool
	Exists        postgres.ColumnBool
	Created       postgres.ColumnTimestamp
	Modified      postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type VideoTable struct {
	videoTable

	EXCLUDED videoTable
}

// AS creates new VideoTable with assigned alias
func (a VideoTable) AS(alias string) *VideoTable {
	return newVideoTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new VideoTable with assigned schema name
func (a VideoTable) FromSchema(schemaName string) *VideoTable {
	return newVideoTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new VideoTable with assigned table prefix
func (a VideoTable) WithPrefix(prefix string) *VideoTable {
	return newVideoTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new VideoTable with assigned table suffix
func (a VideoTable) WithSuffix(suffix string) *VideoTable {
	return newVideoTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newVideoTable(schemaName, tableName, alias string) *VideoTable {
	return &VideoTable{
		videoTable: newVideoTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newVideoTableImpl("", "excluded", ""),
	}
}

func newVideoTableImpl(schemaName, tableName, alias string) videoTable {
	var (
		IDColumn            = postgres.StringColumn("id")
		LibraryPathIDColumn = postgres.StringColumn("library_path_id")
		RelativePathColumn  = postgres.StringColumn("relative_path")
		TitleColumn         = postgres.StringColumn("title")
		FileNameColumn      = postgres.StringColumn("file_name")
		HeightColumn        = postgres.IntegerColumn("height")
		WidthColumn         = postgres.IntegerColumn("width")
		RuntimeColumn       = postgres.FloatColumn("runtime")
		SizeColumn          = postgres.IntegerColumn("size")
		ChecksumColumn      = postgres.StringColumn("checksum")
		AddedColumn         = postgres.TimestampColumn("added")
		DeletedColumn       = postgres.BoolColumn("deleted")
		ExistsColumn        = postgres.BoolColumn("exists")
		CreatedColumn       = postgres.TimestampColumn("created")
		ModifiedColumn      = postgres.TimestampColumn("modified")
		allColumns          = postgres.ColumnList{IDColumn, LibraryPathIDColumn, RelativePathColumn, TitleColumn, FileNameColumn, HeightColumn, WidthColumn, RuntimeColumn, SizeColumn, ChecksumColumn, AddedColumn, DeletedColumn, ExistsColumn, CreatedColumn, ModifiedColumn}
		mutableColumns      = postgres.ColumnList{LibraryPathIDColumn, RelativePathColumn, TitleColumn, FileNameColumn, HeightColumn, WidthColumn, RuntimeColumn, SizeColumn, ChecksumColumn, AddedColumn, DeletedColumn, ExistsColumn, CreatedColumn, ModifiedColumn}
	)

	return videoTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:            IDColumn,
		LibraryPathID: LibraryPathIDColumn,
		RelativePath:  RelativePathColumn,
		Title:         TitleColumn,
		FileName:      FileNameColumn,
		Height:        HeightColumn,
		Width:         WidthColumn,
		Runtime:       RuntimeColumn,
		Size:          SizeColumn,
		Checksum:      ChecksumColumn,
		Added:         AddedColumn,
		Deleted:       DeletedColumn,
		Exists:        ExistsColumn,
		Created:       CreatedColumn,
		Modified:      ModifiedColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
