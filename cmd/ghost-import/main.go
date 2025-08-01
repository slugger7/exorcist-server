package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	"github.com/slugger7/exorcist/internal/job"
)

type UserMap struct {
	ExorcistUser string `json:"exorcist_user"`
	GhostUser    string `json:"ghost_user"`
}

type LibraryConfig struct {
	Exclude []string `json:"exclude"`
}

type Config struct {
	BatchSize int           `json:"batchSize"`
	UserMap   []UserMap     `json:"userMap"`
	Library   LibraryConfig `json:"library"`
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

	if len(ghostLibraries) != int(rows) {
		log.Println("ROWS FOUND VS ROWS ALTERED DID NOT MATCH FOR TRANSFERRING OF LIBRARIES")
	}

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
			RETURNING(media.GhostID, media.ID, media.Path).
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
				Runtime: float64(gVid.Runtime / 1000),
			}

		}

		insertVidStmnt := table.Video.INSERT(table.Video.MediaID, table.Video.GhostID, table.Video.Height, table.Video.Width, table.Video.Runtime).
			MODELS(videoEntities).
			RETURNING(table.Video.ID, table.Video.MediaID, table.Video.GhostID).
			ON_CONFLICT(table.Video.GhostID).DO_NOTHING()

		var insertedVids []model.Video
		if err := insertVidStmnt.Query(ctx.ExorcistDb, &insertedVids); err != nil {
			return err
		}

		if len(insertedVids) == 0 {
			log.Println("No videos were inserted and returned")
			return nil
		}

		rows := int64(len(insertedVids))

		log.Printf("Altered %v rows in exorcist video\n", rows)

		thumbnailJobs := []model.Job{}
		for _, v := range insertedVids {
			var mediaEntity *model.Media
			for _, m := range mediaEntities {
				if *m.GhostID == *v.GhostID {
					mediaEntity = &m
					break
				}
			}

			if mediaEntity == nil {
				continue
			}

			assetPath := filepath.Join(os.Getenv("ASSETS"), v.MediaID.String(), fmt.Sprintf("%v.webp", filepath.Base(mediaEntity.Path)))
			job, err := job.CreateGenerateThumbnailJob(v, nil, assetPath, 0, int(v.Height), int(v.Width))
			if err != nil {
				continue
			}

			thumbnailJobs = append(thumbnailJobs, *job)
		}

		insertStment := table.Job.INSERT(table.Job.JobType, table.Job.Status, table.Job.Data, table.Job.Parent, table.Job.Priority).
			MODELS(thumbnailJobs)

		res, err := insertStment.Exec(ctx.ExorcistDb)
		if err != nil {
			return err
		}

		rows, _ = res.RowsAffected()

		log.Printf("Created %v thumbnail creation jobs", rows)
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

	log.Printf("Found %v playlist videos in ghost", len(gPlaylistVideos))

	playlistStmnt := table.Playlist.SELECT(table.Playlist.ID, table.Playlist.GhostID)

	var playlists []model.Playlist
	if err := playlistStmnt.Query(ctx.ExorcistDb, &playlists); err != nil {
		return err
	}

	if len(playlists) == 0 {
		log.Println("No playlists found")
		return nil
	}

	videoStmnt := table.Media.SELECT(table.Media.ID, table.Media.GhostID)

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

	noPlaylistCount := 0
	noMediaCount := 0
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
			noPlaylistCount++
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
			noMediaCount++

			stmnt := gtable.Videos.SELECT(gtable.Videos.ID).
				WHERE(gtable.Videos.ID.EQ(postgres.Int32(pv.VideoId)))

			var vid []gmodel.Videos
			if err := stmnt.Query(ctx.GhostDb, &vid); err != nil {
				return err
			}

			if len(vid) > 0 {
				log.Println("Video does exist in ghost", pv.VideoId)
			}

			pstmnt := table.Media.SELECT(table.Media.ID, table.Media.GhostID).
				WHERE(table.Media.GhostID.EQ(postgres.Int32(pv.VideoId)))

			var med []model.Media
			if err := pstmnt.Query(ctx.ExorcistDb, &med); err != nil {
				return err
			}

			if len(med) > 0 {
				log.Println("Media was find by ghost id in exorcist")
			}
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

	log.Printf("Created %v playlist media entities to insert, no playlist count %v, no media count %v", len(playlistMedia), noPlaylistCount, noMediaCount)

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

func (ctx Context) transferFavouriteActors() error {
	log.Println("Transferring favourite actors")

	favouriteActorsStmnt := gtable.FavouriteActors.SELECT(gtable.FavouriteActors.AllColumns)

	var gFavouriteActors []gmodel.FavouriteActors
	if err := favouriteActorsStmnt.Query(ctx.GhostDb, &gFavouriteActors); err != nil {
		return err
	}

	if len(gFavouriteActors) == 0 {
		log.Println("No favourite actors in ghost")
		return nil
	}

	log.Printf("Found %v favourite actors in ghost", len(gFavouriteActors))

	userStmnt := table.User.SELECT(table.User.ID, table.User.GhostID)

	var users []model.User
	if err := userStmnt.Query(ctx.ExorcistDb, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users in exorcist")
		return nil
	}

	peopleStmnt := table.Person.SELECT(table.Person.AllColumns)

	var people []model.Person
	if err := peopleStmnt.Query(ctx.ExorcistDb, &people); err != nil {
		return err
	}

	if len(people) == 0 {
		log.Println("No people in exorcist")
		return nil
	}

	favouritePeople := []model.FavouritePerson{}
	for _, fa := range gFavouriteActors {
		var user *model.User
		for _, u := range users {
			if fa.UserId == *u.GhostID {
				user = &u
				break
			}
		}

		if user == nil {
			continue
		}

		var person *model.Person
		for _, p := range people {
			if fa.ActorId == *p.GhostID {
				person = &p
				break
			}
		}

		if person == nil {
			continue
		}

		favouritePeople = append(favouritePeople, model.FavouritePerson{
			GhostID:  &fa.ID,
			UserID:   user.ID,
			PersonID: person.ID,
		})
	}

	if len(favouritePeople) == 0 {
		log.Println("No favourite people to add into exorcist")
		return nil
	}

	insertStmnt := table.FavouritePerson.INSERT(table.FavouritePerson.GhostID, table.FavouritePerson.UserID, table.FavouritePerson.PersonID).
		MODELS(favouritePeople).
		ON_CONFLICT(table.FavouritePerson.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in favourite people in exorcist", rows)

	return nil
}

func (ctx Context) transferFavouriteVideos() error {
	log.Println("Transferring favourite videos")

	favouriteVideosStmnt := gtable.FavouriteVideos.SELECT(gtable.FavouriteVideos.AllColumns)

	var favouriteVideos []gmodel.FavouriteVideos
	if err := favouriteVideosStmnt.Query(ctx.GhostDb, &favouriteVideos); err != nil {
		return err
	}

	if len(favouriteVideos) == 0 {
		log.Println("No favourite videos in ghost")
		return nil
	}

	log.Printf("Found %v favourite videos in ghost", len(favouriteVideos))

	userStmnt := table.User.SELECT(table.User.ID, table.User.GhostID)

	var users []model.User
	if err := userStmnt.Query(ctx.ExorcistDb, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users in exorcist")
		return nil
	}

	mediaStmnt := table.Media.SELECT(table.Media.AllColumns)

	var mediaList []model.Media
	if err := mediaStmnt.Query(ctx.ExorcistDb, &mediaList); err != nil {
		return err
	}

	if len(mediaList) == 0 {
		log.Println("No media in exorcist")
		return nil
	}

	favouriteMedia := []model.FavouriteMedia{}
	for _, fa := range favouriteVideos {
		var user *model.User
		for _, u := range users {
			if fa.UserId == *u.GhostID {
				user = &u
				break
			}
		}

		if user == nil {
			continue
		}

		var media *model.Media
		for _, p := range mediaList {
			if fa.VideoId == *p.GhostID {
				media = &p
				break
			}
		}

		if media == nil {
			continue
		}

		favouriteMedia = append(favouriteMedia, model.FavouriteMedia{
			GhostID: &fa.ID,
			UserID:  user.ID,
			MediaID: media.ID,
		})
	}

	if len(favouriteMedia) == 0 {
		log.Println("No favourite media to add into exorcist")
		return nil
	}

	insertStmnt := table.FavouriteMedia.INSERT(table.FavouriteMedia.GhostID, table.FavouriteMedia.UserID, table.FavouriteMedia.MediaID).
		MODELS(favouriteMedia).
		ON_CONFLICT(table.FavouriteMedia.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in favourite media in exorcist", rows)

	return nil
}

func (ctx Context) transferProgress() error {
	log.Println("Transferring progress")
	progressStmnt := gtable.Progress.SELECT(gtable.Progress.AllColumns)

	var gProgress []gmodel.Progress
	if err := progressStmnt.Query(ctx.GhostDb, &gProgress); err != nil {
		return err
	}

	if len(gProgress) == 0 {
		log.Println("No progress found in ghost")
		return nil
	}

	userStmnt := table.User.SELECT(table.User.GhostID, table.User.ID)

	var users []model.User
	if err := userStmnt.Query(ctx.ExorcistDb, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		log.Println("No users found in exorcist")
		return nil
	}

	mediaStatement := table.Media.SELECT(table.Media.ID, table.Media.GhostID)

	var mediaList []model.Media
	if err := mediaStatement.Query(ctx.ExorcistDb, &mediaList); err != nil {
		return err
	}

	if len(mediaList) == 0 {
		log.Println("No media found in exorcist")
		return nil
	}

	var accErrs error
	var progressList []model.MediaProgress
	for _, gp := range gProgress {
		var user *model.User
		for _, u := range users {
			if gp.UserId == *u.GhostID {
				user = &u
				break
			}
		}

		if user == nil {
			accErrs = fmt.Errorf("no user found in exorcist while transferring progress for ghost id: %v\n %w", gp.UserId, accErrs)
			continue
		}

		var media *model.Media
		for _, m := range mediaList {
			if gp.VideoId == *m.GhostID {
				media = &m
				break
			}
		}

		if media == nil {
			accErrs = fmt.Errorf("no media found in exorcist while transferring progress for ghost id: %v\n%w", gp.VideoId, accErrs)
			continue
		}

		progressList = append(progressList, model.MediaProgress{
			GhostID:   &gp.ID,
			UserID:    user.ID,
			MediaID:   media.ID,
			Timestamp: float64(gp.Timestamp) / 1000,
		})
	}

	if accErrs != nil {
		log.Printf("warning: errors found while creating media progress entities: %v", accErrs.Error())
	}

	if len(progressList) == 0 {
		log.Println("Progress list was empty to insert to exorcist")
	}

	insertStmnt := table.MediaProgress.INSERT(table.MediaProgress.GhostID, table.MediaProgress.UserID, table.MediaProgress.MediaID, table.MediaProgress.Timestamp).
		MODELS(progressList).
		ON_CONFLICT(table.MediaProgress.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows in media pgrogress exorcist", rows)

	return nil
}

func (ctx Context) transferRelatedVideos() error {
	log.Println("Transferring related videos")

	relatedVideosStmnt := gtable.RelatedVideos.SELECT(gtable.RelatedVideos.AllColumns)

	var relatedVideos []gmodel.RelatedVideos
	if err := relatedVideosStmnt.Query(ctx.GhostDb, &relatedVideos); err != nil {
		return err
	}

	if len(relatedVideos) == 0 {
		log.Println("No related videos found in ghost")
		return nil
	}

	log.Printf("Found %v related videos in ghost", len(relatedVideos))

	mediaStmnt := table.Media.SELECT(table.Media.GhostID, table.Media.ID)

	var mediaList []model.Media
	if err := mediaStmnt.Query(ctx.ExorcistDb, &mediaList); err != nil {
		return err
	}

	if len(mediaList) == 0 {
		log.Println("No media in exorcist")
		return nil
	}

	var accErrs error
	var relatedMediaList []model.MediaRelation
	for _, rv := range relatedVideos {
		var media *model.Media
		var relatedTo *model.Media
		for _, m := range mediaList {
			if *m.GhostID == rv.VideoId {
				media = &m
			}
			if *m.GhostID == rv.RelatedToId {
				relatedTo = &m
			}

			if media != nil && relatedTo != nil {
				break
			}
		}

		if media == nil {
			accErrs = fmt.Errorf("media was not found while transferring relations by ghost id in exorcist: %v\n%w", rv.VideoId, accErrs)
			continue
		}

		if relatedTo == nil {
			accErrs = fmt.Errorf("relatedTo was not found while transferring relations by ghost id in exorcist: %v\n%w", rv.RelatedToId, accErrs)
			continue
		}

		relatedMediaList = append(relatedMediaList, model.MediaRelation{
			GhostID:      &rv.ID,
			MediaID:      media.ID,
			RelatedTo:    relatedTo.ID,
			RelationType: model.MediaRelationTypeEnum_Media,
		})
	}

	if accErrs != nil {
		log.Printf("warning: some errors while creating media relation entities: %v", accErrs.Error())
	}

	if len(relatedMediaList) == 0 {
		log.Println("no media relations to add into exorcist")
		return nil
	}

	insertStmnt := table.MediaRelation.INSERT(table.MediaRelation.GhostID, table.MediaRelation.MediaID, table.MediaRelation.RelatedTo, table.MediaRelation.RelationType).
		MODELS(relatedMediaList).
		ON_CONFLICT(table.MediaRelation.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	log.Printf("Altered %v rows in media relation exorcist", rows)

	return nil
}

func (ctx Context) transferVideoActors() error {
	log.Println("Transferring video actors")

	stmnt := gtable.VideoActors.SELECT(gtable.VideoActors.AllColumns)

	var videoActors []gmodel.VideoActors
	if err := stmnt.Query(ctx.GhostDb, &videoActors); err != nil {
		return err
	}

	if len(videoActors) == 0 {
		log.Println("No video actors found in ghost")
		return nil
	}

	log.Printf("Found %v video actors in ghost", len(videoActors))

	mediaStmnt := table.Media.SELECT(table.Media.ID, table.Media.GhostID)

	var mediaList []model.Media
	if err := mediaStmnt.Query(ctx.ExorcistDb, &mediaList); err != nil {
		return err
	}

	if len(mediaList) == 0 {
		log.Println("No media found in exorcist")
		return nil
	}

	peopleStmnt := table.Person.SELECT(table.Person.ID, table.Person.GhostID)

	var people []model.Person
	if err := peopleStmnt.Query(ctx.ExorcistDb, &people); err != nil {
		return err
	}

	if len(people) == 0 {
		log.Println("No people found in exorcist")
		return nil
	}

	var accErrs error
	var mediaPeople []model.MediaPerson
	for _, va := range videoActors {
		var person *model.Person
		for _, p := range people {
			if va.ActorId == *p.GhostID {
				person = &p
				break
			}
		}

		if person == nil {
			accErrs = fmt.Errorf("no person found while transferring video actors in exorcist: %v\n%w", va.ActorId, accErrs)
			continue
		}

		var media *model.Media
		for _, m := range mediaList {
			if va.VideoId == *m.GhostID {
				media = &m
				break
			}
		}

		if media == nil {
			accErrs = fmt.Errorf("no media found while transferring video actors in exorcist: %v\n%w", va.VideoId, accErrs)
			continue
		}

		mediaPeople = append(mediaPeople, model.MediaPerson{
			GhostID:  &va.ID,
			PersonID: person.ID,
			MediaID:  media.ID,
		})
	}

	if accErrs != nil {
		log.Printf("warning: some errors were found while creating media person entities for exorcist: %v", accErrs.Error())
	}

	if len(mediaPeople) == 0 {
		log.Println("No media person entities created while transferring video actors")
		return nil
	}

	insertStmnt := table.MediaPerson.INSERT(table.MediaPerson.GhostID, table.MediaPerson.PersonID, table.MediaPerson.MediaID).
		MODELS(mediaPeople).
		ON_CONFLICT(table.MediaPerson.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows while transferring video actors to exorcist", rows)

	return nil
}

func (ctx Context) transferVideoGenres() error {
	log.Println("Transferring video genres")

	stmnt := gtable.VideoGenres.SELECT(gtable.VideoGenres.AllColumns)

	var vgs []gmodel.VideoGenres
	if err := stmnt.Query(ctx.GhostDb, &vgs); err != nil {
		return err
	}

	if len(vgs) == 0 {
		log.Println("No video genres found in ghost")
		return nil
	}

	stmnt = gtable.Videos.SELECT(gtable.Videos.AllColumns)

	var vs []gmodel.Videos
	if err := stmnt.Query(ctx.GhostDb, &vs); err != nil {
		return err
	}

	var VideoGenres []gmodel.VideoGenres

	for _, vg := range vgs {
		var v *gmodel.Videos
		for _, vid := range vs {
			if vid.ID == vg.VideoId {
				v = &vid
				break
			}
		}

		if v != nil {
			VideoGenres = append(VideoGenres, vg)
		}
	}

	mediaStmnt := table.Media.SELECT(table.Media.ID, table.Media.GhostID)

	var mediaList []model.Media
	if err := mediaStmnt.Query(ctx.ExorcistDb, &mediaList); err != nil {
		return err
	}

	if len(mediaList) == 0 {
		log.Println("No media found in exorcist")
		return nil
	}

	tagStmnt := table.Tag.SELECT(table.Tag.ID, table.Tag.GhostID)

	var tags []model.Tag
	if err := tagStmnt.Query(ctx.ExorcistDb, &tags); err != nil {
		return err
	}

	if len(tags) == 0 {
		log.Println("No tags found in exorcist")
		return nil
	}

	var accErrs error
	var mediaTags []model.MediaTag
	for i, va := range vgs {
		var tag *model.Tag
		for _, p := range tags {
			if va.GenreId == *p.GhostID {
				tag = &p
				break
			}
		}

		if tag == nil {
			accErrs = fmt.Errorf("no tag found while transferring video genres in exorcist: %v\n%w", va.GenreId, accErrs)
			continue
		}

		var media *model.Media
		for _, m := range mediaList {
			if va.VideoId == *m.GhostID {
				media = &m
				break
			}
		}

		if media == nil {
			accErrs = fmt.Errorf("no media found while transferring video genres in exorcist: %v\n%w", va.VideoId, accErrs)
			continue
		}

		mediaTags = append(mediaTags, model.MediaTag{
			GhostID: &va.ID,
			TagID:   tag.ID,
			MediaID: media.ID,
		})

		if i%ctx.Config.BatchSize == 0 {
			if err := ctx.insertMediaTags(mediaTags); err != nil {
				return err
			}

			mediaTags = []model.MediaTag{}

			runtime.GC()
		}
	}

	if len(mediaTags) != 0 {
		if err := ctx.insertMediaTags(mediaTags); err != nil {
			return err
		}
	}

	if accErrs != nil {
		log.Printf("warning: some errors were found while creating media tag entities for exorcist: %v", accErrs.Error())
	}

	if len(mediaTags) == 0 {
		log.Println("No media tag entities created while transferring video genres")
		return nil
	}

	return nil
}

func (ctx *Context) insertMediaTags(mediaTags []model.MediaTag) error {
	insertStmnt := table.MediaTag.INSERT(table.MediaTag.GhostID, table.MediaTag.TagID, table.MediaTag.MediaID).
		MODELS(mediaTags).
		ON_CONFLICT(table.MediaTag.GhostID).DO_NOTHING()

	res, err := insertStmnt.Exec(ctx.ExorcistDb)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()

	log.Printf("Altered %v rows while transferring video genres to exorcist", rows)

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

	err = ctx.transferFavouriteActors()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferFavouriteVideos()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferProgress()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferRelatedVideos()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferVideoActors()
	if err != nil {
		errs.PanicError(err)
	}

	err = ctx.transferVideoGenres()
	if err != nil {
		errs.PanicError(err)
	}
}
