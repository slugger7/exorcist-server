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
	ID      postgres.ColumnString
	MediaID postgres.ColumnString
	Height  postgres.ColumnInteger
	Width   postgres.ColumnInteger
	Runtime postgres.ColumnFloat
	GhostID postgres.ColumnInteger

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
		IDColumn       = postgres.StringColumn("id")
		MediaIDColumn  = postgres.StringColumn("media_id")
		HeightColumn   = postgres.IntegerColumn("height")
		WidthColumn    = postgres.IntegerColumn("width")
		RuntimeColumn  = postgres.FloatColumn("runtime")
		GhostIDColumn  = postgres.IntegerColumn("ghost_id")
		allColumns     = postgres.ColumnList{IDColumn, MediaIDColumn, HeightColumn, WidthColumn, RuntimeColumn, GhostIDColumn}
		mutableColumns = postgres.ColumnList{MediaIDColumn, HeightColumn, WidthColumn, RuntimeColumn, GhostIDColumn}
	)

	return videoTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:      IDColumn,
		MediaID: MediaIDColumn,
		Height:  HeightColumn,
		Width:   WidthColumn,
		Runtime: RuntimeColumn,
		GhostID: GhostIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
