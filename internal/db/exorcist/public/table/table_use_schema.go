//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	FavouriteMedia = FavouriteMedia.FromSchema(schema)
	FavouritePerson = FavouritePerson.FromSchema(schema)
	Image = Image.FromSchema(schema)
	Job = Job.FromSchema(schema)
	Library = Library.FromSchema(schema)
	LibraryPath = LibraryPath.FromSchema(schema)
	Media = Media.FromSchema(schema)
	MediaPerson = MediaPerson.FromSchema(schema)
	MediaProgress = MediaProgress.FromSchema(schema)
	MediaRelation = MediaRelation.FromSchema(schema)
	MediaTag = MediaTag.FromSchema(schema)
	Person = Person.FromSchema(schema)
	PersonAlias = PersonAlias.FromSchema(schema)
	Playlist = Playlist.FromSchema(schema)
	PlaylistMedia = PlaylistMedia.FromSchema(schema)
	SchemaMigrations = SchemaMigrations.FromSchema(schema)
	Tag = Tag.FromSchema(schema)
	TagAlias = TagAlias.FromSchema(schema)
	User = User.FromSchema(schema)
	Video = Video.FromSchema(schema)
}
