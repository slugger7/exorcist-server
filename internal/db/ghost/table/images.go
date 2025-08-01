//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/sqlite"
)

var Images = newImagesTable("", "Images", "")

type imagesTable struct {
	sqlite.Table

	// Columns
	ID   sqlite.ColumnInteger
	Name sqlite.ColumnString
	Path sqlite.ColumnString

	AllColumns     sqlite.ColumnList
	MutableColumns sqlite.ColumnList
}

type ImagesTable struct {
	imagesTable

	EXCLUDED imagesTable
}

// AS creates new ImagesTable with assigned alias
func (a ImagesTable) AS(alias string) *ImagesTable {
	return newImagesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ImagesTable with assigned schema name
func (a ImagesTable) FromSchema(schemaName string) *ImagesTable {
	return newImagesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ImagesTable with assigned table prefix
func (a ImagesTable) WithPrefix(prefix string) *ImagesTable {
	return newImagesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ImagesTable with assigned table suffix
func (a ImagesTable) WithSuffix(suffix string) *ImagesTable {
	return newImagesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newImagesTable(schemaName, tableName, alias string) *ImagesTable {
	return &ImagesTable{
		imagesTable: newImagesTableImpl(schemaName, tableName, alias),
		EXCLUDED:    newImagesTableImpl("", "excluded", ""),
	}
}

func newImagesTableImpl(schemaName, tableName, alias string) imagesTable {
	var (
		IDColumn       = sqlite.IntegerColumn("Id")
		NameColumn     = sqlite.StringColumn("Name")
		PathColumn     = sqlite.StringColumn("Path")
		allColumns     = sqlite.ColumnList{IDColumn, NameColumn, PathColumn}
		mutableColumns = sqlite.ColumnList{NameColumn, PathColumn}
	)

	return imagesTable{
		Table: sqlite.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:   IDColumn,
		Name: NameColumn,
		Path: PathColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
