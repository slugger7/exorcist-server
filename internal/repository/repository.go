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
	jobRepository "github.com/slugger7/exorcist/internal/repository/job"
	libraryRepository "github.com/slugger7/exorcist/internal/repository/library"
	libraryPathRepository "github.com/slugger7/exorcist/internal/repository/library_path"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
	videoRepository "github.com/slugger7/exorcist/internal/repository/video"
)

type IRepository interface {
	Health() map[string]string

	Close() error

	Job() jobRepository.IJobRepository
	Library() libraryRepository.ILibraryRepository
	LibraryPath() libraryPathRepository.ILibraryPathRepository
	Video() videoRepository.IVideoRepository
	User() userRepository.IUserRepository
}

type Repository struct {
	db              *sql.DB
	logger          logger.ILogger
	Env             *environment.EnvironmentVariables
	jobRepo         jobRepository.IJobRepository
	libraryRepo     libraryRepository.ILibraryRepository
	libraryPathRepo libraryPathRepository.ILibraryPathRepository
	videoRepo       videoRepository.IVideoRepository
	userRepo        userRepository.IUserRepository
}

var dbInstance *Repository

func New(env *environment.EnvironmentVariables) IRepository {
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

		dbInstance = &Repository{
			db:              db,
			Env:             env,
			logger:          logger.New(env),
			jobRepo:         jobRepository.New(db, env),
			libraryRepo:     libraryRepository.New(db, env),
			libraryPathRepo: libraryPathRepository.New(db, env),
			videoRepo:       videoRepository.New(db, env),
			userRepo:        userRepository.New(db, env),
		}

		err = dbInstance.runMigrations()
		if err != nil {
			dbInstance.logger.Warningf("Migrations were not run because %v", err)
		}
		dbInstance.logger.Info("Database instance created")
	}

	return dbInstance
}

func (s *Repository) Job() jobRepository.IJobRepository {
	return s.jobRepo
}

func (s *Repository) Library() libraryRepository.ILibraryRepository {
	return s.libraryRepo
}

func (s *Repository) LibraryPath() libraryPathRepository.ILibraryPathRepository {
	return s.libraryPathRepo
}

func (s *Repository) Video() videoRepository.IVideoRepository {
	return s.videoRepo
}

func (s *Repository) User() userRepository.IUserRepository {
	return dbInstance.userRepo
}

func (s *Repository) Db() *sql.DB {
	return dbInstance.db
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *Repository) Health() map[string]string {
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
func (s *Repository) Close() error {
	log.Printf("Disconnected from database: %s", s.Env.DatabaseName)
	return s.db.Close()
}

func (s *Repository) runMigrations() error {
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
