package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"web-storage-service/pkg"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func init() {
	// TODO: Если использовать "github.com/joho/godotenv/autoload",
	// TODO: можно будет убрать init() функции
	err := pkg.LoadEnv(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port = os.Getenv("DB_PORT")
	host = os.Getenv("DB_HOST")
	schema = os.Getenv("DB_SCHEMA")
}

// Service represents a service that interacts with a database.
type Service interface {
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)

	ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)

	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row

	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)

	BeginTx(ctx context.Context) (pgx.Tx, error)

	Health() map[string]string

	Close()
}

type service struct {
	db *pgxpool.Pool
}

var (
	database   string
	password   string
	username   string
	port       string
	host       string
	schema     string
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}
	dbInstance = &service{
		db: pool,
	}
	return dbInstance
}

func (s *service) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return s.db.Exec(ctx, query, args...)
}

func (s *service) ExecContext(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	return s.db.Exec(ctx, query, args...)
}

func (s *service) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *service) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return s.db.QueryRow(ctx, query, args...)
}

func (s *service) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.ReadCommitted,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err))
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := s.db.Stat()
	stats["open_connections"] = strconv.Itoa(int(dbStats.NewConnsCount()))
	stats["idle"] = strconv.Itoa(int(dbStats.IdleConns()))
	stats["acquire_duration"] = dbStats.AcquireDuration().String()
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeDestroyCount(), 10)

	if dbStats.AcquiredConns() > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.MaxLifetimeDestroyCount() > int64(dbStats.AcquiredConns())/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *service) Close() {
	log.Printf("Disconnected from database: %s", database)
	defer s.db.Close()
}
