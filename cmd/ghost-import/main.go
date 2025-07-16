package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	gmodel "github.com/slugger7/exorcist/internal/db/ghost/model"
	gtable "github.com/slugger7/exorcist/internal/db/ghost/table"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type UserMap struct {
	ExorcistUser string `json:"exorcist_user"`
	GhostUser    string `json:"ghost_user"`
}

type LibraryConfig struct {
	Exclude []string `json:"exclude"`
}

type Config struct {
	UserMap []UserMap     `json:"userMap"`
	Library LibraryConfig `json:"library"`
}

type Context struct {
	ExorcistDb *sql.DB
	GhostDb    *sql.DB
	Config     *Config
}

func createPostgresDb() *sql.DB {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Printf("Opening DB: %v\n", psqlconn)

	db, err := sql.Open("postgres", psqlconn)
	errs.PanicError(err)

	return db
}

func createSqlLiteDb() *sql.DB {
	dbPath := os.Getenv("GHOST_DB_PATH")

	db, err := sql.Open("sqlite3", dbPath)
	errs.PanicError(err)

	return db
}

func parseConfig(filePath string) (*Config, error) {
	userMapFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	userMapBytes, err := io.ReadAll(userMapFile)
	if err != nil {
		return nil, err
	}

	var config Config

	err = json.Unmarshal(userMapBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (ctx Context) transferUsers() error {
	gUserStatement := gtable.Users.SELECT(gtable.Users.AllColumns)

	var gUsers []gmodel.Users
	if err := gUserStatement.Query(ctx.GhostDb, &gUsers); err != nil {
		return err
	}

	if len(gUsers) == 0 {
		return fmt.Errorf("no users found in ghost")
	}

	userStatement := table.User.SELECT(table.User.AllColumns)

	var users []model.User
	if err := userStatement.Query(ctx.ExorcistDb, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		return fmt.Errorf("no users found in exorcist")
	}

	var accErrs error
	for _, mu := range ctx.Config.UserMap {
		var guser *gmodel.Users
		for _, gu := range gUsers {
			if gu.Username == mu.GhostUser {
				guser = &gu
				break
			}
		}

		if guser == nil {
			accErrs = errors.Join(accErrs, fmt.Errorf("could not find user in ghost: %v", mu.GhostUser))
			continue
		}

		var user *model.User
		for _, u := range users {
			if u.Username == mu.ExorcistUser {
				user = &u
				break
			}
		}

		if user == nil {
			accErrs = errors.Join(accErrs, fmt.Errorf("could not find user in exorcist: %v", mu.ExorcistUser))
			continue
		}

		user.GhostID = &guser.ID

		updateStmnt := table.User.UPDATE(table.User.GhostID).
			MODEL(user).
			WHERE(table.User.ID.EQ(postgres.UUID(user.ID)))

		res, err := updateStmnt.Exec(ctx.ExorcistDb)
		if err != nil {
			accErrs = errors.Join(accErrs, err)
			continue
		}

		rows, _ := res.RowsAffected()

		log.Printf("Altered %v rows for user %v in exorcist", rows, mu.ExorcistUser)
	}

	return accErrs
}

func (ctx *Context) transferLibraries() error {
	log.Println("Transferring libraries")

	ghostStatement := gtable.Libraries.SELECT(gtable.Libraries.AllColumns)

	var ghostLibraries []gmodel.Libraries

	if err := ghostStatement.Query(ctx.GhostDb, &ghostLibraries); err != nil {
		return err
	}

	log.Printf("Found %v libraries in ghost\n", len(ghostLibraries))

	exorcistLibraries := []model.Library{}

	for _, l := range ghostLibraries {
		if slices.Contains(ctx.Config.Library.Exclude, l.Name) {
			continue
		}
		exorcistLibraries = append(exorcistLibraries, model.Library{
			Name:    l.Name,
			GhostID: &l.ID,
		})
	}

	insertStmnt := table.Library.INSERT(table.Library.GhostID, table.Library.Name).
		MODELS(exorcistLibraries).
		ON_CONFLICT(table.Library.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	log.Printf("Altered %v rows in exorcist libraries\n", rows)

	return nil
}

func (ctx *Context) transferLibraryPaths() error {
	log.Println("Transferring library paths")

	librariesStatement := table.Library.SELECT(table.Library.AllColumns)

	var libraries []model.Library
	if err := librariesStatement.Query(ctx.ExorcistDb, &libraries); err != nil {
		return err
	}

	libraryPaths := []model.LibraryPath{}
	for _, l := range libraries {
		if l.GhostID == nil {
			break
		}

		stmnt := gtable.LibraryPaths.SELECT(gtable.LibraryPaths.AllColumns).
			WHERE(gtable.LibraryPaths.ID.EQ(postgres.Int32(*l.GhostID)))

		var ghostLibraryPaths []gmodel.LibraryPaths
		if err := stmnt.Query(ctx.GhostDb, &ghostLibraryPaths); err != nil {
			return err
		}

		for _, lp := range ghostLibraryPaths {
			libraryPaths = append(libraryPaths, model.LibraryPath{
				GhostID:   &lp.ID,
				LibraryID: l.ID,
				Path:      lp.Path,
			})
		}
	}

	log.Printf("Found %v library paths in ghost\n", len(libraryPaths))

	insertStmnt := table.LibraryPath.INSERT(table.LibraryPath.GhostID, table.LibraryPath.LibraryID, table.LibraryPath.Path).
		MODELS(libraryPaths).
		ON_CONFLICT(table.LibraryPath.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	log.Printf("Altered %v rows in exorcist library paths\n", rows)

	return nil
}

func (ctx *Context) transferVideos() error {
	log.Println("Transferring videos")

	libraryPathStmnt := table.LibraryPath.SELECT(table.LibraryPath.AllColumns).
		WHERE(table.LibraryPath.GhostID.IS_NOT_NULL())

	var libraryPaths []model.LibraryPath
	if err := libraryPathStmnt.Query(ctx.ExorcistDb, &libraryPaths); err != nil {
		return err
	}

	if len(libraryPaths) == 0 {
		return nil
	}

	for i, lp := range libraryPaths {
		log.Printf("Transferring library path %v/%v\n", i+1, len(libraryPaths))
		videoStmnt := gtable.Videos.SELECT(gtable.Videos.AllColumns).
			WHERE(gtable.Videos.LibraryPathId.EQ(postgres.Int32(*lp.GhostID)))

		var gVideos []gmodel.Videos
		if err := videoStmnt.Query(ctx.GhostDb, &gVideos); err != nil {
			return err
		}

		if len(gVideos) == 0 {
			log.Printf("Skipping library path %v due to not having videos to transfer", i+1)
			continue
		}

		log.Printf("Found %v videos to transfer for library path %v\n", len(gVideos), i+1)

		gVideosMap := make(map[int32]gmodel.Videos)
		mediaEntitiesMap := make([]model.Media, len(gVideos))
		for x, gv := range gVideos {
			gVideosMap[gv.ID] = gv

			addedTime, _ := time.Parse(time.DateTime, gv.DateAdded)
			createdTime, _ := time.Parse(time.DateTime, gv.Created)
			mediaEntitiesMap[x] = model.Media{
				LibraryPathID: lp.ID,
				Path:          gv.Path,
				Title:         gv.Title,
				MediaType:     model.MediaTypeEnum_Primary,
				Size:          int64(gv.Size),
				Added:         addedTime,
				Created:       createdTime,
				GhostID:       &gv.ID,
			}
		}

		media := table.Media
		mediaInsertStmnt := media.INSERT(media.LibraryPathID, media.Path, media.Title, media.MediaType, media.Size, media.Added, media.Created, media.GhostID).
			MODELS(mediaEntitiesMap).
			RETURNING(media.GhostID, media.ID).
			ON_CONFLICT(media.GhostID).DO_NOTHING()

		var mediaEntities []model.Media
		if err := mediaInsertStmnt.Query(ctx.ExorcistDb, &mediaEntities); err != nil {
			return err
		}

		if len(mediaEntities) == 0 {
			continue
		}

		log.Printf("Altered %v rows in exorcist media\n", len(mediaEntities))

		videoEntities := make([]model.Video, len(mediaEntities))
		for y, m := range mediaEntities {
			gVid := gVideosMap[*m.GhostID]
			videoEntities[y] = model.Video{
				MediaID: m.ID,
				GhostID: m.GhostID,
				Height:  gVid.Height,
				Width:   gVid.Width,
				Runtime: float64(gVid.Runtime), // TODO: these might be in different formats
			}

			// TODO: create job to generate thumbnail for each media piece
		}

		insertVidStmnt := table.Video.INSERT(table.Video.MediaID, table.Video.GhostID, table.Video.Height, table.Video.Width, table.Video.Runtime).
			MODELS(videoEntities).
			ON_CONFLICT(table.Video.GhostID).DO_NOTHING()

		res, err := insertVidStmnt.Exec(ctx.ExorcistDb)
		if err != nil {
			return err
		}

		rows, _ := res.RowsAffected()

		log.Printf("Altered %v rows in exorcist video\n", rows)
	}
	return nil
}

func (ctx Context) transferActors() error {
	log.Println("Transferring actors")
	actorsStmnt := gtable.Actors.SELECT(gtable.Actors.AllColumns)

	var actors []gmodel.Actors
	if err := actorsStmnt.Query(ctx.GhostDb, &actors); err != nil {
		return err
	}

	if len(actors) == 0 {
		return nil
	}

	log.Printf("Found %v actors in ghost\n", len(actors))

	people := make([]model.Person, len(actors))
	for i, a := range actors {
		people[i] = model.Person{
			GhostID: &a.ID,
			Name:    a.Name,
		}
	}

	insertStmnt := table.Person.INSERT(table.Person.GhostID, table.Person.Name).
		MODELS(people).
		ON_CONFLICT(table.Person.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in exorcist person\n", rows)

	return nil
}

func (ctx Context) transferGenres() error {
	log.Println("Transferring genres")
	genresStmnt := gtable.Genres.SELECT(gtable.Genres.AllColumns)

	var genres []gmodel.Genres
	if err := genresStmnt.Query(ctx.GhostDb, &genres); err != nil {
		return err
	}

	if len(genres) == 0 {
		return nil
	}

	log.Printf("Found %v genres in ghost\n", len(genres))

	tags := make([]model.Tag, len(genres))
	for i, g := range genres {
		tags[i] = model.Tag{
			GhostID: &g.ID,
			Name:    g.Name,
		}
	}

	insertStmnt := table.Tag.INSERT(table.Tag.GhostID, table.Tag.Name).
		MODELS(tags).
		ON_CONFLICT(table.Tag.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in exorcist tag", rows)

	return nil
}

func (ctx Context) transferPlaylists() error {
	log.Println("Trasferring playlists")

	playlistStmnt := gtable.Playlists.SELECT(gtable.Playlists.AllColumns)

	var gPlaylists []gmodel.Playlists
	if err := playlistStmnt.Query(ctx.GhostDb, &gPlaylists); err != nil {
		return err
	}

	if len(gPlaylists) == 0 {
		return nil
	}

	log.Printf("Found %v rows in ghost playlist\n", len(gPlaylists))

	userPlaylistMap := make(map[int32][]gmodel.Playlists)
	for _, gp := range gPlaylists {
		userPlaylistMap[gp.UserId] = append(userPlaylistMap[gp.UserId], gp)
	}

	playlists := []model.Playlist{}
	for id := range userPlaylistMap {
		log.Println(id)
		userStmnt := table.User.SELECT(table.User.ID, table.User.GhostID).
			WHERE(table.User.GhostID.EQ(postgres.Int32(id)))

		var users []model.User
		if err := userStmnt.Query(ctx.ExorcistDb, &users); err != nil {
			return err
		}

		if len(users) == 0 {
			log.Println("Colud not find user")
			return nil
		}

		user := users[0]

		for _, p := range userPlaylistMap[id] {
			createdTime, _ := time.Parse(time.DateTime, p.CreatedAt)
			playlists = append(playlists, model.Playlist{
				GhostID: &p.ID,
				UserID:  user.ID,
				Name:    p.Name,
				Created: createdTime,
			})
		}
	}

	if len(playlists) == 0 {
		log.Println("No playlists left to create")
		return nil
	}

	insertStmnt := table.Playlist.INSERT(table.Playlist.GhostID, table.Playlist.UserID, table.Playlist.Name, table.Playlist.Created).
		MODELS(playlists).
		ON_CONFLICT(table.Playlist.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in exorcist playlist\n", rows)

	return nil
}

func (ctx Context) transferPlaylistVideos() error {
	log.Println("Transferring playlist videos")

	playlistVideoStmnt := gtable.PlaylistVideos.SELECT(gtable.PlaylistVideos.AllColumns)

	var gPlaylistVideos []gmodel.PlaylistVideos
	if err := playlistVideoStmnt.Query(ctx.GhostDb, &gPlaylistVideos); err != nil {
		return err
	}

	if len(gPlaylistVideos) == 0 {
		log.Println("No playlist videos found")
		return nil
	}

	playlistStmnt := table.Playlist.SELECT(table.Playlist.AllColumns)

	var playlists []model.Playlist
	if err := playlistStmnt.Query(ctx.ExorcistDb, &playlists); err != nil {
		return err
	}

	if len(playlists) == 0 {
		log.Println("No playlists found")
		return nil
	}

	videoStmnt := table.Media.SELECT(table.Media.AllColumns)

	var videoMedia []model.Media
	if err := videoStmnt.Query(ctx.ExorcistDb, &videoMedia); err != nil {
		return err
	}

	if len(videoMedia) == 0 {
		log.Println("No videos found")
		return nil
	} else {
		log.Printf("Found %v videos in exorcist to add to playlists", len(videoMedia))
	}

	playlistMedia := []model.PlaylistMedia{}
	for _, pv := range gPlaylistVideos {
		var playlist *model.Playlist
		for _, p := range playlists {
			if *p.GhostID == pv.PlaylistId {
				playlist = &p
				break
			}
		}

		if playlist == nil {
			log.Println("playlist was nil")
			continue
		}

		var media *model.Media
		for _, m := range videoMedia {
			if *m.GhostID == pv.VideoId {
				media = &m
				break
			}
		}

		if media == nil {
			continue
		}

		createdTime, _ := time.Parse(time.DateTime, pv.CreatedAt)

		playlistMedia = append(playlistMedia, model.PlaylistMedia{
			GhostID:    &pv.ID,
			MediaID:    media.ID,
			PlaylistID: playlist.ID,
			Created:    createdTime,
		})
	}

	insertStmnt := table.PlaylistMedia.INSERT(table.PlaylistMedia.GhostID, table.PlaylistMedia.MediaID, table.PlaylistMedia.PlaylistID, table.PlaylistMedia.Created).
		MODELS(playlistMedia).
		ON_CONFLICT(table.PlaylistMedia.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in exorcist playlist media\n", rows)

	return nil
}

func main() {
	err := godotenv.Load()
	errs.PanicError(err)

	pgDb := createPostgresDb()
	defer pgDb.Close()

	err = pgDb.Ping()
	errs.PanicError(err)

	sqlLiteDb := createSqlLiteDb()
	defer sqlLiteDb.Close()

	err = sqlLiteDb.Ping()
	errs.PanicError(err)

	config, err := parseConfig("./.temp/ghost_import.config.json")
	if err != nil {
		errs.PanicError(err)
	}

	ctx := Context{
		ExorcistDb: pgDb,
		GhostDb:    sqlLiteDb,
		Config:     config,
	}

	err = ctx.transferUsers()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferLibraries()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferLibraryPaths()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferVideos()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferActors()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferGenres()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferPlaylists()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferPlaylistVideos()
	if err != nil {
		errs.PanicError(err)
	}

	// TODO: favourite_actors
	// TODO: favourite_videos
	// TODO: progress
	// TODO: related video
	// TODO: video actors
	// TODO: video genres
}
