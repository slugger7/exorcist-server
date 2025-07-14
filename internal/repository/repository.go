package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/logger"
	imageRepository "github.com/slugger7/exorcist/internal/repository/image"
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	mediaRepository "github.com/slugger7/exorcist/internal/repository/media"
	personRepository "github.com/slugger7/exorcist/internal/repository/person"
	tagRepository "github.com/slugger7/exorcist/internal/repository/tag"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type IRepository interface {
	Health() map[string]string

	Close() error

	Job() jobRepository.IJobRepository
	Library() libraryRepository.LibraryRepository
	LibraryPath() libraryPathRepository.ILibraryPathRepository
	Video() videoRepository.IVideoRepository
	User() userRepository.IUserRepository
	Image() imageRepository.IImageRepository
	Media() mediaRepository.IMediaRepository
	Person() personRepository.PersonRepository
	Tag() tagRepository.TagRepository
}

type repository struct {
	db              *sql.DB
	logger          logger.ILogger
	env             *environment.EnvironmentVariables
	jobRepo         jobRepository.IJobRepository
	libraryRepo     libraryRepository.LibraryRepository
	libraryPathRepo libraryPathRepository.ILibraryPathRepository
	videoRepo       videoRepository.IVideoRepository
	userRepo        userRepository.IUserRepository
	imageRepo       imageRepository.IImageRepository
	mediaRepo       mediaRepository.IMediaRepository
	personRepo      personRepository.PersonRepository
	tagRepo         tagRepository.TagRepository
}

var dbInstance *repository

func New(env *environment.EnvironmentVariables, context context.Context) IRepository {
	if dbInstance == nil {
		psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.DatabaseHost,
			env.DatabasePort,
			env.DatabaseUser,
			env.DatabasePassword,
			env.DatabaseName)
		if env.AppEnv == environment.AppEnvEnum.Local {
			log.Printf("connection_string: %v", psqlconn)
		}
		db, err := sql.Open("postgres", psqlconn)
		errs.PanicError(err)

		dbInstance = &repository{
			db:              db,
			env:             env,
			logger:          logger.New(env),
			jobRepo:         jobRepository.New(db, env, context),
			libraryRepo:     libraryRepository.New(db, env, context),
			libraryPathRepo: libraryPathRepository.New(db, env, context),
			videoRepo:       videoRepository.New(db, env, context),
			userRepo:        userRepository.New(db, env, context),
			imageRepo:       imageRepository.New(db, env, context),
			mediaRepo:       mediaRepository.New(db, env, context),
			personRepo:      personRepository.New(env, db, context),
			tagRepo:         tagRepository.New(env, db, context),
		}

		err = dbInstance.runMigrations()
		if err != nil {
			dbInstance.logger.Warningf("Migrations were not run because %v", err)
		}
		dbInstance.logger.Info("Database instance created")
	}

	return dbInstance
}

func (s *repository) Job() jobRepository.IJobRepository {
	s.logger.Debug("Getting job repo")
	return s.jobRepo
}

func (s *repository) Library() libraryRepository.LibraryRepository {
	s.logger.Debug("Getting library repo")
	return s.libraryRepo
}

func (s *repository) LibraryPath() libraryPathRepository.ILibraryPathRepository {
	s.logger.Debug("Getting library path repo")
	return s.libraryPathRepo
}

func (s *repository) Video() videoRepository.IVideoRepository {
	s.logger.Debug("Getting video repo")
	return s.videoRepo
}

func (s *repository) User() userRepository.IUserRepository {
	s.logger.Debug("Getting user repo")
	return dbInstance.userRepo
}

func (s *repository) Image() imageRepository.IImageRepository {
	s.logger.Debug("Getting image repo")
	return dbInstance.imageRepo
}

func (s *repository) Media() mediaRepository.IMediaRepository {
	s.logger.Debug("Getting media repo")
	return dbInstance.mediaRepo
}

func (s *repository) Person() personRepository.PersonRepository {
	s.logger.Debug("Getting person repo")
	return dbInstance.personRepo
}

func (s *repository) Tag() tagRepository.TagRepository {
	s.logger.Debug("Getting tag repo")
	return dbInstance.tagRepo
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *repository) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		s.logger.Errorf("db down: %v", err) // Log the error and terminate the program
		errs.PanicError(err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *repository) Close() error {
	log.Printf("Disconnected from database: %s", s.env.DatabaseName)
	return s.db.Close()
}

func (s *repository) runMigrations() error {
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	s.logger.Info("Running migrations")
	err = m.Up()
	if err != nil {
		return err
	}
	s.logger.Info("Migrations completed")
	return nil
}
