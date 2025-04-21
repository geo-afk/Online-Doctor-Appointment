package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/geo-afk/Online-Doctor-Appointment/internal/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	DB() *postgres.Queries
	Close() error
}

type service struct {
	query *postgres.Queries
	pool  *pgxpool.Pool
}

func New() Service {

	dbUrl := os.Getenv("DATABASE_URL")
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database pool: %v\n", err)
	}

	dbInstance := &service{
		query: postgres.New(pool),
		pool:  pool,
	}

	return dbInstance

}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Perform a simple query to check the connection
	if s.pool != nil {
		poolStats := s.pool.Stat()
		stats["acquired"] = strconv.Itoa(int(poolStats.AcquiredConns()))
		stats["idle"] = strconv.Itoa(int(poolStats.IdleConns()))
		stats["total"] = strconv.Itoa(int(poolStats.TotalConns()))
		return stats
	} else {
		stats["status"] = "unknown"
		stats["error"] = "no database connection available"
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"
	return stats
}

func (s *service) DB() *postgres.Queries {
	return s.query
}

func (s *service) Close() error {
	s.pool.Close()
	return nil
}
