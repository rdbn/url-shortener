package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"
	"url-shortener/internal/domain"

	"github.com/lib/pq"
)

type URLPostgres struct {
	db *sql.DB
}

func NewURLPostgres(
	postgresDSN string,
	maxOpenConns, maxIdleConns int,
	connMaxLifetime time.Duration,
) *URLPostgres {
	postgresDB, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		panic(err)
	}

	postgresDB.SetMaxOpenConns(maxOpenConns)
	postgresDB.SetMaxIdleConns(maxIdleConns)
	postgresDB.SetConnMaxLifetime(connMaxLifetime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = postgresDB.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	return &URLPostgres{db: postgresDB}
}

func (s *URLPostgres) Save(ctx context.Context, url *domain.URL) error {
	query := "INSERT INTO urls (id, short_code, original_url, clicks) VALUES ($1, $2, $3, $4)"

	_, err := s.db.ExecContext(
		ctx,
		query,
		url.ID,
		url.ShortCode,
		url.OriginalURL,
		url.Clicks,
	)
	if err != nil {
		isDuplicate := false
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			isDuplicate = pgErr.Code == "23505"
		}

		if isDuplicate {
			return domain.ErrDuplicate
		}

		log.Fatal(err)

		return domain.ErrSaveURL
	}

	return nil
}

func (s *URLPostgres) Get(ctx context.Context, shortCode string) (*domain.URL, error) {
	query := "SELECT id, short_code, original_url, created_at, clicks, updated_at FROM urls WHERE short_code = $1"

	var url domain.URL

	err := s.db.QueryRowContext(ctx, query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.CreatedAt,
		&url.Clicks,
		&url.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}

	if err != nil {
		return nil, domain.ErrGetURL
	}

	return &url, nil
}

func (s *URLPostgres) SaveBatch(ctx context.Context, batch map[string]domain.URL) error {
	query := "UPDATE urls SET clicks = $1, updated_at = now() WHERE short_code = $2"

	tx, err := s.db.Begin()
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			panic(err)
		}
	}(tx)
	if err != nil {
		return err
	}

	for shortCode, urlDomain := range batch {
		_, err := tx.ExecContext(ctx, query, urlDomain.Clicks, shortCode)
		if err != nil {
			return domain.ErrBatchIncrementClicks
		}
	}

	return nil
}

func (s *URLPostgres) Close() {
	err := s.db.Close()
	if err != nil {
		panic(err)
	}
}
