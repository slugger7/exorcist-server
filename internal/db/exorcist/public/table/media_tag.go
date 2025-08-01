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

var MediaTag = newMediaTagTable("public", "media_tag", "")

type mediaTagTable struct {
	postgres.Table

	// Columns
	ID       postgres.ColumnString
	Created  postgres.ColumnTimestamp
	Modified postgres.ColumnTimestamp
	MediaID  postgres.ColumnString
	TagID    postgres.ColumnString
	GhostID  postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type MediaTagTable struct {
	mediaTagTable

	EXCLUDED mediaTagTable
}

// AS creates new MediaTagTable with assigned alias
func (a MediaTagTable) AS(alias string) *MediaTagTable {
	return newMediaTagTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new MediaTagTable with assigned schema name
func (a MediaTagTable) FromSchema(schemaName string) *MediaTagTable {
	return newMediaTagTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new MediaTagTable with assigned table prefix
func (a MediaTagTable) WithPrefix(prefix string) *MediaTagTable {
	return newMediaTagTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new MediaTagTable with assigned table suffix
func (a MediaTagTable) WithSuffix(suffix string) *MediaTagTable {
	return newMediaTagTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newMediaTagTable(schemaName, tableName, alias string) *MediaTagTable {
	return &MediaTagTable{
		mediaTagTable: newMediaTagTableImpl(schemaName, tableName, alias),
		EXCLUDED:      newMediaTagTableImpl("", "excluded", ""),
	}
}

func newMediaTagTableImpl(schemaName, tableName, alias string) mediaTagTable {
	var (
		IDColumn       = postgres.StringColumn("id")
		CreatedColumn  = postgres.TimestampColumn("created")
		ModifiedColumn = postgres.TimestampColumn("modified")
		MediaIDColumn  = postgres.StringColumn("media_id")
		TagIDColumn    = postgres.StringColumn("tag_id")
		GhostIDColumn  = postgres.IntegerColumn("ghost_id")
		allColumns     = postgres.ColumnList{IDColumn, CreatedColumn, ModifiedColumn, MediaIDColumn, TagIDColumn, GhostIDColumn}
		mutableColumns = postgres.ColumnList{CreatedColumn, ModifiedColumn, MediaIDColumn, TagIDColumn, GhostIDColumn}
	)

	return mediaTagTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:       IDColumn,
		Created:  CreatedColumn,
		Modified: ModifiedColumn,
		MediaID:  MediaIDColumn,
		TagID:    TagIDColumn,
		GhostID:  GhostIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
