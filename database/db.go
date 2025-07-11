package database

import (
	"fmt"
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"monitron-server/config"

	migrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitDB(cfg *config.Config) *gorm.DB {
	dbURI := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)

	var err error
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
			log.Fatal().Err(err).Msg("Failed to connect to database")

	log.Info().Msg("Successfully connected to database!")

	// Run migrations using the sqlx driver for golang-migrate
	m, err := migrate.New(
		"file://database/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create migrate instance")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal().Err(err).Msg("Failed to run migrations")

	log.Info().Msg("Database migrations applied successfully!")
	return db
}

func CloseDB(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
				log.Error().Err(err).Msg("Error getting *sql.DB from GORM for closing")
				return
		}
		sqlDB.Close()
		log.Info().Msg("Database connection closed.")
	}
}
